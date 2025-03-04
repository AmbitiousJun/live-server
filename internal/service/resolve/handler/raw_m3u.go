package handler

import (
	"fmt"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
)

func init() {
	resolve.RegisterHandler(&rawM3UHandler{})
}

// rawM3UHandler 对 m3u 播放地址直接进行处理
type rawM3UHandler struct{}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *rawM3UHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 获取环境变量
	url, ok := env.Get(params.UrlEnv)
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("获取不到环境变量: [%s]", params.UrlEnv)
	}
	params.Headers = nil
	return resolve.M3U8Result(url, params)
}

// Name 处理器名称
func (h *rawM3UHandler) Name() string {
	return "raw_m3u"
}

// HelpDoc 处理器说明文档
func (h *rawM3UHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 将 m3u 直播地址设置到环境变量中, 变量名随意")
	sb.WriteString("\n2. 类似于 remote_m3u 处理器的形式调用即可")
	sb.WriteString("\n3. 该处理器适用于在本地服务直接代理远程的 m3u 播放地址（非订阅地址）")
	sb.WriteString("\n4. 使用该处理器时, 频道名称 ch 可任意传递一个非空值")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *rawM3UHandler) SupportM3UProxy() bool {
	return true
}

// Enabled 标记处理器是否是启用状态
func (h *rawM3UHandler) Enabled() bool {
	return true
}
