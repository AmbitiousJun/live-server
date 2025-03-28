package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/base64s"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

// yangHandler 处理 Yang1989 M3U
type yangHandler struct {
	resolve.CommonM3U8

	// serverHost yang 服务器地址
	serverHost string

	// m3uAddr m3u 订阅地址
	m3uAddr string

	// reqHeaders 请求需要携带的请求头
	reqHeaders http.Header

	// errorRespSeg 远程解析异常的相应片段
	errorRespSeg string

	// subC 缓存订阅数据的 http 客户端
	subC *https.CacheClient

	// chC 缓存频道数据的 http 客户端
	chC *https.CacheClient
}

func init() {
	y := new(yangHandler)
	y.serverHost = base64s.MustDecodeString("dHYuaWlsbC50b3A=")
	y.m3uAddr = fmt.Sprintf("https://%s/m3u/Gather", y.serverHost)
	y.reqHeaders = make(http.Header)
	y.reqHeaders.Set("User-Agent", "okhttp")
	y.errorRespSeg = "404"
	y.subC = https.NewCacheClient(1, time.Hour)
	y.chC = https.NewCacheClient(50, time.Minute*10)
	resolve.RegisterHandler(y)
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (y *yangHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 解析订阅信息
	infos, err := y.ResolveSub(y.subC, y.m3uAddr, y.reqHeaders)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("解析订阅地址失败: %v", err)
	}

	// 获取用户请求的频道
	resInfo, err := y.MatchChannel(infos, params.ChName, params.Format)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("匹配频道失败: %v", err)
	}

	// 手动请求一次频道地址, 重定向到真实地址并进行缓存
	finalUrl, resp, err := y.chC.Request(http.MethodGet, resInfo.Url, y.reqHeaders, nil, true)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求频道地址 %s 失败: %v", resInfo.Url, err)
	}
	defer resp.Body.Close()
	if !https.IsSuccessCode(resp.StatusCode) {
		return resolve.HandleResult{}, fmt.Errorf("请求频道地址 %s 失败: %s", resInfo.Url, resp.Status)
	}

	if strings.Contains(finalUrl, y.errorRespSeg) {
		return resolve.HandleResult{}, fmt.Errorf("频道丢失: %s, 请等待官方修复", params.ChName)
	}

	return resolve.HandleResult{
		Type: resolve.ResultRedirect,
		Url:  finalUrl,
	}, nil
}

// Name 处理器名称
func (y *yangHandler) Name() string {
	return "yang"
}

// HelpDoc 处理器说明文档
func (y *yangHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 解除 YanG-1989 大佬订阅的 UA 限制, 仅限个人部署方便调用, 请勿传播")
	sb.WriteString("\n2. 调用示例: ${clientOrigin}/handler/yang/ch/CCTV1")
	sb.WriteString("\n3. 具体可用的频道名称列表请从 <a href=\"https://github.com/YanG-1989/m3u\" target=\"_blank\">官方订阅</a> 中进行查看")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (y *yangHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (y *yangHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (y *yangHandler) Enabled() bool {
	return true
}
