package resolve

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/m3u8"
	"github.com/AmbitiousJun/live-server/internal/service/net"
	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/gin-gonic/gin"
)

const (
	Env_CustomTsProxyEnableKey = "custom_ts_proxy_enable" // 是否启用自定义的代理接口
	Env_CustomTsProxyHostKey   = "custom_ts_proxy_host"   // 自定义代理接口地址
	DefaultProxyUA             = "libmpv"                 // 代理的默认客户端标识
)

var (

	// cacheableProxyClient 使用带缓存特性的 http 客户端代理 m3u
	cacheableProxyClient = https.NewCacheClient(1000, time.Second*5)

	// cacheableTsProxyClient 使用带缓存特性的 http 客户端代理 ts
	cacheableTsProxyClient = https.NewCacheClient(100, time.Second*10)
)

// ProxyM3U 代理 m3u 地址
//
// 代理成功时会返回代理后的 m3u 文本
func ProxyM3U(m3uLink string, header http.Header, proxyTs bool) (string, error) {
	// 设置默认的客户端标识
	if header == nil {
		header = make(http.Header)
	}
	if header.Get("User-Agent") == "" {
		header.Set("User-Agent", DefaultProxyUA)
	}

	// 请求远程
	finalLink, resp, err := cacheableProxyClient.Request(http.MethodGet, m3uLink, header, nil, true)
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
	m3uInfo, err := m3u8.ReadContent(m3u8.ExtractUrl(finalLink), string(bodyBytes))
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
	// 校验客户端 ip 是否可受信任
	clientIp := c.ClientIP()
	if net.IsBlackIp(clientIp) {
		c.String(http.StatusNotFound, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
		return
	}
	ipInfo, ok := net.GetIpAddrInfo(clientIp)
	if !ok || !whitearea.Passable(ipInfo) {
		c.String(http.StatusNotFound, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
		return
	}

	customEnable, ok := env.Get(Env_CustomTsProxyEnableKey)
	if ok && customEnable == "1" {
		customHost, ok := env.Get(Env_CustomTsProxyHostKey)
		if !ok {
			log.Printf(colors.ToRed("代理失败，请先设置自定义代理接口环境变量: %s"), Env_CustomTsProxyHostKey)
			c.String(http.StatusBadRequest, "代理失败，请检查日志")
			return
		}

		cu, err := url.Parse(customHost)
		if err != nil {
			log.Printf(colors.ToRed("代理失败，自定义代理接口无法解析: %s => %v"), customHost, err)
			c.String(http.StatusBadRequest, "代理失败，请检查日志")
			return
		}

		q := cu.Query()
		q.Set("remote", c.Query("remote"))
		cu.RawQuery = q.Encode()
		c.Redirect(http.StatusFound, cu.String())
		return
	}

	// 解码远程 url 地址
	remoteBytes, err := base64.StdEncoding.DecodeString(c.Query("remote"))
	if err != nil {
		log.Println(colors.ToRed("代理切片失败, 参数必须是 base64 编码"))
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	remote := string(remoteBytes)

	header := make(http.Header)
	header.Set("User-Agent", DefaultProxyUA)
	_, resp, err := cacheableTsProxyClient.Request(http.MethodGet, remote, header, nil, true)
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

	// 设置允许缓存
	resp.Header.Set("Cache-Control", "public, max-age=31536000")
	resp.Header.Del("Expires")
	resp.Header.Del("Last-Modified")
	resp.Header.Del("Date")
	resp.Header.Set("Content-Type", "text/html")
	c.Status(resp.StatusCode)
	https.CloneHeader(c, resp.Header)
	io.Copy(c.Writer, resp.Body)
}
