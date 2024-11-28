package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/ytdlp"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

const (
	// Env_YoutubeCustomFormatEnable 是否允许解析自定义格式
	Env_YoutubeCustomFormatEnable = "youtube_custom_format_enable"

	// YoutubeResPrefix youtube 资源前缀
	YoutubeResPrefix = "https://www.youtube.com/watch?v="
)

// youtubeHandler Youtube 直播处理器
type youtubeHandler struct {
	cacher *youtubeCacher
}

func init() {
	resolve.RegisterHandler(&youtubeHandler{cacher: NewYoutubeCacher()})
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (y *youtubeHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	formatCode := y.chooseFormat(params.Format)
	playlist, err := y.cacher.GetM3U(params.ChName, formatCode)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("获取 playlist 失败: %v", err)
	}

	if !params.ProxyM3U {
		return resolve.HandleResult{Type: resolve.ResultRedirect, Url: playlist}, nil
	}

	content, err := resolve.ProxyM3U(playlist, nil, params.ProxyTs, params.ClientIp)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("代理 m3u 失败: %s, err: %v", params.ChName, err)
	}

	respHeader := make(http.Header)
	respHeader.Set("Content-Type", "application/vnd.apple.mpegurl")
	return resolve.HandleResult{
		Type:   resolve.ResultProxy,
		Code:   http.StatusOK,
		Header: respHeader,
		Body:   []byte(content),
	}, nil
}

// Name 处理器名称
func (y *youtubeHandler) Name() string {
	return "youtube"
}

// HelpDoc 处理器说明文档
func (y *youtubeHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 频道名称即为 youtube 直播的 id (在浏览器地址上查看)")
	sb.WriteString("\n2. 使用该处理器时，服务端必须具有外网环境")
	sb.WriteString("\n3. 如果客户端没有外网环境，可开启 【m3u + ts 切片】 代理模式，让服务端代理流量")
	sb.WriteString("\n4. 默认只支持 HD(720p) 直播格式解析")
	sb.WriteString("\n5. 设置环境变量 【youtube_custom_format_enable=1】 可开启自定义直播格式解析, 目前只支持 FHD")
	sb.WriteString("\n6. 使用自定义的直播格式需要在调用处理器时传入 query 参数，如：?format=FHD")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (y *youtubeHandler) SupportM3UProxy() bool {
	return true
}

// chooseFormat 选择一个合适的 yt-dlp formatCode
func (y *youtubeHandler) chooseFormat(wantFmt string) string {
	// 默认使用 720 p
	res := "95"

	// 判断是否允许自定义
	enable, ok := env.Get(Env_YoutubeCustomFormatEnable)
	if !ok || enable != "1" {
		return res
	}

	switch wantFmt {
	case "FHD":
		res = "96"
	}

	return res
}

var (

	// youtubeCacheTimeout youtube 频道缓存过期时间
	youtubeCacheTimeout = time.Hour.Milliseconds()
)

// youtubePreCacheRequest 频道预缓存请求
type youtubePreCacheRequest struct {
	chId       string      // 频道 id
	formatCode string      // 频道格式
	resChan    chan string // 用于接收请求完成的 m3u 地址的通道
}

// youtubeCacher m3u 缓存器
type youtubeCacher struct {

	// channels 频道列表, 以 youtube 的 id 作为 key
	channels map[string]*struct {
		url      string // m3u 地址
		lastRead int64  // 最后一次读取时间戳
	}

	// preCacheReqChan 预缓存请求通道
	preCacheReqChan chan youtubePreCacheRequest

	// mu 并发控制
	mu sync.RWMutex
}

// NewYoutubeCacher 创建一个 youtube m3u 缓存器
func NewYoutubeCacher() *youtubeCacher {
	yc := new(youtubeCacher)
	yc.channels = make(map[string]*struct {
		url      string
		lastRead int64
	})
	yc.preCacheReqChan = make(chan youtubePreCacheRequest, 100)

	// 每隔 30 分钟维护一次内存
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for {
			select {
			case <-ticker.C:
				yc.UpdateAll()
			case req := <-yc.preCacheReqChan:
				if req.resChan == nil {
					break
				}
				ch, err := yc.RequestChannel(req.chId, req.formatCode)
				if err != nil {
					log.Printf(colors.ToRed("请求频道信息失败: %v"), err)
					ch = ""
				}
				req.resChan <- ch
				close(req.resChan)
			}
		}
	}()

	return yc
}

// CalcCacheKey 计算频道缓存 key
func (yc *youtubeCacher) CalcCacheKey(chId, formatCode string) string {
	return chId + ":" + formatCode
}

// ResolveCacheKey 解析缓存 key
func (yc *youtubeCacher) ResolveCacheKey(cacheKey string) (chId, formatCode string, valid bool) {
	idx := strings.LastIndex(cacheKey, ":")
	if idx == -1 {
		valid = false
		return
	}
	chId, formatCode, valid = cacheKey[:idx], cacheKey[idx+1:], true
	return
}

// Get 获取指定频道 m3u
//
// 该方法可以接受客户端并发调用
func (yc *youtubeCacher) GetM3U(chId, formatCode string) (m3u string, err error) {
	defer func() {
		if err == nil {
			// 更新最后读取时间
			if chCache, ok := yc.channels[yc.CalcCacheKey(chId, formatCode)]; ok {
				chCache.lastRead = time.Now().UnixMilli()
			}
		}
	}()

	// 缓冲区存在该频道, 直接返回
	yc.mu.RLock()
	chCache, ok := yc.channels[yc.CalcCacheKey(chId, formatCode)]
	yc.mu.RUnlock()
	if ok {
		return chCache.url, nil
	}

	// 排队请求新频道
	req := youtubePreCacheRequest{
		chId:       chId,
		formatCode: formatCode,
		resChan:    make(chan string),
	}
	select {
	case yc.preCacheReqChan <- req:
	default:
		return "", fmt.Errorf("服务器繁忙, 请稍后再试")
	}

	// 等待频道信息返回
	m3u = <-req.resChan
	if m3u == "" {
		return "", errors.New("请求频道信息失败")
	}

	return m3u, nil
}

// UpdateAll 更新缓存, 同时淘汰太长时间未读的列表
func (yc *youtubeCacher) UpdateAll() {
	successCnt, failCnt, removeCnt := 0, 0, 0
	defer func() {
		if successCnt <= 0 && failCnt <= 0 && removeCnt <= 0 {
			return
		}
		log.Printf(colors.ToGreen("youtube 缓存刷新完成, 成功: %d, 失败: %d, 移除: %d"), successCnt, failCnt, removeCnt)
	}()

	toRemoves := make([]string, 0)
	for key, value := range yc.channels {
		if value.lastRead+youtubeCacheTimeout < time.Now().UnixMilli() {
			// 缓存过期
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		chId, formatCode, ok := yc.ResolveCacheKey(key)
		if !ok {
			// 无效的 key
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		_, err := yc.CacheM3U(chId, formatCode)
		if err != nil {
			failCnt++
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		successCnt++
	}

	if len(toRemoves) == 0 {
		return
	}

	// 移除过期缓存
	yc.mu.Lock()
	defer yc.mu.Unlock()

	for _, key := range toRemoves {
		delete(yc.channels, key)
	}
}

// RequestChannel 请求指定频道
//
// 如果缓存中已经有指定频道, 直接返回缓存值
//
// 否则调用 yt-dlp 获取一个新的并进行缓存
func (yc *youtubeCacher) RequestChannel(chId, formatCode string) (string, error) {
	key := yc.CalcCacheKey(chId, formatCode)

	yc.mu.RLock()
	chCache, ok := yc.channels[key]
	yc.mu.RUnlock()
	if ok {
		return chCache.url, nil
	}

	m3u, err := yc.CacheM3U(chId, formatCode)
	if err != nil {
		return "", err
	}

	return m3u, nil
}

// CacheM3U 调用 yt-dlp 获取频道的 m3u 地址
func (yc *youtubeCacher) CacheM3U(chId, formatCode string) (string, error) {
	res, err := ytdlp.Extract(YoutubeResPrefix+chId, formatCode)
	if err != nil {
		return "", fmt.Errorf("调用 yt-dlp 失败: %v", err)
	}

	yc.mu.Lock()
	defer yc.mu.Unlock()

	// 将地址维护到内存中
	chCache, ok := yc.channels[yc.CalcCacheKey(chId, formatCode)]
	if !ok {
		chCache = &struct {
			url      string
			lastRead int64
		}{}
	}

	chCache.url = res
	yc.channels[yc.CalcCacheKey(chId, formatCode)] = chCache
	return res, nil
}
