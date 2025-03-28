package base64s

import (
	"encoding/base64"
	"log"
)

// MustDecode 解码 base64 字符串成字节切片
//
// 解码失败时触发 panic
func MustDecode(enc string) []byte {
	res, err := base64.StdEncoding.DecodeString(enc)
	if err != nil {
		log.Panicf("base64 解码异常, enc: %s, err: %v", enc, err)
	}
	return res
}

// MustDecodeString 解码 base64 字符串成原文字符串
//
// 解码失败时触发 panic
func MustDecodeString(enc string) string {
	return string(MustDecode(enc))
}
