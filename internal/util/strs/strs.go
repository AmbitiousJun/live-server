package strs

import (
	"sort"
	"strings"
)

// AllNotEmpty 所有字符串都不为空时, 返回 true
func AllNotEmpty(strs ...string) bool {
	return !AnyEmpty(strs...)
}

// AnyEmpty 有任意一个字符串为空时, 返回 true
func AnyEmpty(strs ...string) bool {
	for _, str := range strs {
		if str = strings.TrimSpace(str); str == "" {
			return true
		}
	}
	return false
}

// Sort 将一个字符串进行字典序排序
func Sort(str string) string {
	runes := []rune(str)
	sort.SliceStable(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

// ReverseString 翻转字符串
func ReverseString(s string) string {
	runes := []rune(s) // 转换为rune切片
	// 双指针翻转
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes) // 转换回字符串
}