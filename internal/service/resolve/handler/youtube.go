package handler

import (
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

	// Env_YoutubePlaylistCache 在环境变量中缓存 youtube 的解析列表,
	// 避免频繁调用 yt-dlp 进行解析
	Env_YoutubePlaylistCache = "youtube_playlist_cache"

	// Env_YoutubeCustomFormatEnable 是否允许解析 FHD 格式
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

	content, err := resolve.ProxyM3U(playlist, nil, params.ProxyTs)
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

// youtubeCacher m3u 缓存器
type youtubeCacher struct {

	// channels 频道列表, 以 youtube 的 id 作为 key
	channels map[string]*struct {
		url      string // m3u 地址
		lastRead int64  // 最后一次读取时间戳
	}

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

	// 每隔 30 分钟维护一次内存
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for range ticker.C {
			yc.UpdateAll()
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

	// 获取最新的 m3u 并进行缓存
	return yc.CacheM3U(chId, formatCode)
}

// CacheM3U 调用 yt-dlp 获取频道的 m3u 地址
func (yc *youtubeCacher) CacheM3U(chId, formatCode string) (string, error) {
	yc.mu.Lock()
	defer yc.mu.Unlock()

	res, err := ytdlp.Extract(YoutubeResPrefix+chId, formatCode)
	if err != nil {
		return "", fmt.Errorf("调用 yt-dlp 失败: %v", err)
	}

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

// UpdateAll 更新缓存, 同时淘汰太长时间未读的列表
func (yc *youtubeCacher) UpdateAll() {
	log.Println(colors.ToBlue("youtube 缓存开始刷新..."))
	successCnt, failCnt, removeCnt := 0, 0, 0
	defer func() {
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
