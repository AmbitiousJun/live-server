package web

import "strings"

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
