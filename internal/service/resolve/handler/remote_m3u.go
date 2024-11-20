package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/subm3u"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

func init() {
	resolve.RegisterHandler(&remoteM3UHandler{
		cc: https.NewCacheClient(1000, time.Minute*10),
	})
}

// remoteM3UHandler 远程 m3u8 直播源处理器
type remoteM3UHandler struct {
	cc *https.CacheClient
}

func (h *remoteM3UHandler) Name() string {
	return "remote_m3u"
}

func (h *remoteM3UHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 获取环境变量
	url, ok := env.Get(params.UrlEnv)
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("获取不到环境变量: %s", params.UrlEnv)
	}

	// 请求远程 m3u 文本
	_, resp, err := h.cc.Request(http.MethodGet, url, nil, nil, true)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求远程地址失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析文本内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("读取远程响应失败: %v", err)
	}
	infos, err := subm3u.ReadContent(string(bodyBytes))
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("解析远程响应失败: %v, 原始响应: %s", err, string(bodyBytes))
	}

	// 获取用户请求的频道
	destInfo, ok := infos[params.ChName]
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("匹配频道名称失败: %s, 请检查远程地址是否有效", params.ChName)
	}

	// 如果无需代理, 直接重定向
	if !params.ProxyM3U {
		return resolve.HandleResult{Type: resolve.ResultRedirect, Url: destInfo.Url}, nil
	}

	content, err := resolve.ProxyM3U(destInfo.Url, nil, params.ProxyTs)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("代理 m3u 失败: %v", err)
	}

	respHeader := make(http.Header)
	respHeader.Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	respHeader.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	respHeader.Set("Pragma", "no-cache")
	respHeader.Set("Expires", "0")
	return resolve.HandleResult{
		Type:   resolve.ResultProxy,
		Code:   http.StatusOK,
		Body:   []byte(content),
		Header: respHeader,
	}, nil
}

// HelpDoc 处理器说明文档
func (h *remoteM3UHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 将有效的 m3u 在线地址设置到程序的环境变量中，变量名随意，如：remote_m3u_v6")
	sb.WriteString("\n2. 调用处理器时，传递有效的频道名称和环境变量名，即可观看")
	sb.WriteString("\n3. 环境变量名传递方式：在调用地址后边加上 query 参数，如：${clientOrigin}/handler/remote_m3u/ch/CCTV1?url_env=remote_m3u_v6")
	sb.WriteString("\n4. 频道名传递方式：程序会按照 [tvg-name, 后缀别名, tvg-id] 的顺序依次读取，以首个不为空的参数作为频道名称")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *remoteM3UHandler) SupportM3UProxy() bool {
	return true
}
