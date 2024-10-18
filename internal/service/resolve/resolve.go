package resolve

import "net/http"

// allHandlers 存放所有的直播处理器
var allHandlers = []Handler{
	new(fengHandler),
	new(remoteM3UHandler),
	new(thirdGdtvHandler),
	new(sttvHandler),
}

// handlerMap 将处理器的名称作为 key 存放到 map 中, 便于快速读取
var handlerMap map[string]Handler

func init() {
	handlerMap = make(map[string]Handler)
	for _, handler := range allHandlers {
		handlerMap[handler.Name()] = handler
	}
}

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
}

// GetHandler 根据处理器名称获取处理器
func GetHandler(name string) (Handler, bool) {
	handler, ok := handlerMap[name]
	if !ok {
		return nil, false
	}
	return handler, true
}
