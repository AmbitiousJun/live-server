package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/warp"
	"github.com/AmbitiousJun/live-server/internal/service/ytdlp"
	"github.com/AmbitiousJun/live-server/internal/util/base64s"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"golang.org/x/sync/errgroup"
)

// youtubeParams 频道请求参数
type youtubeParams struct {
	chId       string // 频道 id
	formatCode string // 频道格式
}

// youtubeHandler Youtube 直播处理器
type youtubeHandler struct {

	// customFormatEnableEnv 解析自定义格式开关的环境变量名
	customFormatEnableEnv string

	// chPrefix 频道地址前缀
	chPrefix string

	// cacher 解析并缓存频道地址的处理器
	cacher *resolve.Cacher[youtubeParams]

	// warpListenerId warp 监听器 id, 用于动态检测修复被封禁 ip
	warpListenerId string
}

func init() {

	y := &youtubeHandler{
		customFormatEnableEnv: "youtube_custom_format_enable",
		chPrefix:              base64s.MustDecodeString("aHR0cHM6Ly93d3cueW91dHViZS5jb20vd2F0Y2g/dj0="),
	}
	y.initCacher()

	l := warp.Listener{
		CheckIP: func() error {
			_, audio := y.chooseFormat("")
			_, err := ytdlp.Extract(y.chPrefix+"6IquAgfvYmc", audio)
			if err != nil && strings.Contains(err.Error(), "Sign in to confirm you’re not a bot") {
				return fmt.Errorf("youtube 不可用: %v", err)
			}
			return nil
		},
	}
	y.warpListenerId = warp.RegisterListener(l)

	resolve.RegisterHandler(y)
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (y *youtubeHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	videoCode, audioCode := y.chooseFormat(params.Format)

	var eg errgroup.Group
	var err error
	var videoPlaylist, audioPlaylist string

	eg.Go(func() error {
		videoPlaylist, err = y.cacher.Request(youtubeParams{chId: params.ChName, formatCode: videoCode})
		if err != nil {
			return fmt.Errorf("获取 video playlist 失败: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		audioPlaylist, err = y.cacher.Request(youtubeParams{chId: params.ChName, formatCode: audioCode})
		if err != nil {
			return fmt.Errorf("获取 audio playlist 失败: %w", err)
		}
		return nil
	})

	if err = eg.Wait(); err != nil {
		go warp.ReportError(y.warpListenerId)
		return resolve.HandleResult{}, nil
	}

	// 不支持重定向, 必须启用 proxy_m3u 参数
	if !params.ProxyM3U {
		return resolve.HandleResult{}, errors.New("该解析器暂时不支持直接重定向, 请启用 proxy_m3u 参数")
	}

	// 拼接请求头
	headers := []string{
		"User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36",
		"Referer", "https://www.youtube.com/",
		"Origin", "https://www.youtube.com",
	}

	// 构造 query 参数
	urlQueries := url.Values{}
	urlQueries.Set("headers", base64.StdEncoding.EncodeToString([]byte(strings.Join(headers, constant.HeadersSeg))))
	urlQueries.Set("client_host", base64.StdEncoding.EncodeToString([]byte(params.ClientHost)))
	if params.ProxyTs {
		urlQueries.Set("proxy_ts", "1")
	}
	switch params.TsProxyMode {
	case resolve.ModeCustom:
		urlQueries.Set("ts_proxy_mode", string(params.TsProxyMode))
	default:
		urlQueries.Set("ts_proxy_mode", string(resolve.ModeLocal))
	}

	// 构造视频 m3u 代理地址
	urlQueries.Set("remote", base64.StdEncoding.EncodeToString([]byte(videoPlaylist)))
	videoUrl, _ := url.Parse("/proxy.m3u8")
	videoUrl.RawQuery = urlQueries.Encode()

	// 构造音频 m3u 代理地址
	urlQueries.Set("remote", base64.StdEncoding.EncodeToString([]byte(audioPlaylist)))
	audioUrl, _ := url.Parse("/proxy.m3u8")
	audioUrl.RawQuery = urlQueries.Encode()

	// 构造响应
	responseHeader := make(http.Header)
	responseHeader.Set("Content-Type", "application/vnd.apple.mpegurl")

	responseBody := fmt.Sprintf(`#EXTM3U

#EXT-X-MEDIA:TYPE=AUDIO,GROUP-ID="audio",NAME="Audio",DEFAULT=YES,AUTOSELECT=YES,URI="%s"

#EXT-X-STREAM-INF:BANDWIDTH=3000000,AUDIO="audio"
%s
`, audioUrl.String(), videoUrl.String())

	return resolve.HandleResult{
		Type:   resolve.ResultProxy,
		Code:   http.StatusOK,
		Header: responseHeader,
		Body:   []byte(responseBody),
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
	sb.WriteString("\n7. 解析失败时，请自行在浏览器登录 youtube 账号，并使用 <a href=\"https://chromewebstore.google.com/detail/editthiscookie/cmbkolgnkghmgajbbapoicfhjlabmpef?utm_source=ext_app_menu\" target=\"__blank\">EditThisCookie</a> 插件导出 Netscape 格式的 cookie 文件放到程序读取目录下解决")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (y *youtubeHandler) SupportM3UProxy() bool {
	return true
}

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (y *youtubeHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (y *youtubeHandler) Enabled() bool {
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
			res, err := ytdlp.Extract(y.chPrefix+p.chId, p.formatCode)
			if err != nil {
				return "", fmt.Errorf("调用 yt-dlp 失败: %v", err)
			}
			return res, nil
		}),

		resolve.WithUpdateComplete[youtubeParams](func(success, fail, remove int) {
			log.Printf(colors.ToGreen("youtube 缓存刷新完成, 成功: %d, 失败: %d, 移除: %d"), success, fail, remove)
		}),

		resolve.WithCacheTimeout[youtubeParams](time.Hour*3),
		resolve.WithRemoveInterval[youtubeParams](time.Minute*10),
		resolve.WithUpdateInterval[youtubeParams](time.Hour*4),
	)
}

// chooseFormat 选择一个合适的 yt-dlp formatCode
func (y *youtubeHandler) chooseFormat(wantFmt string) (video, audio string) {
	// 默认使用 720 p
	video, audio = "232", "234"

	// 判断是否允许自定义
	enable, ok := env.Get(y.customFormatEnableEnv)
	if !ok || enable != "1" {
		return
	}

	switch wantFmt {
	case "FHD":
		return "270", audio
	}

	return
}
