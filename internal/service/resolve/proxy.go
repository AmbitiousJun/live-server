package resolve

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	HeadersSeg                 = "[[[:]]]"                // 分割请求头的分隔符
)

var (

	// cacheableProxyClient 使用带缓存特性的 http 客户端代理 m3u
	cacheableProxyClient = https.NewCacheClient(1000, time.Second*5)

	// cacheableTsProxyClient 使用带缓存特性的 http 客户端代理 ts
	cacheableTsProxyClient = https.NewCacheClient(50, time.Second*30)
)

// ProxyM3U 代理 m3u 地址
//
// 代理成功时会返回代理后的 m3u 文本
func ProxyM3U(m3uLink string, header http.Header, proxyTs bool, tsProxyMode TsProxyMode, clientHost string) (string, error) {
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

	basePath := clientHost + "/proxy.ts"
	if customHost, ok := getCustomTsProxyHost(tsProxyMode); ok {
		basePath = customHost
	}

	var headerStr string
	if header != nil {
		kvs := []string{}
		for k, vs := range header {
			kvs = append(kvs, k, strings.Join(vs, ", "))
		}
		headerStr = base64.StdEncoding.EncodeToString([]byte(strings.Join(kvs, HeadersSeg)))
	}

	tsLink, _ := url.Parse(basePath)

	// 将 ts 切片地址更改为本地代理地址
	return m3uInfo.ContentFunc(func(tsIdx int, tsUrl string) string {
		remoteStr := base64.StdEncoding.EncodeToString([]byte(tsUrl))
		q := tsLink.Query()
		q.Set("remote", remoteStr)
		if headerStr != "" {
			q.Set("headers", headerStr)
		}

		tsLink.RawQuery = q.Encode()
		return tsLink.String()
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

	// 解码远程 url 地址
	remoteBytes, err := base64.StdEncoding.DecodeString(c.Query("remote"))
	if err != nil {
		log.Println(colors.ToRed("代理切片失败, 参数 [remote] 必须是 base64 编码"))
		c.String(http.StatusBadRequest, "参数错误")
		return
	}
	remote := string(remoteBytes)

	header := make(http.Header)
	header.Set("User-Agent", DefaultProxyUA)
	headerBytes, err := base64.StdEncoding.DecodeString(c.Query("headers"))
	if err != nil {
		log.Println(colors.ToRed("代理切片失败, 参数 [headers] 必须是 base64 编码"))
		c.String(http.StatusBadRequest, "参数错误")
		return
	}

	headerStr := string(headerBytes)
	if headerStr != "" {
		kvs := strings.Split(headerStr, HeadersSeg)
		for i := 0; i+1 < len(kvs); i += 2 {
			header.Set(kvs[i], kvs[i+1])
		}
	}

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
	c.Header("Cache-Control", "max-age=3600")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(resp.StatusCode)
	io.Copy(c.Writer, resp.Body)
}

// getCustomTsProxyHost 获取用户配置的自定义切片代理地址, 如果未开启自定义代理, 则返回 false
func getCustomTsProxyHost(proxyMode TsProxyMode) (string, bool) {
	if proxyMode == ModeLocal {
		return "", false
	}

	customHost, ok := env.Get(Env_CustomTsProxyHostKey)
	if proxyMode == ModeCustom {
		return customHost, true
	}
	if !ok {
		return "", false
	}

	customEnable, ok := env.Get(Env_CustomTsProxyEnableKey)
	if !ok || customEnable != "1" {
		return "", false
	}
	return customHost, true
}
