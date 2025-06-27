package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/base64s"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
	"github.com/AmbitiousJun/live-server/internal/util/ratelimits"
)

func init() {
	f := &fengHandler{
		tokenEnv:      "feng_token",
		authUrl:       base64s.MustDecodeString("aHR0cHM6Ly9tLmZlbmdzaG93cy5jb20vYXBpL3YzL2h1Yi9saXZlL2F1dGgtdXJsP2xpdmVfcWE9RkhE"),
		authUpdateUrl: base64s.MustDecodeString("aHR0cDovL20uZmVuZ3Nob3dzLmNvbS91c2VyL29hdXRoL3VwZGF0ZQ=="),
		numBucket:     ratelimits.NewBucket(1, time.Second*10, 3),
		rateBucket:    ratelimits.NewBucket(1, time.Millisecond*2500, 1),
		channels: map[string]string{
			"fhzx": "7c96b084-60e1-40a9-89c5-682b994fb680", // 凤凰资讯
			"fhzw": "f7f48462-9b13-485b-8101-7b54716411ec", // 凤凰中文
			"fhxg": "15e02d92-1698-416c-af2f-3e9a872b4d78", // 凤凰香港
		},
	}
	f.autoRefreshFengToken()
	resolve.RegisterHandler(f)
}

// fengHandler 凤凰卫视电视直播处理器
type fengHandler struct {

	// tokenEnv 凤凰秀接口请求 token 参数的环境变量名
	tokenEnv string

	// authUrl 直播源授权接口
	authUrl string

	// authUpdateUrl token 刷新接口
	//
	// 刷新时需要携带一个有效的 token,
	// 并且成功刷新后会使得原有的 token 立即失效
	authUpdateUrl string

	// numBucket 限制一段时间间隔最多请求数
	numBucket ratelimits.Bucket

	// rateBucket 限制两个请求之间的最小间隔
	rateBucket ratelimits.Bucket

	// channels 支持的频道列表
	channels map[string]string
}

func (f *fengHandler) Name() string {
	return "feng"
}

func (f *fengHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 1 检查 token 变量
	token, ok := env.Get(f.tokenEnv)
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("请先设置环境变量: %s", f.tokenEnv)
	}

	// 2 判断频道是否支持
	liveId, ok := f.channels[params.ChName]
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("不支持的频道: %s", params.ChName)
	}

	f.Wait()

	// 3 请求授权接口
	u, _ := url.Parse(f.authUrl)
	q := u.Query()
	q.Set("live_id", liveId)
	u.RawQuery = q.Encode()
	header := make(http.Header)
	header.Set("Token", token)
	header.Set("fengshows-client", "app(ios,5041701);iPhone11,8;16.1.2")
	header.Set("client-lang", "zh-CN")
	header.Set("User-Agent", "FengWatch/5.4.17 (iPhone; iOS 16.1.2; Scale/2.00)")
	resp, err := https.Get(u.String()).Header(header).Do()
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("请求失败: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return resolve.HandleResult{}, fmt.Errorf("授权失败, 响应码: %d", resp.StatusCode)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("读取响应异常: %v", err)
	}
	resJson, err := jsons.New(string(bytes))
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("JSON 转换失败: %v", err)
	}

	liveUrl, ok := resJson.Attr("data").Attr("live_url").String()
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("获取直播地址失败, 原始响应: %s", resJson)
	}

	return resolve.HandleResult{Type: resolve.ResultRedirect, Url: liveUrl}, nil
}

// HelpDoc 处理器说明文档
func (f *fengHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 手机安装《凤凰秀》app，登录 app 后使用抓包工具手动抓取位于请求头的 jwt 登录 token")
	sb.WriteString("\n2. 将 token 设置到环境变量中即可正常观看, key: feng_token")
	sb.WriteString("\n3. 设置完 token 之后就不要再去打开 app 了，否则现有 token 失效")
	sb.WriteString("\n4. 程序每隔 6 小时自动刷新 token")
	sb.WriteString("\n5. 支持的频道: fhzw(凤凰中文)、fhzx(凤凰资讯)、fhxg(凤凰香港)")
	sb.WriteString("\n6. 该处理器设置了请求速率限制, 每分钟允许请求 6 次，仅自用不适合分享，请避免滥用")
	sb.WriteString("\n7. 如不会自己抓 token，可以在此页面：<a target=\"_blank\" href=\"${clientOrigin}/feng/auth?secret=\">${clientOrigin}/feng/auth?secret={程序密钥}</a> 使用手机号登录授权获取，不保证可用性，有问题可以到 issue 区反馈")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
// 如果返回 true, 会自动在帮助文档中加入标记
func (f *fengHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (f *fengHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (f *fengHandler) Enabled() bool {
	return true
}

// Wait 请求速率限制
func (f *fengHandler) Wait() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		f.numBucket.Consume(1)
	}()
	go func() {
		defer wg.Done()
		f.rateBucket.Consume(1)
	}()
	wg.Wait()
}

// autoRefreshFengToken 定时自动刷新 token
func (f *fengHandler) autoRefreshFengToken() {
	env.SetAutoRefresh(f.tokenEnv, func(curVal string) (string, error) {
		header := make(http.Header)
		header.Set("Token", curVal)
		resp, err := https.Post(f.authUpdateUrl).Header(header).Do()
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
