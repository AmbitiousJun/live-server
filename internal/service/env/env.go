package env

import "sync"

// globalEnv 存放环境变量
var globalEnv = sync.Map{}

// Set 设置环境变量
func Set(key, value string) {
	globalEnv.Store(key, value)
}

// Get 获取环境变量
func Get(key string) (string, bool) {
	val, ok := globalEnv.Load(key)
	if !ok {
		return "", false
	}
	return val.(string), true
}
