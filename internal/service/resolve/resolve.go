package resolve

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

// handlerMap 将处理器的名称作为 key 存放到 map 中, 便于快速读取
var handlerMap = map[string]Handler{}

// handlerMapOpMutex 操作 handlerMap 的同步锁
var handlerMapOpMutex = sync.RWMutex{}

// HandleParams 处理参数
type HandleParams struct {
	ChName   string // 频道简称
	UrlEnv   string // 存储远程地址的环境变量名
	ProxyM3U bool   // 是否代理 m3u
	ProxyTs  bool   // 是否代理 ts
	Format   string // 要处理的直播格式
	ClientIp string // 客户端 ip
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
}

// RegisterHandler 注册处理器到内存中
func RegisterHandler(handler Handler) {
	if handler == nil {
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
	if !ok {
		return nil, false
	}
	return handler, true
}

// HelpDoc 输出所有解析处理器的帮助文档
func HelpDoc() string {
	sb := strings.Builder{}
	version := "v1.13.13"

	// 程序说明
	sb.WriteString("\nlive-server " + version + " 帮助文档\n")
	sb.WriteString("\n仓库地址：https://github.com/AmbitiousJun/live-server")
	sb.WriteString("\n项目说明：本项目仅限个人测试使用，切勿大肆传播！！")
	sb.WriteString("\n")

	// 程序密钥相关
	sb.WriteString("\n程序密钥说明：")
	sb.WriteString("\n1. 每次启动运行时，会初始化一次密钥并在日志中输出，注意查看")
	sb.WriteString("\n2. 程序密钥会保存在【数据目录/secret.txt】文件中，初始化时如果文件不存在会自动生成一个随机密钥")
	sb.WriteString("\n3. 可通过修改 secret.txt 文件自定义密钥")
	sb.WriteString("\n4. 对于一些不便公开的接口，需要将程序密钥设置为接口的 secret 参数才可以成功调用")
	sb.WriteString("\n")

	// 地域白名单相关
	sb.WriteString("\n地域白名单说明：")
	sb.WriteString("\n1. 如果服务部署在公网上，推荐使用这个功能")
	sb.WriteString("\n2. 如果没有设置过白名单，默认允许所有 ip 进行访问")
	sb.WriteString("\n3. 功能生效的前提是程序能正确获取到客户端 ip 的归属地信息（可以在运行日志中加以确认）")
	sb.WriteString("\n4. 接口调用格式参考下方说明，其中 area 传递格式是将地域用 '/' 符号逐级隔开，如（广东/佛山/南海）")
	sb.WriteString("\n5. 场景举例说明：在只设置了一个白名单的前提下，如果【area=广东】，则除了广东地区之外的所有 ip 都会被屏蔽，如果【area=广东/佛山】，则广东省内除了佛山市之外的其他市的 ip 都会被屏蔽")
	sb.WriteString("\n6. 白名单只需要设置一次，数据会持久化在【数据目录/white_area.json】文件中")
	sb.WriteString("\n")

	// 接口调用相关
	sb.WriteString("\n接口调用说明：")
	sb.WriteString("\n如果不熟悉接口调用, 可前往小白专用配置页进行配置: ${clientOrigin}/config?secret={程序密钥}")
	sb.WriteString("\n1. 设置环境变量(GET) => ${clientOrigin}/env?key={变量名}&value={变量值}&secret={程序密钥}")
	sb.WriteString("\n2. 帮助文档(GET) => ${clientOrigin}/help")
	sb.WriteString("\n3. 调用处理器(GET) => ${clientOrigin}/handler/{处理器名}/ch/{频道名}[可选的 query 参数，如：?url_env=remote_m3u_v6]")
	sb.WriteString("\n4. ip 黑名单(GET) => ${clientOrigin}/black_ip?ip={要加入黑名单的地址}&secret={程序密钥}")
	sb.WriteString("\n5. 设置地域白名单(GET) => ${clientOrigin}/white_area/set?area={要加入白名单的地域}&secret={程序密钥}")
	sb.WriteString("\n6. 移除地域白名单(GET) => ${clientOrigin}/white_area/del?area={要移除白名单的地域}&secret={程序密钥}")
	sb.WriteString("\n")

	// 代理参数
	sb.WriteString("\n代理参数说明：")
	sb.WriteString("\n1. 如果处理器支持 M3U 代理, 则可以在调用处理器时传递代理参数进行代理")
	sb.WriteString("\n2. 代理参数 ①：proxy_m3u => 是否代理 m3u，传递 1 则开启代理，其他值无效")
	sb.WriteString("\n3. 代理参数 ②：proxy_ts  => 是否代理 ts 切片，传递 1 则开启代理，其他值无效")
	sb.WriteString("\n4. 开启切片代理时，会消耗服务器流量")
	sb.WriteString("\n5. 代理功能可以正常使用的前提是服务器的网络环境是能够和直播源进行连通的")
	sb.WriteString("\n6. 举例：/handler/third_gdtv/ch/xwpd?proxy_m3u=1")
	sb.WriteString("\n")

	// 自定义切片代理接口
	sb.WriteString("\n自定义切片代理接口说明：")
	sb.WriteString("\n1. 原理：借助 Cloudflare 的免费 Worker 来代理切片，避免 live-server 所在服务器的流量消耗")
	sb.WriteString("\n2. 准备工作：(1) Cloudflare 账号 (2) Github 账号 (3) 有自定义域名托管在 Cloudflare 上（网上有教程）(4) 知道如何将 Github 仓库连接部署到 Cloudflare Pages 上（网上有教程）")
	sb.WriteString("\n3. fork live-server 仓库到自己的 Github 账号上")
	sb.WriteString("\n4. 在 Cloudflare Pages 上连接这个 fork 仓库进行部署")
	sb.WriteString("\n5. 部署好之后为这个服务设置一个自定义子域名，主域名必须是托管在 Cloudflare 上的，然后等待一段时间")
	sb.WriteString("\n6. 浏览器访问 [https://{自定义子域名}] 如果返回 [Empty remote] 表示部署成功")
	sb.WriteString("\n7. 设置 live-server 的环境变量 [custom_ts_proxy_host] 值为 [https://{自定义的子域名}]")
	sb.WriteString("\n8. 设置 live-server 的环境变量 [custom_ts_proxy_enable] 值为 [1] 使用自定义的代理接口")
	sb.WriteString("\n")

	// 处理器文档
	sb.WriteString("\n处理器说明：")
	handlerMapOpMutex.RLock()
	defer handlerMapOpMutex.RUnlock()
	for name, handler := range handlerMap {
		sb.WriteString("\n\n=====")
		sb.WriteString("\n处理器名：")
		sb.WriteString(name)
		if handler.SupportM3UProxy() {
			sb.WriteString("\n【该处理器支持 M3U 代理】")
		}
		sb.WriteString("\n处理器说明：")
		sb.WriteString(handler.HelpDoc())
	}

	sb.WriteString("\n\n\n\n\n\n\n")

	return sb.String()
}
