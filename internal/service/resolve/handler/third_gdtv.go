package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
)

func init() {
	resolve.RegisterHandler(new(thirdGdtvHandler))
}

// thirdGdtvHandler 第三方广东新闻处理器
type thirdGdtvHandler struct{}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (thirdGdtvHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	u := fmt.Sprintf("https://php.17186.eu.org/gdtv/web/%s.m3u8", params.ChName)

	// 无需代理
	if !params.ProxyM3U {
		return resolve.HandleResult{Type: resolve.ResultRedirect, Url: u}, nil
	}

	// 使用特定 UA 请求远程
	header := make(http.Header)
	header.Set("User-Agent", "libmpv")
	content, err := resolve.ProxyM3U(u, header, params.ProxyTs)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("代理 m3u 失败: %v", err)
	}

	respHeader := make(http.Header)
	header.Set("Content-Type", "application/vnd.apple.mpegurl")
	return resolve.HandleResult{
		Type:   resolve.ResultProxy,
		Code:   http.StatusOK,
		Header: respHeader,
		Body:   []byte(content),
	}, nil
}

// Name 处理器名称
func (thirdGdtvHandler) Name() string {
	return "third_gdtv"
}

// HelpDoc 处理器说明文档
func (thirdGdtvHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 第三方的荔枝网代理接口，原地址：https://php.17186.eu.org/gdtv/web/{频道名称}.m3u8")
	sb.WriteString("\n2. 该处理器不保证可用性，只是简单地绕过 UA 限制进行代理，部分 ip 会被第三方接口屏蔽无法使用")
	sb.WriteString("\n3. 可尝试频道：xwpd(广东新闻)、gdws(广东卫视)、gdzj(广东珠江)、gdgg(广东民生)、gdty(广东体育)、nfws(大湾区卫视)、gdgj(大湾区卫视海外版)、jjkj(经济科教)、gdzy(4k超高清)、gdse(广东少儿)、jjkt(嘉佳卡通)、nfgw(南方购物)、lnxq(岭南戏曲)、xdjy(现代教育)、ydpd(广东移动)")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
// 如果返回 true, 会自动在帮助文档中加入标记
func (thirdGdtvHandler) SupportM3UProxy() bool {
	return true
}
