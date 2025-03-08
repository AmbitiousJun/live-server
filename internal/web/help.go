package web

import (
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/constant"
)

// trimDollarSuffix 移除请求的 $ 后缀
// 从后往前查找第一个出现的 $ 符号位置, 移除它及其之后的子串后返回
// 若查找不到 $ 符号, 则返回原始串
func trimDollarSuffix(rawUrl string) string {
	lastDollarIdx := strings.LastIndex(rawUrl, "$")
	if lastDollarIdx == -1 {
		return rawUrl
	}
	return rawUrl[:lastDollarIdx]
}

// resolveHeaders 解析 headers 拼接字符串为 http.Header 对象
//
// headers 格式: key1${constant.HeadersSeg}value1${constant.HeadersSeg}key2${constant.HeadersSeg}value2
func resolveHeaders(headers string) http.Header {
	res := make(http.Header)
	segments := strings.Split(headers, constant.HeadersSeg)
	for i := 0; i < len(segments)-1; i += 2 {
		res.Add(segments[i], segments[i+1])
	}
	return res
}
