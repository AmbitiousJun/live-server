package subm3u

import (
	"sync"

	"github.com/AmbitiousJun/live-server/internal/util/encrypts"
)

const (
	MaxCacheSize = 10 // 内存中维护的订阅列表的最大个数
)

// cacheMap 存储订阅列表缓存
//
// key: 订阅列表 md5 => value: 整理好的订阅列表 map
var cacheMap = map[string]map[string][]Info{}

// opCacheMutex 列表缓存并发控制
var opCacheMutex = sync.RWMutex{}

// getCache 根据文本的原始内容获取缓存
func getCache(content string) (map[string][]Info, bool) {
	if MaxCacheSize < 1 {
		return nil, false
	}

	md5Hash := encrypts.Md5Hash(content)
	opCacheMutex.RLock()
	defer opCacheMutex.RUnlock()

	c, ok := cacheMap[md5Hash]
	if ok {
		return c, true
	}

	return nil, false
}

// updateCache 更新缓存
func updateCache(content string, infoMap map[string][]Info) {
	if MaxCacheSize < 1 {
		return
	}

	md5Hash := encrypts.Md5Hash(content)
	opCacheMutex.Lock()
	defer opCacheMutex.Unlock()

	// 内存中已经包含该列表, 直接更新
	if _, ok := cacheMap[md5Hash]; ok {
		cacheMap[md5Hash] = infoMap
		return
	}

	// 内存随机淘汰
	removeNum := len(cacheMap) - (MaxCacheSize - 1)
	removeKeys := []string{}
	for key := range cacheMap {
		if removeNum < 1 {
			break
		}
		removeKeys = append(removeKeys, key)
		removeNum--
	}
	for _, key := range removeKeys {
		delete(cacheMap, key)
	}

	cacheMap[md5Hash] = infoMap
}
