package https

import (
	"bytes"
	"errors"
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
	caches      map[string]httpRespCache // 存放缓存响应
	pendings    map[string]*requestSync  // 存放正在进行的请求
	mu          sync.RWMutex             // 操作缓存时需要获取的锁
	expiredTime int64                    // 过期时间 (毫秒)
	maxCacheNum int                      // 最大缓存个数
}

// httpRespCache http 请求的响应缓存
type httpRespCache struct {
	code      int
	header    http.Header
	body      []byte
	finalUrl  string // 记录请求重定向后的最终地址
	timestamp int64  // 加入缓存的时间
}

// requestSync 对相同 cache key 的请求进行同步
type requestSync struct {
	mu   sync.Mutex
	cond *sync.Cond
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
		pendings:    make(map[string]*requestSync),
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
		return cache.finalUrl, cc.useCacheResp(cache), nil
	}

	// 判断并标记请求为进行中状态
	if state, pending := cc.getOrMarkPending(cacheKey); pending {
		// 二次检查缓存是否命中, 防止 broadcast 提前触发导致线程睡死
		if cache, ok := cc.getCache(cacheKey); ok {
			cc.removePending(cacheKey)
			return cache.finalUrl, cc.useCacheResp(cache), nil
		}

		// 等待并发的同个请求执行完成
		state.mu.Lock()
		defer state.mu.Unlock()
		state.cond.Wait()

		// 再次判断缓存是否命中
		if cache, ok := cc.getCache(cacheKey); ok {
			cc.removePending(cacheKey)
			return cache.finalUrl, cc.useCacheResp(cache), nil
		}
		return "", nil, errors.New("cache client 状态异常, 无法命中缓存")
	}

	defer func() {
		// 通知等待当前请求的其他 goroutine
		if state, ok := cc.pendingState(cacheKey); ok {
			state.cond.Broadcast()
		}
		cc.removePending(cacheKey)
	}()

	// 发起请求
	finalUrl, resp, err := Request(method, url, header, io.NopCloser(bytes.NewBufferString(strBody)), autoRedirect)
	var bodyBytes []byte

	// 错误响应
	if err != nil || !IsSuccessCode(resp.StatusCode) {
		return finalUrl, resp, err
	}

	// 缓存请求
	defer func() {
		var statusCode int
		var respHeader http.Header
		if resp != nil {
			statusCode = resp.StatusCode
			if resp.Header != nil {
				respHeader = resp.Header.Clone()
			}
		}
		cc.putCache(cacheKey, httpRespCache{
			code:      statusCode,
			header:    respHeader,
			body:      bodyBytes,
			finalUrl:  finalUrl,
			timestamp: time.Now().UnixMilli(),
		})
	}()

	// 拷贝一份响应体出来
	bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return finalUrl, resp, fmt.Errorf("读取响应体失败: %v", err)
	}
	resp.Body.Close()
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

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
	cc.mu.Lock()
	defer cc.mu.Unlock()
	delete(cc.caches, key)
}

// putCache 设置缓存
func (cc *CacheClient) putCache(key string, cache httpRespCache) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

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
func (cc *CacheClient) getCache(key string) (httpRespCache, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if cache, ok := cc.caches[key]; ok {
		return cache, true
	}

	return httpRespCache{}, false
}

// pendingState 获取请求状态
func (cc *CacheClient) pendingState(key string) (*requestSync, bool) {
	cc.mu.RLock()
	defer cc.mu.RUnlock()

	if state, ok := cc.pendings[key]; ok && state != nil {
		return state, true
	}

	return nil, false
}

// getOrMarkPending 获取请求的进行中状态, 如果不存在则同时标记为进行中
//
// 如果请求是首次被标记为进行中状态, 第二个参数会返回 false
func (cc *CacheClient) getOrMarkPending(key string) (*requestSync, bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()

	// 已经是进行中状态直接返回
	if state, ok := cc.pendings[key]; ok && state != nil {
		return state, true
	}

	// 将请求标记为进行中状态并返回
	state := new(requestSync)
	state.cond = sync.NewCond(&state.mu)
	cc.pendings[key] = state
	return state, false
}

// removePending 移除请求进行中的状态
func (cc *CacheClient) removePending(key string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	delete(cc.pendings, key)
}

// useCacheResp 使用缓存响应
func (cc *CacheClient) useCacheResp(cache httpRespCache) *http.Response {
	resp := new(http.Response)
	resp.StatusCode = cache.code
	resp.Body = io.NopCloser(bytes.NewBuffer(cache.body))
	if cache.header != nil {
		resp.Header = cache.header.Clone()
	}
	return resp
}
