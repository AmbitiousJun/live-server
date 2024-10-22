package web

import (
	"errors"
	"log"

	"github.com/AmbitiousJun/live-server/internal/bloom"
	"github.com/toniphan21/go-bf"
)

const (
	BloomItemNums  = 10000             // 布隆过滤器存储的元素个数
	BloomErrorRate = 0.001             // 布隆过滤器的容错率
	Env_BlackIps   = "bloom_black_ips" // 存储布隆过滤器数据的环境变量值
)

// blackIpBf 存放黑名单 ip 的布隆过滤器实例
var blackIpBf bf.BloomFilter

func init() {
	filter, err := bf.New(
		bf.WithAccuracy(BloomErrorRate, BloomItemNums),
		bf.WithStorage(&bloom.EnvStorageFactory{EnvKey: Env_BlackIps}),
	)
	if err != nil {
		log.Printf("黑名单布隆过滤器初始化失败: %v", err)
		return
	}
	blackIpBf = filter
}

// AddBlackIp 添加黑名单 ip
func AddBlackIp(ip string) error {
	if blackIpBf == nil {
		return errors.New("黑名单布隆过滤器未初始化")
	}

	blackIpBf.Add([]byte(ip))
	return nil
}

// IsBlackIp 判断一个 ip 是否在黑名单中
func IsBlackIp(ip string) bool {
	if blackIpBf == nil {
		return true
	}

	return blackIpBf.Exists([]byte(ip))
}
