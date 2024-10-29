package resolve

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AmbitiousJun/live-server/internal/util/https"
)

func init() {
	registerHandler(new(thirdGdtvHandler))
}

// thirdGdtvHandler 第三方广东新闻处理器
type thirdGdtvHandler struct{}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (thirdGdtvHandler) Handle(params HandleParams) (HandleResult, error) {
	// 1 使用特定 UA 请求远程
	header := make(http.Header)
	header.Set("User-Agent", "libmpv")
	u := fmt.Sprintf("https://php.17186.eu.org/gdtv/web/%s.m3u8", params.ChName)
	resp, err := https.Request(http.MethodGet, u, header, nil)
	if err != nil {
		return HandleResult{}, fmt.Errorf("请求远程失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HandleResult{}, fmt.Errorf("请求远程出错, 响应码: %d", resp.StatusCode)
	}

	// 2 代理响应
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return HandleResult{}, fmt.Errorf("读取远程响应失败: %v", err)
	}
	return HandleResult{
		Type:   ResultProxy,
		Code:   resp.StatusCode,
		Header: resp.Header,
		Body:   bodyBytes,
	}, nil
}

// Name 处理器名称
func (thirdGdtvHandler) Name() string {
	return "third_gdtv"
}
