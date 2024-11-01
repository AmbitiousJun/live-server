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
	ChName string // 频道简称
	UrlEnv string // 存储远程地址的环境变量名
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

	// 程序说明
	sb.WriteString("\nlive-server 帮助文档\n")
	sb.WriteString("\n仓库地址：https://github.com/AmbitiousJun/live-server")
	sb.WriteString("\n项目说明：本项目仅限个人测试使用，切勿用于大肆传播！！")
	sb.WriteString("\n")

	// 接口调用相关
	sb.WriteString("\n接口调用说明：")
	sb.WriteString("\n1. 设置环境变量(GET) => :5666/env?key={变量名}&value={变量值}")
	sb.WriteString("\n2. 帮助文档(GET) => :5666/help")
	sb.WriteString("\n3. 调用处理器(GET) => :5666/handler/{处理器名}/ch/{频道名}[可选的 query 参数，如：?url_env=remote_m3u_v6]")
	sb.WriteString("\n4. ip 黑名单(GET) => :5666/black_ip?ip={要加入黑名单的地址}")
	sb.WriteString("\n")

	// 处理器文档
	sb.WriteString("\n处理器说明：")
	handlerMapOpMutex.RLock()
	defer handlerMapOpMutex.RUnlock()
	for name, handler := range handlerMap {
		sb.WriteString("\n=====")
		sb.WriteString("\n处理器名：")
		sb.WriteString(name)
		sb.WriteString("\n处理器说明：")
		sb.WriteString(handler.HelpDoc())
	}

	sb.WriteString("\n\n\n\n\n\n\n")

	return sb.String()
}
