package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

func init() {
	resolve.RegisterHandler(new(thirdGdtvHandler))
}

// thirdGdtvHandler 第三方广东新闻处理器
type thirdGdtvHandler struct{}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (thirdGdtvHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 1 使用特定 UA 请求远程
	header := make(http.Header)
	header.Set("User-Agent", "libmpv")
	u := fmt.Sprintf("https://php.17186.eu.org/gdtv/web/%s.m3u8", params.ChName)
	resp, err := https.Request(http.MethodGet, u, header, nil)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求远程失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resolve.HandleResult{}, fmt.Errorf("请求远程出错, 响应码: %d", resp.StatusCode)
	}

	// 2 代理响应
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("读取远程响应失败: %v", err)
	}
	return resolve.HandleResult{
		Type:   resolve.ResultProxy,
		Code:   resp.StatusCode,
		Header: resp.Header,
		Body:   bodyBytes,
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
