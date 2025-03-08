package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/subm3u"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
)

func init() {
	resolve.RegisterHandler(&remoteM3UHandler{
		cc: https.NewCacheClient(1000, time.Hour),
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
	destInfos, ok := infos[params.ChName]
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("匹配频道名称失败: %s, 请检查远程地址是否有效", params.ChName)
	}

	// 根据 format 参数筛选出目标直播源
	resInfo := destInfos[0]
	formatIdx, err := strconv.Atoi(params.Format)
	if strs.AllNotEmpty(params.Format) && err != nil {
		log.Printf(colors.ToYellow("format 索引解析失败: %v, 有效范围: [1, %d], 本次请求使用默认值"), err, len(destInfos))
	}
	if err == nil {
		formatIdx--
		if formatIdx < 0 || formatIdx >= len(destInfos) {
			return resolve.HandleResult{}, fmt.Errorf("format 索引传递错误, 有效范围: [1, %d]", len(destInfos))
		}
		resInfo = destInfos[formatIdx]
	}

	return resolve.M3U8Result(resInfo.Url, params)
}

// HelpDoc 处理器说明文档
func (h *remoteM3UHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 将有效的 m3u 在线地址设置到程序的环境变量中，变量名随意，如：remote_m3u_v6")
	sb.WriteString("\n2. 调用处理器时，传递有效的频道名称和环境变量名，即可观看")
	sb.WriteString("\n3. 环境变量名传递方式：在调用地址后边加上 query 参数，如：${clientOrigin}/handler/remote_m3u/ch/CCTV1?url_env=remote_m3u_v6")
	sb.WriteString("\n4. 频道名传递方式：")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;① 假设有如下 m3u 的直播源 A 信息: #EXTINF:-1 tvg-name=\"寰宇新闻1\" tvg-id=\"寰宇新闻2\" group-title=\"港澳台及国外\",寰宇新闻3")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;② 传递频道名称 ch 参数时, 使用 [寰宇新闻1, 寰宇新闻2, 寰宇新闻3] 中的任意一个都能够匹配到当前频道源")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;③ 假设除了直播源 A 之外还有直播源 B 信息: #EXTINF:-1 tvg-name=\"寰宇新闻4\" tvg-id=\"寰宇新闻3\" group-title=\"港澳台及国外\",寰宇新闻2")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;④ A 和 B 中有部分信息是重叠的, 这时可以通过 format 参数来指定具体的直播源")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;⑤ 比如: ${clientOrigin}/handler/remote_m3u/ch/寰宇新闻2?url_env=remote_m3u_v6&format=1 可以指定使用直播源 A")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;⑥ 同理, 传递 format=2 可以指定使用直播源 B")
	sb.WriteString("\n&nbsp;&nbsp;&nbsp;⑦ 如果出现信息重叠且没有传递 format 参数, 默认使用第一个直播源")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *remoteM3UHandler) SupportM3UProxy() bool {
	return true
}

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *remoteM3UHandler) SupportCustomHeaders() bool {
	return true
}

// Enabled 标记处理器是否是启用状态
func (h *remoteM3UHandler) Enabled() bool {
	return true
}
