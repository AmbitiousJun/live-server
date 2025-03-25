package handler

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
)

// combineChCache 聚合频道缓存
type combineChCache struct {

	// idx 频道索引缓存
	idx int

	// createTime 缓存创建时间
	createTime time.Time
}

// combineHandler 接口聚合处理器
type combineHandler struct {

	// caches 缓存
	caches map[string]combineChCache

	// cacheTimeout 频道缓存过期时间
	cacheTimeout time.Duration

	// cacheOpMutex 控制缓存并发读写
	cacheOpMutex sync.RWMutex

	// chSeg 分割多个直播源的分隔符
	chSeg string
}

func init() {
	y := &combineHandler{
		chSeg:        "{{{:}}}",
		cacheTimeout: time.Second * 10,
		caches:       make(map[string]combineChCache),
	}
	resolve.RegisterHandler(y)
	go func() {
		ticker := time.NewTicker(time.Second)
		for range ticker.C {
			if len(y.caches) == 0 {
				continue
			}

			y.cacheOpMutex.Lock()
			toDeletes := make([]string, 0)
			for k, v := range y.caches {
				if time.Since(v.createTime) > y.cacheTimeout {
					toDeletes = append(toDeletes, k)
				}
			}
			for _, del := range toDeletes {
				delete(y.caches, del)
			}
			y.cacheOpMutex.Unlock()
		}
	}()
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *combineHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 1 取出环境变量
	if strs.AnyEmpty(params.UrlEnv) {
		return resolve.HandleResult{}, errors.New("请传入环境变量")
	}
	envVal, ok := env.Get(params.UrlEnv)
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("环境变量 %s 不存在", params.UrlEnv)
	}

	// 2 解析直播源列表
	chs := strings.Split(envVal, h.chSeg)
	if len(chs) == 0 {
		return resolve.HandleResult{}, fmt.Errorf("直播源列表为空, 请使用 %s 分隔多个直播源", h.chSeg)
	}

	// 3 缓存检测
	if cacheIdx, ok := h.readCacheIdx(params.UrlEnv); ok && cacheIdx < len(chs) && cacheIdx >= 0 {
		return resolve.HandleResult{Type: resolve.ResultRedirect, Url: chs[cacheIdx]}, nil
	}

	// 4 随机选择一个直播源
	randIdx := rand.IntN(len(chs))
	h.putCacheIdx(params.UrlEnv, randIdx)
	return resolve.HandleResult{
		Type: resolve.ResultRedirect,
		Url:  chs[randIdx],
	}, nil
}

// Name 处理器名称
func (h *combineHandler) Name() string {
	return "combine"
}

// HelpDoc 处理器说明文档
func (h *combineHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 本处理器作用是聚合多个直播源, 并在客户端请求时随机返回重定向")
	sb.WriteString("\n2. 先将自定义的多个直播源设置到程序的环境变量中, 变量名称任意, 如: combine_hyxw")
	sb.WriteString("\n3. 多个源之间使用 " + h.chSeg + " 进行分隔")
	sb.WriteString("\n4. 变量设置示例: ${clientOrigin}/handler/aktv/ch/寰宇新聞台?proxy_m3u=1&proxy_ts=1&ts_proxy_mode=custom$AKTV" + h.chSeg + "${clientOrigin}/handler/youtube/ch/6IquAgfvYmc?proxy_m3u=1&proxy_ts=1&ts_proxy_mode=custom&format=FHD$YT" + h.chSeg + "${clientOrigin}/handler/345/ch/hyxw?proxy_m3u=1&proxy_ts=1&format=1&ts_proxy_mode=custom$345")
	sb.WriteString("\n5. 频道名称 ch 可任意传递一个非空值")
	sb.WriteString("\n6. 请求示例: ${clientOrigin}/handler/combine/ch/1?url_env=combine_hyxw")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *combineHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *combineHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (h *combineHandler) Enabled() bool {
	return true
}

// readCacheIdx 读取缓存中的频道索引
func (h *combineHandler) readCacheIdx(env string) (int, bool) {
	h.cacheOpMutex.RLock()
	defer h.cacheOpMutex.RUnlock()

	if cache, ok := h.caches[env]; ok {
		return cache.idx, true
	}

	return 0, false
}

// putCacheIdx 写入频道索引缓存
func (h *combineHandler) putCacheIdx(env string, idx int) {
	h.cacheOpMutex.Lock()
	defer h.cacheOpMutex.Unlock()
	h.caches[env] = combineChCache{idx: idx, createTime: time.Now()}
}
