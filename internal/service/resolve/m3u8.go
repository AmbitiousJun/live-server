package resolve

import (
	"fmt"
	"net/http"
)

// M3U8Result 根据处理器参数返回 m3u 地址的处理结果
func M3U8Result(url string, params HandleParams) (HandleResult, error) {
	// 如果无需代理, 直接重定向
	if !params.ProxyM3U {
		return HandleResult{Type: ResultRedirect, Url: url}, nil
	}

	content, err := ProxyM3U(url, params.Headers, params.ProxyTs, params.TsProxyMode, params.ClientHost)
	if err != nil {
		return HandleResult{}, fmt.Errorf("代理 m3u 失败: %v", err)
	}

	respHeader := make(http.Header)
	respHeader.Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	return HandleResult{
		Type:   ResultProxy,
		Code:   http.StatusOK,
		Body:   []byte(content),
		Header: respHeader,
	}, nil
}
