package bloom

import (
	"errors"
	"log"
	"math/big"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/toniphan21/go-bf"
)

// EnvStoreBigIntBase 在环境变量中存储的进制数
const EnvStoreBigIntBase = 62

// EnvStorage 布隆过滤器接入 env 环境变量, 使用字符串模拟
type EnvStorage struct {
	capacity uint32 // 码位个数
	envKey   string // 环境变量键值
	mu       sync.RWMutex
}

// Set 标记指定码位
func (es *EnvStorage) Set(index uint32) {
	es.mu.Lock()
	defer es.mu.Unlock()

	var bigInt big.Int
	bloomStrHex, ok := env.Get(es.envKey)
	if ok {
		bigInt.SetString(bloomStrHex, EnvStoreBigIntBase)
	}
	bigInt.SetBit(&bigInt, int(index), 1)
	env.Set(es.envKey, bigInt.Text(EnvStoreBigIntBase))
}

// Get 获取指定码位的标记
func (es *EnvStorage) Get(index uint32) bool {
	es.mu.RLock()
	defer es.mu.RUnlock()

	bloomStrHex, ok := env.Get(es.envKey)
	if !ok {
		return false
	}
	var bigInt big.Int
	if _, ok := bigInt.SetString(bloomStrHex, EnvStoreBigIntBase); !ok {
		log.Printf("环境变量异常: 无法通过 Hex 恢复布隆过滤器")
		return false
	}

	return bigInt.Bit(int(index)) == 1
}

// Capacity 布隆过滤器容量
func (es *EnvStorage) Capacity() uint32 {
	return es.capacity
}

type EnvStorageFactory struct {
	EnvKey string // 环境变量键值
}

func (esf *EnvStorageFactory) Make(capacity uint32) (bf.Storage, error) {
	if esf.EnvKey == "" {
		return nil, errors.New("EnvKey 不能为空")
	}
	return &EnvStorage{capacity: capacity, envKey: esf.EnvKey}, nil
}
