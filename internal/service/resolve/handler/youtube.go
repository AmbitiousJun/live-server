package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/ytdlp"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
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
type youtubeHandler struct{}

func init() {
	resolve.RegisterHandler(new(youtubeHandler))
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (y youtubeHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 判断是否有缓存的 playlist
	playlist, ok := env.Get(y.playlistCacheKey(params.ChName, params.Format))
	if ok && !y.isPlaylistValid(playlist) {
		playlist = ""
	}

	// 命中缓存失败, 调用 yt-dlp 进行解析
	if playlist == "" {
		formatCode := y.chooseFormat(params.Format)
		res, err := ytdlp.Extract(YoutubeResPrefix+params.ChName, formatCode)
		if err != nil {
			return resolve.HandleResult{}, fmt.Errorf("解析频道失败: %s, err: %v", params.ChName, err)
		}
		playlist = res
		key := y.playlistCacheKey(params.ChName, params.Format)
		env.Set(key, playlist)
		y.autoRemovePlaylistCache(key)
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
func (youtubeHandler) Name() string {
	return "youtube"
}

// HelpDoc 处理器说明文档
func (youtubeHandler) HelpDoc() string {
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
func (youtubeHandler) SupportM3UProxy() bool {
	return true
}

// isPlaylistValid 检查播放列表是否有效
func (youtubeHandler) isPlaylistValid(playlist string) bool {
	_, resp, err := https.Request(http.MethodGet, playlist, nil, nil, true)
	if err != nil {
		log.Printf(colors.ToRed("验证 playlist 有效性失败: %v"), err)
		return false
	}
	defer resp.Body.Close()
	return https.IsSuccessCode(resp.StatusCode)
}

// chooseFormat 选择一个合适的 yt-dlp formatCode
func (youtubeHandler) chooseFormat(wantFmt string) string {
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

// playlistCacheKey 获取播放列表在环境变量中的缓存键值
//
// format 参数只有在开启了自定义 format 功能之后才会被拼接到 key 中
func (youtubeHandler) playlistCacheKey(chName, format string) string {
	key := Env_YoutubePlaylistCache + ":" + chName

	enable, ok := env.Get(Env_YoutubeCustomFormatEnable)
	if !ok || enable != "1" {
		return key
	}

	return key + ":" + format
}

// autoRemovePlaylistCache 30 分钟后自动删除播放列表缓存
func (youtubeHandler) autoRemovePlaylistCache(key string) {
	env.SetAutoRefresh(key, func(curVal string) (string, error) {
		return "", env.ErrRemoveAndStop
	}, time.Minute*30)
}
