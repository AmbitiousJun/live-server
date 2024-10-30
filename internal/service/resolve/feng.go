package resolve

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (

	// Env_FengToken 凤凰卫视接口请求 token 环境变量
	Env_FengToken = "feng_token"

	// fengAuthUrl 凤凰卫视直播源授权接口
	fengAuthUrl = "https://m.fengshows.com/api/v3/hub/live/auth-url?live_qa=FHD"

	// fengUpdateUrl 凤凰卫视 token 刷新接口
	//
	// 刷新时需要携带一个有效的 token,
	// 并且成功刷新后会使得原有的 token 立即失效
	fengUpdateUrl = "http://m.fengshows.com/user/oauth/update"
)

func init() {
	autoRefreshFengToken()
	registerHandler(new(fengHandler))
}

// autoRefreshFengToken 定时自动刷新 token
func autoRefreshFengToken() {
	env.SetAutoRefresh(Env_FengToken, func(curVal string) (string, error) {
		header := make(http.Header)
		header.Set("Token", curVal)
		resp, err := https.Request(http.MethodPost, fengUpdateUrl, header, nil)
		if err != nil {
			return "", fmt.Errorf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("错误的响应码: %d", resp.StatusCode)
		}

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("读取响应体失败: %v", err)
		}

		bodyJson, err := jsons.New(string(bodyBytes))
		if err != nil {
			return "", fmt.Errorf("解析响应失败, err: %v, 原始响应: %s", err, string(bodyBytes))
		}

		newVal, ok := bodyJson.Attr("data").Attr("token").String()
		if !ok {
			return "", fmt.Errorf("解析响应失败, 获取不到新 token, 原始响应: %s", string(bodyBytes))
		}

		return newVal, nil
	}, time.Hour*6)
}

// fengChannels 记录支持的凤凰卫视频道
var fengChannels = map[string]string{
	"fhzx": "7c96b084-60e1-40a9-89c5-682b994fb680", // 凤凰资讯
	"fhzw": "f7f48462-9b13-485b-8101-7b54716411ec", // 凤凰中文
	"fhxg": "15e02d92-1698-416c-af2f-3e9a872b4d78", // 凤凰香港
}

// fengHandler 凤凰卫视电视直播处理器
type fengHandler struct{}

func (fengHandler) Name() string {
	return "feng"
}

func (fengHandler) Handle(params HandleParams) (HandleResult, error) {
	// 1 检查 token 变量
	token, ok := env.Get(Env_FengToken)
	if !ok {
		return HandleResult{}, fmt.Errorf("请先设置环境变量: %s", Env_FengToken)
	}

	// 2 判断频道是否支持
	liveId, ok := fengChannels[params.ChName]
	if !ok {
		return HandleResult{}, fmt.Errorf("不支持的频道: %s", params.ChName)
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
		return HandleResult{}, fmt.Errorf("请求失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HandleResult{}, fmt.Errorf("授权失败, 响应码: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return HandleResult{}, fmt.Errorf("读取响应异常: %v", err)
	}
	resJson, err := jsons.New(string(bytes))
	if err != nil {
		return HandleResult{}, fmt.Errorf("JSON 转换失败: %v", err)
	}

	liveUrl, ok := resJson.Attr("data").Attr("live_url").String()
	if !ok {
		return HandleResult{}, fmt.Errorf("获取直播地址失败, 原始响应: %s", resJson)
	}

	return HandleResult{Type: ResultRedirect, Url: liveUrl}, nil
}

// HelpDoc 处理器说明文档
func (fengHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 手机安装《凤凰秀》app，登录 app 后使用抓包工具手动抓取位于请求头的 jwt 登录 token")
	sb.WriteString("\n2. 将 token 设置到环境变量中即可正常观看, key: feng_token")
	sb.WriteString("\n3. 设置完 token 之后就不要再去打开 app 了，否则现有 token 失效")
	sb.WriteString("\n4. 程序每隔 6 小时自动刷新 token")
	sb.WriteString("\n5. 支持的频道: fhzw(凤凰中文)、fhzx(凤凰资讯)、fhxg(凤凰香港)")
	return sb.String()
}
