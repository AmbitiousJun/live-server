package resolve

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/AmbitiousJun/live-server/internal/service/m3u8"
	"github.com/AmbitiousJun/live-server/internal/service/net"
	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
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
		return "", fmt.Errorf("请求远程地址失败: %s", resp.Status)
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
		q.Set("remote", base64.StdEncoding.EncodeToString([]byte(tsUrl)))
		u.RawQuery = q.Encode()
		return u.String()
	}), nil
}

// ProxyTs 代理 ts 切片
func ProxyTs(c *gin.Context) {
	remoteBytes, err := base64.StdEncoding.DecodeString(c.Query("remote"))
	if err != nil {
		log.Println(colors.ToRed("代理切片失败, 参数必须是 base64 编码"))
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	remote := string(remoteBytes)

	// 校验客户端 ip 是否可受信任
	clientIp := c.ClientIP()
	if net.IsBlackIp(clientIp) {
		c.String(http.StatusForbidden, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
		return
	}
	ipInfo, ok := net.GetIpAddrInfo(clientIp)
	if !ok || !whitearea.Passable(ipInfo) {
		c.String(http.StatusForbidden, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
		return
	}

	_, resp, err := https.Request(http.MethodGet, remote, nil, nil, true)
	if err != nil {
		log.Printf(colors.ToRed("代理切片失败: %v"), err)
		c.String(http.StatusInternalServerError, "代理切片失败")
		return
	}
	defer resp.Body.Close()

	// 非成功响应, 有可能是被拦截
	if !https.IsSuccessCode(resp.StatusCode) {
		log.Printf(colors.ToRed("请求远程地址失败: %s"), resp.Status)
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
