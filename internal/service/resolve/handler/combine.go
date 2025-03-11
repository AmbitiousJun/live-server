package handler

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
)

// combineHandler 接口聚合处理器
type combineHandler struct {

	// chSeg 分割多个直播源的分隔符
	chSeg string
}

func init() {
	resolve.RegisterHandler(&combineHandler{
		chSeg: constant.HeadersSeg,
	})
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (h *combineHandler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	// 1 取出环境变量
	if strs.AnyEmpty(params.UrlEnv) {
		return resolve.HandleResult{}, errors.New("请传入环境变量")
	}
	envVal, ok := env.Get(params.UrlEnv)
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("环境变量 %s 不存在", params.UrlEnv)
	}

	// 2 解析直播源列表
	chs := strings.Split(envVal, h.chSeg)
	if len(chs) == 0 {
		return resolve.HandleResult{}, fmt.Errorf("直播源列表为空, 请使用 %s 分隔多个直播源", h.chSeg)
	}

	// 3 随机选择一个直播源
	return resolve.HandleResult{
		Type: resolve.ResultRedirect,
		Url:  chs[rand.IntN(len(chs))],
	}, nil
}

// Name 处理器名称
func (h *combineHandler) Name() string {
	return "combine"
}

// HelpDoc 处理器说明文档
func (h *combineHandler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 本处理器作用是聚合多个直播源, 并在客户端请求时随机返回重定向")
	sb.WriteString("\n2. 先将自定义的多个直播源设置到程序的环境变量中, 变量名称任意, 如: combine_hyxw")
	sb.WriteString("\n3. 多个源之间使用 " + h.chSeg + " 进行分隔")
	sb.WriteString("\n4. 频道名称 ch 可任意传递一个非空值")
	sb.WriteString("\n4. 请求示例: ${clientOrigin}/handler/combine/ch/1?url_env=combine_hyxw")
	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *combineHandler) SupportM3UProxy() bool {
	return false
}

// SupportCustomHeaders 是否支持自定义请求头
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (h *combineHandler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (h *combineHandler) Enabled() bool {
	return true
}
