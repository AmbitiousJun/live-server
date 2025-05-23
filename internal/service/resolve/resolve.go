package resolve

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

// handlerMap 将处理器的名称作为 key 存放到 map 中, 便于快速读取
var handlerMap = map[string]Handler{}

// handlerMapOpMutex 操作 handlerMap 的同步锁
var handlerMapOpMutex = sync.RWMutex{}

type TsProxyMode string

var (
	ModeCustom TsProxyMode = "custom" // 使用自定义的 ts 代理接口策略
	ModeLocal  TsProxyMode = "local"  // 使用本地代理策略
)

// HandleParams 处理参数
type HandleParams struct {
	ChName      string      // 频道简称
	UrlEnv      string      // 存储远程地址的环境变量名
	ProxyM3U    bool        // 是否代理 m3u
	ProxyTs     bool        // 是否代理 ts
	TsProxyMode             // ts 代理模式
	Format      string      // 要处理的直播格式
	ClientIp    string      // 客户端 ip
	ClientHost  string      // 客户端请求的主机前缀
	Headers     http.Header // 请求头
}

// ResultType 处理器的处理结果
type ResultType string

const (
	ResultRedirect ResultType = "redirect" // 重定向
	ResultProxy    ResultType = "proxy"    // 本地代理
)

// HandleResult 处理器的处理结果
type HandleResult struct {
	Type   ResultType  // 响应类型
	Url    string      // 响应地址, 用于重定向
	Code   int         // 响应状态码, 用于本地代理
	Header http.Header // 响应头, 用于本地代理
	Body   []byte      // 响应体, 用于本地代理
}

// Handler 直播响应处理器
type Handler interface {

	// Handle 处理直播, 返回一个用于重定向的远程地址
	Handle(HandleParams) (HandleResult, error)

	// Name 处理器名称
	Name() string

	// HelpDoc 处理器说明文档
	HelpDoc() string

	// SupportProxy 是否支持 m3u 代理
	//
	// 如果返回 true, 会自动在帮助文档中加入标记
	SupportM3UProxy() bool

	// SupportCustomHeaders 是否支持自定义请求头
	//
	// 如果返回 true, 会自动在帮助文档中加入标记
	SupportCustomHeaders() bool

	// Enabled 标记处理器是否是启用状态
	Enabled() bool
}

// RegisterHandler 注册处理器到内存中
func RegisterHandler(handler Handler) {
	if handler == nil || !handler.Enabled() {
		return
	}
	handlerMapOpMutex.Lock()
	defer handlerMapOpMutex.Unlock()
	handlerMap[handler.Name()] = handler
	log.Printf(colors.ToBlue("处理器 %s 初始化完成"), handler.Name())
}

// GetHandler 根据处理器名称获取处理器
func GetHandler(name string) (Handler, bool) {
	handlerMapOpMutex.RLock()
	defer handlerMapOpMutex.RUnlock()
	handler, ok := handlerMap[name]
	if !ok || !handler.Enabled() {
		return nil, false
	}
	return handler, true
}

// HelpDoc 输出所有解析处理器的帮助文档
func HelpDoc() string {
	sb := strings.Builder{}

	// 程序说明
	sb.WriteString("\n<strong>live-server " + constant.Version + " 帮助文档</strong>\n")
	sb.WriteString("\n仓库地址：<a href=\"" + constant.RepoAddr + "\" target=\"_blank\">" + constant.RepoAddr + "</a>")
	sb.WriteString("\n使用说明：本项目仅限个人测试使用，切勿大肆分享传播！！！")
	sb.WriteString("\n调用方式：<a href=\"${clientOrigin}/handler/youtube/ch/6IquAgfvYmc?proxy_m3u=1&proxy_ts=1$YT\" target=\"_blank\">${clientOrigin}/handler/{处理器名}/ch/{频道名称}[?{query 参数}][${频道注释}]</a>")
	sb.WriteString("\n")

	// 程序密钥相关
	sb.WriteString("\n<strong>程序密钥说明：</strong>")
	sb.WriteString("\n1. 每次启动运行时，会初始化一次密钥并在日志中输出，注意查看")
	sb.WriteString("\n2. 程序密钥会保存在【数据目录/secret.txt】文件中，初始化时如果文件不存在会自动生成一个随机密钥")
	sb.WriteString("\n3. 可通过修改 secret.txt 文件自定义密钥")
	sb.WriteString("\n4. 对于一些不便公开的接口，需要将程序密钥设置为接口的 secret 参数才可以成功调用")
	sb.WriteString("\n")

	// 地域白名单相关
	sb.WriteString("\n<strong>地域白名单说明：</strong>")
	sb.WriteString("\n1. 如果服务部署在公网上，推荐使用这个功能")
	sb.WriteString("\n2. 如果没有设置过白名单，默认允许所有 ip 进行访问")
	sb.WriteString("\n3. 功能生效的前提是程序能正确获取到客户端 ip 的归属地信息（可以在运行日志中加以确认）")
	sb.WriteString("\n4. 接口调用格式参考下方说明，其中 area 传递格式是将地域用 '/' 符号逐级隔开，如（广东/佛山/南海）")
	sb.WriteString("\n5. 场景举例说明：在只设置了一个白名单的前提下，如果【area=广东】，则除了广东地区之外的所有 ip 都会被屏蔽，如果【area=广东/佛山】，则广东省内除了佛山市之外的其他市的 ip 都会被屏蔽")
	sb.WriteString("\n6. 白名单只需要设置一次，数据会持久化在【数据目录/white_area.json】文件中")
	sb.WriteString("\n")

	// 接口调用相关
	sb.WriteString("\n<strong>接口调用说明：</strong>")
	sb.WriteString("\n请前往配置页进行配置: <a target=\"_blank\" href=\"${clientOrigin}/config?secret=\">${clientOrigin}/config?secret={程序密钥}</a>")
	sb.WriteString("\n")

	// 代理参数
	sb.WriteString("\n<strong>代理参数说明：</strong>")
	sb.WriteString("\n1. 如果处理器支持 M3U 代理, 则可以在调用处理器时传递代理参数进行代理")
	sb.WriteString("\n2. 代理参数 ①：proxy_m3u => 是否代理 m3u，传递 1 则开启代理，其他值无效")
	sb.WriteString("\n3. 代理参数 ②：proxy_ts  => 是否代理 ts 切片，传递 1 则开启代理，其他值无效")
	sb.WriteString("\n4. 代理参数 ③：ts_proxy_mode => ts 代理模式，传递 custom 则使用自定义的代理接口，传递 local 则使用本地代理，传递该参数时会覆盖全局配置")
	sb.WriteString("\n5. 开启切片代理时，会消耗服务器流量")
	sb.WriteString("\n6. 代理功能可以正常使用的前提是服务器的网络环境是能够和直播源进行连通的")
	sb.WriteString("\n7. 举例：${clientOrigin}/handler/345/ch/cctv13?proxy_m3u=1&proxy_ts=1&ts_proxy_mode=custom&format=1")
	sb.WriteString("\n")

	// 自定义请求头
	sb.WriteString("\n<strong>自定义请求头说明：</strong>")
	sb.WriteString("\n1. 如果处理器支持自定义请求头, 则可以在调用处理器时传递 headers 参数进行自定义")
	sb.WriteString("\n2. headers 格式: key1[[[:]]]value1[[[:]]]key2[[[:]]]value2")
	sb.WriteString("\n3. 举例：${clientOrigin}/handler/raw_m3u/ch/1?url_env=test_ch&proxy_m3u=1&proxy_ts=1&headers=Referer[[[:]]]https://www.baidu.com[[[:]]]User-Agent[[[:]]]okhttp")
	sb.WriteString("\n")

	// 自定义切片代理接口
	sb.WriteString("\n<strong>自定义切片代理接口说明：</strong>")
	sb.WriteString("\n1. 原理：借助 Cloudflare 的免费 Worker 来代理切片，避免 live-server 所在服务器的流量消耗")
	sb.WriteString("\n2. 准备工作：(1) Cloudflare 账号 (2) Github 账号 (3) 有自定义域名托管在 Cloudflare 上（网上有教程）(4) 知道如何将 Github 仓库连接部署到 Cloudflare Pages 上（网上有教程）")
	sb.WriteString("\n3. fork live-server 仓库到自己的 Github 账号上")
	sb.WriteString("\n4. 在 Cloudflare Pages 上连接这个 fork 仓库进行部署 (若对 Worker 比较熟悉, 建议直接复制本项目的 Worker 源代码到 Cloudflare 手动部署 Worker 而不是 Page)")
	sb.WriteString("\n5. 部署好之后为这个服务设置一个自定义子域名，主域名必须是托管在 Cloudflare 上的，然后等待一段时间")
	sb.WriteString("\n6. 浏览器访问 [https://{自定义子域名}] 如果返回 [Empty remote] 表示部署成功")
	sb.WriteString("\n7. 设置 live-server 的环境变量 [custom_ts_proxy_host] 值为 [https://{自定义的子域名}]")
	sb.WriteString("\n8. 设置 live-server 的环境变量 [custom_ts_proxy_enable] 值为 [1] 使用自定义的代理接口")
	sb.WriteString("\n")

	// 自定义处理异常视频
	sb.WriteString("\n<strong>自定义处理异常视频说明：</strong>")
	sb.WriteString("\n1. 处理器处理异常时，默认会返回 400 错误异常, 播放器会直接停止播放")
	sb.WriteString("\n2. 如需友好地让客户端知道发生了错误, 可以自己准备一个 mp4 视频文件并托管在任意可访问的 http 服务上, 获得一个类似 [https://raw.githubusercontent.com/AmbitiousJun/AmbitiousJun/refs/heads/master/404/GEM.mp4] 的视频地址")
	sb.WriteString("\n3. 将视频地址设置到 live-server 的环境变量 [fallback_mp4] 上即可自动生效")
	sb.WriteString("\n")

	// warp 说明
	sb.WriteString("\n<strong>warp 自动修复被封禁 ip 说明：</strong>")
	sb.WriteString("\n1. 此功能适用场景: 1) 通过二进制运行 live-server 2) vps 使用 <a href=\"https://github.com/yonggekkk/warp-yg\" target=\"_blank\">ygkkk</a> 的 warp 脚本, 且是已开启状态")
	sb.WriteString("\n2. 描述：当处理器检测到 ip 被封禁时, 自动运行脚本刷新可用 ip 进行解锁")
	sb.WriteString("\n3. 使用方式: 设置 live-server 的环境变量 [" + constant.Warp_ExecPathEnvKey + "] 值为 warp 脚本的绝对路径, 如：[/root/CFwarp.sh], 程序会自动调用")
	sb.WriteString("\n")

	// 处理器文档
	sb.WriteString("\n<strong>处理器说明：</strong>")
	handlerMapOpMutex.RLock()
	defer handlerMapOpMutex.RUnlock()
	for name, handler := range handlerMap {
		if !handler.Enabled() {
			continue
		}

		sb.WriteString("\n\n=====")
		sb.WriteString("\n处理器名：")
		sb.WriteString(name)
		if handler.SupportM3UProxy() {
			sb.WriteString("\n【该处理器支持 M3U 代理】")
		}
		if handler.SupportCustomHeaders() {
			sb.WriteString("\n【该处理器支持自定义 headers】")
		}
		sb.WriteString("\n处理器说明：")
		sb.WriteString(handler.HelpDoc())
	}

	sb.WriteString("\n\n\n\n\n\n\n")

	return sb.String()
}
