package web

import "github.com/AmbitiousJun/live-server/internal/service/env"

// FallbackMp4Env 在环境变量中存储处理器异常提示 mp4 视频地址的环境变量名
const FallbackMp4Env = "fallback_mp4"

// GetFallbackMp4Url 获取用户配置在环境变量中的处理器异常提示 mp4 视频地址
func GetFallbackMp4Url() (string, bool) {
	if mp4, ok := env.Get(FallbackMp4Env); ok {
		return mp4, true
	}
	return "", false
}
