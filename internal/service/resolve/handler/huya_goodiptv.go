package handler

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

func init() {
	resolve.RegisterHandler(&huyaGoodIptvHandler{
		jsonApi: "https://www.goodiptv.club/huya/${id}?type=json",
		cc:      https.NewCacheClient(1000, time.Minute*10),
	})
}

// huyaGoodIptvHandler 基于 goodiptv 接口的虎牙 flv 直播处理器
type huyaGoodIptvHandler struct {

	// jsonApi 获取 json 数据的 api 地址
	jsonApi string

	// cc 缓存请求数据
	cc *https.CacheClient
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *huyaGoodIptvHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 请求 json
	roomRequestUrl := strings.ReplaceAll(h.jsonApi, "${id}", params.ChName)
	_, resp, err := h.cc.Request(http.MethodGet, roomRequestUrl, nil, nil, true)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求远程地址失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析 json
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("读取响应失败: %v", err)
	}
	body := string(bodyBytes)
	resJson, err := jsons.New(body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, body)
	}

	// 从 map 中随机取出一个 flv 地址进行返回
	flvs, ok := resJson.Attr("flv").Done()
	if !ok || flvs.Type() != jsons.JsonTypeObj || flvs.Empty() {
		return resolve.HandleResult{}, fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, body)
	}
	var result string
	flvs.RangeObj(func(key string, value *jsons.Item) error {
		if str, ok := value.Ti().String(); ok && str != "" {
			result = str
			return jsons.ErrBreakRange
		}
		return nil
	})
	if result == "" {
		return resolve.HandleResult{}, fmt.Errorf("解析响应失败: %v, 原始响应: %s", err, body)
	}

	return resolve.HandleResult{
		Type: resolve.ResultRedirect,
		Url:  result,
	}, nil
}

// Name 处理器名称
func (h *huyaGoodIptvHandler) Name() string {
	return "huya_goodiptv"
}

// HelpDoc 处理器说明文档
func (h *huyaGoodIptvHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 基于肥羊大佬的 goodiptv 接口进行开发")
	sb.WriteString("\n2. 将虎牙房间号 id 传递到处理器接口的 ch 参数即可")
	sb.WriteString("\n3. 播放器需要支持 flv 格式播放")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *huyaGoodIptvHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *huyaGoodIptvHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (h *huyaGoodIptvHandler) Enabled() bool {
	return true
}
