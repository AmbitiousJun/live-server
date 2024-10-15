package resolve

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/m3u8"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

// remoteM3UHandler 远程 m3u8 直播源处理器
type remoteM3UHandler struct{}

func (remoteM3UHandler) Name() string {
	return "remote_m3u"
}

func (remoteM3UHandler) Handle(params HandleParams) (HandleResult, error) {
	// 获取环境变量
	url, ok := env.Get(params.UrlEnv)
	if !ok {
		return HandleResult{}, fmt.Errorf("获取不到环境变量: %s", params.UrlEnv)
	}

	// 请求远程 m3u 文本
	resp, err := https.Request(http.MethodGet, url, nil, nil)
	if err != nil {
		return HandleResult{}, fmt.Errorf("请求远程地址失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析文本内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return HandleResult{}, fmt.Errorf("读取远程响应失败: %v", err)
	}
	infos, err := m3u8.ReadContent(string(bodyBytes))
	if err != nil {
		return HandleResult{}, fmt.Errorf("解析远程响应失败: %v, 原始响应: %s", err, string(bodyBytes))
	}

	// 获取用户请求的频道
	if destInfo, ok := infos[params.ChName]; ok {
		return HandleResult{Type: ResultRedirect, Url: destInfo.Url}, nil
	}

	return HandleResult{}, fmt.Errorf("匹配频道名称失败: %s, 请检查远程地址是否有效", params.ChName)
}
