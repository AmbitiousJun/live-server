package https

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/util/encrypts"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
)

// CacheClient 具有缓存功能的 http 请求客户端
type CacheClient struct {
	caches       map[string]httpRespCache // 存放缓存响应
	cacheOpMutex sync.RWMutex             // 操作缓存时需要获取的锁
	expiredTime  int64                    // 过期时间 (毫秒)
	maxCacheNum  int                      // 最大缓存个数
}

// httpRespCache http 请求的响应缓存
type httpRespCache struct {
	code      int
	header    http.Header
	body      []byte
	finalUrl  string // 记录请求重定向后的最终地址
	timestamp int64  // 加入缓存的时间
}

// NewCacheClient 初始化一个具备缓存特性的 http 客户端
//
// 一个请求缓存的过期时间可能存在误差
func NewCacheClient(maxCacheNum int, expiredTime time.Duration) *CacheClient {
	// 校验参数
	if maxCacheNum < 0 {
		log.Panicf("最大缓存个数不能为负数: %d", maxCacheNum)
	}

	// 创建实例
	client := &CacheClient{
		expiredTime: expiredTime.Milliseconds(),
		maxCacheNum: maxCacheNum,
		caches:      make(map[string]httpRespCache),
	}

	// 创建一个定时器, 每隔指定时间清除缓存
	ticker := time.NewTicker(expiredTime)
	go func() {
		for range ticker.C {
			toRemoves := make([]string, 0)
			// 将过期的缓存筛选出来
			for key, cache := range client.caches {
				if cache.timestamp+client.expiredTime > time.Now().UnixMilli() {
					continue
				}
				toRemoves = append(toRemoves, key)
			}
			// 逐个移除过期缓存
			for _, key := range toRemoves {
				client.delCache(key)
			}
		}
	}()

	return client
}

// Request 发起 http 请求获取响应
//
// 如果一个请求有多次重定向并且进行了 autoRedirect,
// 则最后一次重定向的 url 会作为第一个参数返回
func (cc *CacheClient) Request(method, url string, header http.Header, body io.ReadCloser, autoRedirect bool) (string, *http.Response, error) {
	// 计算当前请求的缓存 key
	strBody := ""
	if body != nil {
		reqBytes, err := io.ReadAll(body)
		if err != nil {
			return "", nil, fmt.Errorf("读取请求体失败: %v", err)
		}
		strBody = string(reqBytes)
	}
	cacheKey := cc.cacheKey(url, method, strBody, header, autoRedirect)

	// 判断是否命中缓存
	if cache, ok := cc.getCache(cacheKey); ok {
		resp := new(http.Response)
		resp.StatusCode = cache.code
		resp.Body = io.NopCloser(bytes.NewBuffer(cache.body))
		resp.Header = cache.header.Clone()
		return cache.finalUrl, resp, nil
	}

	// 发起请求
	finalUrl, resp, err := Request(method, url, header, io.NopCloser(bytes.NewBufferString(strBody)), autoRedirect)

	// 请求异常或者非成功响应码, 不进行进一步处理
	if err != nil || !IsSuccessCode(resp.StatusCode) {
		return finalUrl, resp, err
	}

	// 拷贝一份响应体出来
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return finalUrl, resp, fmt.Errorf("读取响应体失败: %v", err)
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// 拷贝响应头
	cloneHeader := resp.Header.Clone()

	// 异步进行缓存
	go func() {
		cc.putCache(cacheKey, httpRespCache{
			code:      resp.StatusCode,
			header:    cloneHeader,
			body:      bodyBytes,
			finalUrl:  finalUrl,
			timestamp: time.Now().UnixMilli(),
		})
	}()

	return finalUrl, resp, nil
}

// cacheKey 计算某个请求的缓存 hash key
func (cc *CacheClient) cacheKey(url, method, body string, header http.Header, autoRedirect bool) string {
	// 将 code、header、body 拼接在一起, 按照字典序排序后计算出 md5 hash 值
	sb := strings.Builder{}
	sb.WriteString(method)
	sb.WriteString(body)
	sb.WriteString(fmt.Sprintf("%v", autoRedirect))
	for key, values := range header {
		sb.WriteString(fmt.Sprintf("%s: %v", key, values))
	}
	hash := encrypts.Md5Hash(strs.Sort(sb.String()))

	// 将请求原始的 uri 拼接在 hash 值前面, 防止出现偶然的 hash 碰撞
	return url + hash
}

// delCache 删除缓存
func (cc *CacheClient) delCache(key string) {
	cc.cacheOpMutex.Lock()
	defer cc.cacheOpMutex.Unlock()
	delete(cc.caches, key)
}

// putCache 设置缓存
func (cc *CacheClient) putCache(key string, cache httpRespCache) {
	cc.cacheOpMutex.Lock()
	defer cc.cacheOpMutex.Unlock()

	// 如果缓存有这个 key, 直接覆盖
	if _, ok := cc.caches[key]; ok {
		cc.caches[key] = cache
		return
	}

	// 如果缓存已经满了, 就放弃缓存
	if len(cc.caches) >= cc.maxCacheNum {
		return
	}

	cc.caches[key] = cache
}

// getCache 获取缓存
//
// 为了保持较好的性能, 如果读锁获取不到, 不进行阻塞等待
func (cc *CacheClient) getCache(key string) (httpRespCache, bool) {
	if !cc.cacheOpMutex.TryRLock() {
		return httpRespCache{}, false
	}
	defer cc.cacheOpMutex.RUnlock()

	if cache, ok := cc.caches[key]; ok {
		return cache, true
	}

	return httpRespCache{}, false
}
