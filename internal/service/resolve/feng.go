package resolve

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (

	// Env_FengToken 凤凰卫视接口请求 token 环境变量
	Env_FengToken = "feng_token"
	// fengAuthUrl 凤凰卫视直播源授权接口
	fengAuthUrl = "https://m.fengshows.com/api/v3/hub/live/auth-url?live_qa=FHD"
)

// fengChannels 记录支持的凤凰卫视频道
var fengChannels = map[string]string{
	"fhzx": "7c96b084-60e1-40a9-89c5-682b994fb680", // 凤凰资讯
	"fhzw": "f7f48462-9b13-485b-8101-7b54716411ec", // 凤凰中文
}

// fengHandler 凤凰卫视电视直播处理器
type fengHandler struct{}

func (fengHandler) Name() string {
	return "feng"
}

func (fengHandler) Handle(params HandleParams) (string, error) {
	// 1 检查 token 变量
	token, ok := env.Get(Env_FengToken)
	if !ok {
		return "", fmt.Errorf("请先设置环境变量: %s", Env_FengToken)
	}

	// 2 判断频道是否支持
	liveId, ok := fengChannels[params.ChName]
	if !ok {
		return "", fmt.Errorf("不支持的频道: %s", params.ChName)
	}

	// 3 请求授权接口
	u, _ := url.Parse(fengAuthUrl)
	q := u.Query()
	q.Set("live_id", liveId)
	u.RawQuery = q.Encode()
	header := make(http.Header)
	header.Set("Token", token)
	resp, err := https.Request(http.MethodGet, u.String(), header, nil)
	if err != nil {
		return "", fmt.Errorf("请求失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("授权失败, 响应码: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应异常: %v", err)
	}
	resJson, err := jsons.New(string(bytes))
	if err != nil {
		return "", fmt.Errorf("JSON 转换失败: %v", err)
	}

	liveUrl, ok := resJson.Attr("data").Attr("live_url").String()
	if !ok {
		return "", fmt.Errorf("获取直播地址失败, 原始响应: %s", resJson)
	}

	return liveUrl, nil
}
