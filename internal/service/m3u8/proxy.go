package m3u8

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

const (

	// CacheNoReadExpiredDuration 超出这个时间没有客户端来读取时,
	// 就将缓存从缓存池中移除, 并停止刷新缓存
	CacheNoReadExpiredDuration = time.Minute * 10

	// CacheMaintainWorkerMaxCount 控制维护缓存的最大协程数量
	CacheMaintainWorkerMaxCount = 50
)

// ProxyParams 代理参数
type ProxyParams struct {

	// Url 要代理的远程 m3u8 地址
	Url string

	// Header 可选, 发起代理的请求头
	Header http.Header
}

// cacheHolder 存放 m3u8 缓存文本
type cacheHolder struct {

	// finalUrl 最终重定向后的地址
	finalUrl string

	// content 缓存文本
	content string

	// lastReadTime 最后一次读取时间
	lastReadTime time.Time

	// ready 区分尚未获取内容与成功获取到空文本
	ready bool
}

// proxyCacheMu 保护任务注册、缓存内容和读取时间, 避免读取与过期删除竞争。
var proxyCacheMu sync.Mutex

// proxyCacheTasks 按缓存 key 存放任务及缓存, 包括尚未首次获取成功的任务。
// 所有访问均由 proxyCacheMu 保护, 同时用于任务去重和数量限制。
var proxyCacheTasks = make(map[string]*cacheHolder)

// ReadCacheContent 读取缓存池中的 m3u8 文本
//
// 发现没有缓存时, 阻塞等待 10 秒后仍没有结果时, 返回错误
func ReadCacheContent(params ProxyParams) (finalUrl, m3u8Content string, err error) {
	key := calcCacheKey(params)

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for {
		proxyCacheMu.Lock()
		if holder, ok := proxyCacheTasks[key]; ok && holder.ready {
			holder.lastReadTime = time.Now()
			content := holder.content
			finalUrl := holder.finalUrl
			proxyCacheMu.Unlock()
			return finalUrl, content, nil
		}
		proxyCacheMu.Unlock()
		select {
		case <-timer.C:
			return "", "", fmt.Errorf("等待 m3u8 缓存超时 (10 秒): %s", params.Url)
		case <-ticker.C:
		}
	}
}

// PushCacheTask 推送缓存任务, 判断当前没有协程在维护这个任务时,
// 就创建一个新的异步协程, 每隔 1 秒使用 http 请求最新 m3u8 文本并刷新缓存,
// 如果缓存超过 CacheNoReadExpiredDuration 时间没有被读取, 就删除缓存, 终止协程
//
// 每次 http 请求时, 超时时间为 30 秒, 超时就放弃请求并打印日志, 等待下一次计时器周期重新请求
//
// 需要确保总共正在维护的协程数不能超过 CacheMaintainWorkerMaxCount 个
func PushCacheTask(params ProxyParams) error {
	u, err := url.Parse(params.Url)
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return fmt.Errorf("无效的 m3u8 地址: %s", params.Url)
	}

	// 异步任务持有独立副本, 调用方返回后可以继续修改自己的请求头。
	if params.Header != nil {
		params.Header = params.Header.Clone()
	}
	key := calcCacheKey(params)

	proxyCacheMu.Lock()
	defer proxyCacheMu.Unlock()

	if _, ok := proxyCacheTasks[key]; ok {
		return nil
	}
	if len(proxyCacheTasks) >= CacheMaintainWorkerMaxCount {
		return fmt.Errorf("m3u8 缓存维护任务已达上限: %d", CacheMaintainWorkerMaxCount)
	}
	holder := &cacheHolder{lastReadTime: time.Now()}
	proxyCacheTasks[key] = holder
	go maintainCache(key, params, holder)
	return nil
}

// maintainCache 首次立即请求, 此后按秒刷新; 同一任务的请求不会重叠。
func maintainCache(key string, params ProxyParams, holder *cacheHolder) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		proxyCacheMu.Lock()
		if time.Since(holder.lastReadTime) >= CacheNoReadExpiredDuration {
			// 在同一临界区删除缓存, 避免退出前又被客户端读取续期。
			delete(proxyCacheTasks, key)
			proxyCacheMu.Unlock()
			return
		}
		proxyCacheMu.Unlock()

		finalUrl, content, err := fetchCacheContent(params)
		if err != nil {
			log.Printf(colors.ToRed("刷新 m3u8 缓存失败, url: %s, err: %v"), params.Url, err)
		} else {
			proxyCacheMu.Lock()
			holder.content = content
			holder.finalUrl = finalUrl
			holder.ready = true
			proxyCacheMu.Unlock()
		}
		<-ticker.C
	}
}

// fetchCacheContent 超时包含重定向和响应体读取, 失败时不覆盖已有缓存。
func fetchCacheContent(params ProxyParams) (finalUrl, content string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	finalUrl, resp, err := https.Get(params.Url).Header(params.Header).Context(ctx).DoRedirect()
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if !https.IsSuccessCode(resp.StatusCode) {
		return "", "", fmt.Errorf("远程响应状态异常: %s", resp.Status)
	}
	contentBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("读取 m3u8 响应失败: %w", err)
	}
	return finalUrl, string(contentBytes), nil
}

// calcCacheKey 将请求地址和请求头拼接成一个大字符串, 按字典序排序后计算哈希值
func calcCacheKey(params ProxyParams) string {
	// 按字段排序而非打散字符, 避免字符相同但语义不同的地址或请求头串用缓存。
	keys := make([]string, 0, len(params.Header))
	for key := range params.Header {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var data strings.Builder
	fmt.Fprintf(&data, "%d:%s", len(params.Url), params.Url)
	for _, key := range keys {
		values := params.Header[key]
		fmt.Fprintf(&data, "%d:%s%d:", len(key), key, len(values))
		for _, value := range values {
			fmt.Fprintf(&data, "%d:%s", len(value), value)
		}
	}
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data.String())))
}
