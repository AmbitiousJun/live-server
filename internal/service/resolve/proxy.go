package resolve

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/AmbitiousJun/live-server/internal/service/m3u8"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/gin-gonic/gin"
)

// ProxyM3U 代理 m3u 地址
//
// 代理成功时会返回代理后的 m3u 文本
func ProxyM3U(m3uLink string, header http.Header, proxyTs bool) (string, error) {
	// 请求远程
	finalLink, resp, err := https.Request(http.MethodGet, m3uLink, header, nil, true)
	if err != nil {
		return "", fmt.Errorf("请求远程地址失败: %s, err: %v", m3uLink, err)
	}
	defer resp.Body.Close()

	// 非成功响应, 有可能是被拦截
	if !https.IsSuccessCode(resp.StatusCode) {
		return "", fmt.Errorf("请求远程地址失败: %d %s", resp.StatusCode, resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败, err: %v", err)
	}

	// 解析 m3u
	m3uInfo, err := m3u8.ReadContent(m3u8.ExtractPrefix(finalLink), string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("解析 m3u 失败: %s, err: %v", m3uLink, err)
	}

	// 不代理切片, 直接返回原始文本
	if !proxyTs {
		return m3uInfo.Content(), nil
	}

	// 将 ts 切片地址更改为本地代理地址
	return m3uInfo.ContentFunc(func(tsIdx int, tsUrl string) string {
		u, _ := url.Parse("/proxy_ts")
		q := u.Query()
		q.Set("remote", tsUrl)
		u.RawQuery = q.Encode()
		return u.String()
	}), nil
}

// ProxyTs 代理 ts 切片
func ProxyTs(c *gin.Context) {
	remote := c.Query("remote")

	_, resp, err := https.Request(http.MethodGet, remote, nil, nil, true)
	if err != nil {
		log.Printf(colors.ToRed("代理切片失败: %v"), err)
		c.String(http.StatusInternalServerError, "代理切片失败")
		return
	}
	defer resp.Body.Close()

	// 非成功响应, 有可能是被拦截
	if !https.IsSuccessCode(resp.StatusCode) {
		log.Printf(colors.ToRed("请求远程地址失败: %d %s"), resp.StatusCode, resp.Status)
		c.String(http.StatusInternalServerError, "代理切片失败失败")
		return
	}

	c.Status(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}
	io.Copy(c.Writer, resp.Body)
}
