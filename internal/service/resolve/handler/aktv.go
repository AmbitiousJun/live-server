package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/base64s"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

func init() {
	h := &aktvHandler{
		host:       base64s.MustDecodeString("aHR0cHM6Ly9ha3R2LnNwYWNl"),
		errorTsSeg: "error",
		cc:         https.NewCacheClient(1, time.Hour*2),
	}
	h.m3uAddr = h.host + "/live.m3u"
	resolve.RegisterHandler(h)
}

// aktvHandler 处理 AKTV 直播源
type aktvHandler struct {
	resolve.CommonM3U8

	// host 主站地址
	host string

	// m3uAddr m3u 订阅地址
	m3uAddr string

	// errorTsSeg 频道异常的 ts 切片地址片段
	errorTsSeg string

	// cc 缓存订阅信息的 http 客户端
	cc *https.CacheClient
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *aktvHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 1 获取订阅信息
	infos, err := h.ResolveSub(h.cc, h.m3uAddr, nil)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("获取订阅信息失败: %v", err)
	}

	// 2 匹配频道
	chInfo, err := h.MatchChannel(infos, params.ChName, params.Format)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("匹配频道失败: %v, 可用频道: [%s]", err, strings.Join(h.ChannelSlice(infos), ", "))
	}

	// 3 代理原始的 m3u 文本, 检查频道是否失效
	content, err := resolve.ProxyM3U(chInfo.Url, nil, false, resolve.ModeCustom, params.ClientHost)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("检查频道有效性失败: %v", err)
	}
	if strings.Contains(content, h.errorTsSeg) {
		return resolve.HandleResult{}, fmt.Errorf("频道异常: [%s], 请稍后再试", params.ChName)
	}

	return resolve.M3U8Result(chInfo.Url, params)
}

// Name 处理器名称
func (h *aktvHandler) Name() string {
	return "aktv"
}

// HelpDoc 处理器说明文档
func (h *aktvHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 处理 AKTV 源，同时检查源是否失效")
	sb.WriteString("\n2. 请求示例：${clientOrigin}/handler/aktv/ch/翡翠台?proxy_m3u=1&proxy_ts=1&format=1")
	sb.WriteString("\n3. 可用频道列表请查看 <a href=\"" + h.host + "\" target=\"_blank\">官方订阅</a>")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *aktvHandler) SupportM3UProxy() bool {
	return true
}

// SupportCustomHeaders 是否支持自定义请求头
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *aktvHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (h *aktvHandler) Enabled() bool {
	return true
}
