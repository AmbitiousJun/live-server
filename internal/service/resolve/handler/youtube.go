package handler

import (
	"fmt"
	"log"
	"strings"
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

// youtubeParams 频道请求参数
type youtubeParams struct {
	chId       string // 频道 id
	formatCode string // 频道格式
}

// youtubeHandler Youtube 直播处理器
type youtubeHandler struct {
	cacher *resolve.Cacher[youtubeParams]
}

func init() {
	y := new(youtubeHandler)
	y.initCacher()
	resolve.RegisterHandler(y)
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (y *youtubeHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	formatCode := y.chooseFormat(params.Format)
	playlist, err := y.cacher.Request(youtubeParams{chId: params.ChName, formatCode: formatCode})
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("获取 playlist 失败: %v", err)
	}

	return resolve.M3U8Result(playlist, params)
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

// initCacher 初始化频道缓存器
func (y *youtubeHandler) initCacher() {
	y.cacher = resolve.NewCacher(

		resolve.WithCalcCacheKey(func(p youtubeParams) string {
			return p.chId + ":" + p.formatCode
		}),

		resolve.WithRecoverCacheKey(func(cacheKey string) (youtubeParams, bool) {
			idx := strings.LastIndex(cacheKey, ":")
			if idx == -1 {
				return youtubeParams{}, false
			}
			return youtubeParams{chId: cacheKey[:idx], formatCode: cacheKey[idx+1:]}, true
		}),

		resolve.WithFetchValue(func(p youtubeParams) (string, error) {
			res, err := ytdlp.Extract(YoutubeResPrefix+p.chId, p.formatCode)
			if err != nil {
				return "", fmt.Errorf("调用 yt-dlp 失败: %v", err)
			}
			return res, nil
		}),

		resolve.WithUpdateComplete[youtubeParams](func(success, fail, remove int) {
			log.Printf(colors.ToGreen("youtube 缓存刷新完成, 成功: %d, 失败: %d, 移除: %d"), success, fail, remove)
		}),

		resolve.WithCacheTimeout[youtubeParams](time.Hour),
		resolve.WithRemoveInterval[youtubeParams](time.Minute*10),
		resolve.WithUpdateInterval[youtubeParams](time.Hour+time.Minute*30),
	)
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
