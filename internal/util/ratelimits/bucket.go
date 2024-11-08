package ratelimits

import (
	"sync"
	"time"
)

// bucket 令牌桶实现
type bucket struct {

	// ok 标记当前令牌桶是否是有效状态
	ok bool

	// tokenBuf 存放 token, 通道大小代表最大 token 数
	tokenBuf chan struct{}

	// mu 读写并发控制
	mu sync.Mutex
}

// NewBucket 创建一个令牌桶实例
//
// 令牌桶会启动一个单独的 goroutine, 负责生成令牌
// 每隔 sched 时间, 生成 tokenPerProduce 个 token
//
// maxToken 指定了桶中的最大 token 数
func NewBucket(tokenPerProduce uint, sched time.Duration, maxToken uint) Bucket {
	if tokenPerProduce <= 0 || sched <= 0 || maxToken <= 0 {
		panic("invalid bucket params")
	}

	b := bucket{ok: true}
	b.tokenBuf = make(chan struct{}, maxToken)

	go func() {
		ticker := time.NewTicker(sched)
		for range ticker.C {
			// 通道满则丢弃 token
			select {
			case b.tokenBuf <- struct{}{}:
			default:
			}
		}
	}()

	return &b
}

// Consume 消耗一定数量的令牌
// 如果桶中没有足量的令牌, 该方法会阻塞线程直至成功消耗相应数量的令牌
func (b *bucket) Consume(token uint) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.ok {
		panic("current bucket has been destroyed")
	}

	// 获取指定个数的 token
	for i := 1; i <= int(token); i++ {
		<-b.tokenBuf
	}
}

// TryConsume 非阻塞消耗一定数量的令牌
// 返回是否消耗成功
func (b *bucket) TryConsume(token uint) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.ok {
		panic("current bucket has been destroyed")
	}

	if len(b.tokenBuf) < int(token) {
		return false
	}

	// 获取指定个数的 token
	for i := 1; i <= int(token); i++ {
		<-b.tokenBuf
	}
	return true
}

// Destroy 销毁令牌桶
func (b *bucket) Destroy() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if !b.ok {
		return
	}

	close(b.tokenBuf)
	b.ok = false
}
