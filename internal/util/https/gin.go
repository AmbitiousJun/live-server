package https

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractReqBody 克隆并提取请求体
// 不影响 c 对象之后再次读取请求体
func ExtractReqBody(c *gin.Context) ([]byte, error) {
	if c == nil {
		return nil, nil
	}
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// ClientRequestHost 获取客户端请求的 Host
func ClientRequestHost(c *gin.Context) string {
	if c == nil {
		return ""
	}

	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

// ClientRequestUrl 获取客户端请求的完整地址
func ClientRequestUrl(c *gin.Context) string {
	return fmt.Sprintf("%s%s", ClientRequestHost(c), c.Request.URL.String())
}

// IsErrorResponse 判断一个请求响应是否是错误响应
//
// 判断标准是响应码以 4xx 5xx 开头
func IsErrorResponse(c *gin.Context) bool {
	if c == nil {
		return true
	}
	code := c.Writer.Status()
	str := strconv.Itoa(code)
	return strings.HasPrefix(str, "4") || strings.HasPrefix(str, "5")
}

// CloneHeader 克隆 http 头部到 gin 的响应头中
func CloneHeader(c *gin.Context, header http.Header) {
	if c == nil || header == nil {
		return
	}
	for key, values := range header {
		c.Writer.Header().Del(key)
		for _, value := range values {
			c.Header(key, value)
		}
	}
}

// ProxyRequest 代理请求, 返回远程响应
func ProxyRequest(c *gin.Context, remote string) (*http.Response, error) {
	if c == nil || remote == "" {
		return nil, errors.New("参数为空")
	}

	// 1 解析远程地址
	rmtUrl, err := url.Parse(remote + c.Request.URL.String())
	if err != nil {
		return nil, fmt.Errorf("解析远程地址失败: %v", err)
	}

	// 2 发送请求
	return Request(c.Request.Method, rmtUrl.String()).
		Header(c.Request.Header).
		Body(c.Request.Body).
		Do()
}

// ProxyPass 代理转发请求
func ProxyPass(c *gin.Context, remote string) error {
	if c == nil || remote == "" {
		return errors.New("参数为空")
	}

	// 1 代理请求
	resp, err := ProxyRequest(c, remote)
	if err != nil {
		return fmt.Errorf("代理请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 2 回写响应头
	c.Status(resp.StatusCode)
	CloneHeader(c, resp.Header)

	// 3 回写响应体
	_, err = io.Copy(c.Writer, resp.Body)
	return err
}
