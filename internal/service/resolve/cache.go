package resolve

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

// PreCacheReq 预缓存请求结构, 用于并发排队
type PreCacheReq[T any] struct {

	// T 调用方的业务参数
	T T

	// resChan 用于接收最终需要缓存的结果值
	resChan chan struct {
		res string
		err error
	}
}

// Cacher 通用的解析缓存器, 避免客户端频繁请求
// 可通过泛型指定业务请求参数
type Cacher[T any] struct {

	// CalcCacheKey 计算 cacheKey 的钩子, 需要调用方实现
	CalcCacheKey func(T) string

	// RecoverCacheKey 恢复 cacheKey 的钩子, 需要调用方实现
	RecoverCacheKey func(string) (T, bool)

	// FetchValue 获取指定入参的最新值钩子, 需要调用方实现
	FetchValue func(T) (string, error)

	// UpdateComplete 缓存刷新完成调用的钩子, 需要调用方实现
	UpdateComplete func(success, fail, remove int)

	// valid 标记当前缓存处理器是否是正在启用状态
	valid bool

	// caches 存放所有的缓存值
	caches map[string]*struct {
		value    string    // 缓存值
		lastRead time.Time // 最后读取的时间
	}

	// cacheTimeout 缓存过期时间
	cacheTimeout time.Duration

	// removeInterval 间隔多少时间检查并移除过期缓存
	removeInterval time.Duration

	// removeTicker 缓存移除定时器
	removeTicker *time.Ticker

	// updateInterval 间隔多少时间刷新缓存
	updateInterval time.Duration

	// updateTicker 缓存刷新定时器
	updateTicker *time.Ticker

	// preCacheReqChan 预缓存请求通道
	preCacheReqChan chan PreCacheReq[T]

	// destroyChan 对象销毁信号通道
	destroyChan chan struct{}

	// mu 并发控制
	mu sync.RWMutex
}

// CacherOption 用于初始化 Cacher 对象
type CacherOption[T any] func(*Cacher[T])

// WithCalcCacheKey 设置 cacheKey 计算方法
func WithCalcCacheKey[T any](f func(T) string) CacherOption[T] {
	return func(c *Cacher[T]) { c.CalcCacheKey = f }
}

// WithRecoverCacheKey 设置 cacheKey 恢复方法
func WithRecoverCacheKey[T any](f func(string) (T, bool)) CacherOption[T] {
	return func(c *Cacher[T]) { c.RecoverCacheKey = f }
}

// WithFetchValue 设置获取最新值的方法, 用于刷新缓存
func WithFetchValue[T any](f func(T) (string, error)) CacherOption[T] {
	return func(c *Cacher[T]) { c.FetchValue = f }
}

// WithUpdateComplete 设置缓存更新完成的统计数据方法
func WithUpdateComplete[T any](f func(success, fail, remove int)) CacherOption[T] {
	return func(c *Cacher[T]) { c.UpdateComplete = f }
}

// WithCacheTimeout 设置缓存过期时间
//
// 某个值超过这个时间没有读取时, 就将其移除
//
// 受缓存过期检查间隔影响, 这个值不是百分百准确的
func WithCacheTimeout[T any](timeout time.Duration) CacherOption[T] {
	return func(c *Cacher[T]) { c.cacheTimeout = timeout }
}

// WithRemoveInterval 设置缓存过期检查间隔
func WithRemoveInterval[T any](removeInterval time.Duration) CacherOption[T] {
	return func(c *Cacher[T]) { c.removeInterval = removeInterval }
}

// WithUpdateInterval 设置缓存刷新时间间隔
func WithUpdateInterval[T any](updateInterval time.Duration) CacherOption[T] {
	return func(c *Cacher[T]) { c.updateInterval = updateInterval }
}

// NewCacher 初始化一个缓存器
func NewCacher[T any](opts ...CacherOption[T]) *Cacher[T] {
	c := new(Cacher[T])
	c.caches = make(map[string]*struct {
		value    string
		lastRead time.Time
	})
	c.preCacheReqChan = make(chan PreCacheReq[T], 100)
	c.removeInterval = time.Minute * 30
	c.updateInterval = time.Hour + time.Minute*30

	for _, opt := range opts {
		opt(c)
	}

	c.removeTicker = time.NewTicker(c.removeInterval)
	c.updateTicker = time.NewTicker(c.updateInterval)
	c.destroyChan = make(chan struct{})

	// 独立线程刷新缓存
	var daemonFunc func()
	daemonFunc = func() {

		var reqPtr *PreCacheReq[T]

		// 防止守护线程崩溃
		defer func() {
			if r := recover(); r != nil {
				log.Printf(colors.ToRed("缓存处理器守护线程出现异常: %v"), r)
				if reqPtr != nil {
					reqPtr.resChan <- struct {
						res string
						err error
					}{err: errors.New("守护线程出现异常")}
					close(reqPtr.resChan)
					reqPtr = nil
				}
				go daemonFunc()
			}
		}()

		for {
			select {
			case <-c.destroyChan:
				return
			case <-c.removeTicker.C:
				c.updateAll(true)
			case <-c.updateTicker.C:
				c.updateAll(false)
			case req := <-c.preCacheReqChan:
				if req.resChan == nil {
					break
				}
				reqPtr = &req
				value, err := c.request(req.T, false)
				req.resChan <- struct {
					res string
					err error
				}{res: value, err: err}
				close(req.resChan)
				reqPtr = nil
			}
		}
	}

	go daemonFunc()

	c.valid = true
	return c
}

// Destroy 销毁缓存器
func (c *Cacher[T]) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.valid {
		return
	}
	c.valid = false

	// 停止守护线程
	c.destroyChan <- struct{}{}
	close(c.destroyChan)

	// 关闭预缓存排队通道
	close(c.preCacheReqChan)

	// 停止定时器
	c.removeTicker.Stop()
	c.updateTicker.Stop()

	// 清空所有缓存
	clear(c.caches)
}

// Request 从缓存器中获取指定入参对应的值,
// 由缓存器自动控制是否使用缓存
//
// 该方法可以被并发调用
func (c *Cacher[T]) Request(params T) (value string, err error) {
	defer func() {
		if err == nil {
			// 更新最后读取时间
			if cache, ok := c.caches[c.CalcCacheKey(params)]; ok {
				cache.lastRead = time.Now()
			}
		}
	}()

	c.mu.RLock()
	if !c.valid {
		c.mu.RUnlock()
		return "", errors.New("缓存器已被销毁")
	}

	// 存在缓存, 直接返回
	cache, ok := c.caches[c.CalcCacheKey(params)]
	c.mu.RUnlock()
	if ok {
		return cache.value, nil
	}

	// 排队请求新频道
	req := PreCacheReq[T]{
		T: params,
		resChan: make(chan struct {
			res string
			err error
		}),
	}
	select {
	case c.preCacheReqChan <- req:
	default:
		return "", fmt.Errorf("服务器繁忙, 请稍后再试")
	}

	// 等待频道信息返回
	result := <-req.resChan
	if result.err != nil {
		return "", fmt.Errorf("刷新缓存异常: %v", result.err)
	}
	return result.res, nil
}

// updateAll 逐个检查所有缓存, 进行淘汰和刷新
//
// removeOnly 指定为 true 时, 只会淘汰不会刷新
func (c *Cacher[T]) updateAll(removeOnly bool) {
	successCnt, failCnt, removeCnt := 0, 0, 0
	defer func() {
		if successCnt <= 0 && failCnt <= 0 && removeCnt <= 0 {
			return
		}
		c.UpdateComplete(successCnt, failCnt, removeCnt)
	}()

	toRemoves := make([]string, 0)
	for key, value := range c.caches {
		if value.lastRead.Add(c.cacheTimeout).Before(time.Now()) {
			// 缓存过期
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		if removeOnly {
			continue
		}

		params, ok := c.RecoverCacheKey(key)
		if !ok {
			// 无效的 key
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		_, err := c.request(params, true)
		if err != nil {
			failCnt++
			toRemoves = append(toRemoves, key)
			removeCnt++
			continue
		}

		successCnt++
	}

	if len(toRemoves) == 0 {
		return
	}

	// 移除过期缓存
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, key := range toRemoves {
		delete(c.caches, key)
	}
}

// requestAndCache 请求缓存值
//
// 如果 forceUpdate 传递为 true, 则忽略缓存, 强制请求最新值并刷新缓存
func (c *Cacher[T]) request(params T, forceUpdate bool) (string, error) {
	cacheKey := c.CalcCacheKey(params)

	// 判断缓存中有没有这个值
	if !forceUpdate {
		c.mu.RLock()
		cache, ok := c.caches[cacheKey]
		c.mu.RUnlock()
		if ok {
			return cache.value, nil
		}
	}

	// 获取最新的值
	value, err := c.FetchValue(params)
	if err != nil {
		return "", fmt.Errorf("无法获取最新值, params: %v, err: %v", params, err)
	}

	// 设置到缓存
	c.mu.Lock()
	defer c.mu.Unlock()
	cache, ok := c.caches[cacheKey]
	if !ok {
		cache = &struct {
			value    string
			lastRead time.Time
		}{}
	}

	cache.value = value
	c.caches[cacheKey] = cache
	return value, nil
}
