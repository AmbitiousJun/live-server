package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (
	// sttvLiveDetailUrl 获取直播源的接口地址
	sttvLiveDetailUrl = "https://sttv42-api.strtv.cn/api/getlivedetail.php?gid=${gid}&type=hdtv"

	// sttvTmpM3uPrefix 临时 m3u 的环境变量前缀
	sttvTmpM3uPrefix = "sttv_tmp_m3u"
)

var sttvChannels = map[string]string{
	"st1": "1169873", // 汕头 1 台
}

func init() {
	resolve.RegisterHandler(new(sttvHandler))
}

// sttvHandler 汕头橄榄台处理器
type sttvHandler struct{}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *sttvHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 验证频道
	gid, ok := sttvChannels[params.ChName]
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("不支持的频道: %s", params.ChName)
	}

	// 拼接路径, 发起请求
	u := strings.ReplaceAll(sttvLiveDetailUrl, "${gid}", gid)
	resp, err := https.Get(u).Do()
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析响应
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("解析响应失败: %v", err)
	}

	// 接口被拦截, 返回环境变量中临时的 m3u 地址
	if resp.StatusCode == http.StatusForbidden {
		m3u, ok := env.Get(sttvTmpM3uPrefix + "_" + params.ChName)
		if ok {
			return resolve.HandleResult{Type: resolve.ResultRedirect, Url: m3u}, nil
		}
	}

	if resp.StatusCode != http.StatusOK {
		return resolve.HandleResult{}, fmt.Errorf("错误的响应码: %d, 原始响应: %s", resp.StatusCode, string(bodyBytes))
	}

	resJson, err := jsons.New(string(bodyBytes))
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("JSON 转换失败: %v, 原始响应: %s", err, string(bodyBytes))
	}

	liveUrl, ok := resJson.Attr("data").Attr("liveurl").String()
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("获取直播源失败, 原始响应: %s", string(bodyBytes))
	}

	// 将 m3u 地址缓存到环境变量中
	env.Set(sttvTmpM3uPrefix+"_"+params.ChName, liveUrl)
	return resolve.HandleResult{
		Type: resolve.ResultRedirect,
		Url:  liveUrl,
	}, nil
}

// Name 处理器名称
func (h *sttvHandler) Name() string {
	return "sttv"
}

// HelpDoc 处理器说明文档
func (h *sttvHandler) HelpDoc() string {
	return "\n目前已失效，勿用"
}

// SupportProxy 是否支持 m3u 代理
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *sttvHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *sttvHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (h *sttvHandler) Enabled() bool {
	return false
}
