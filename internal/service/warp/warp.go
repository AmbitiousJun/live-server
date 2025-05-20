package warp

import (
	"log"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/util/randoms"
)

// Listener ip 检查监听器结构
type Listener struct {

	// CheckIP 检查当前 ip 是否可用, 返回 nil 表示可用
	CheckIP func() error

	// id 监听器唯一 id
	id string

	// busy 标记当前监听器是否正在用于 ip 检测
	busy bool
}

var (
	// allListeners 存放所有注册的监听器
	allListeners = map[string]*Listener{}

	// mu 并发控制锁
	mu sync.Mutex

	// m 全局唯一管理器实例
	m = newManager(10, 15)
)

func init() {
	go func() {
		for r := range m.resultC {
			handleFixResult(r)
		}
	}()
}

// handleFixResult 处理检测结果
func handleFixResult(result fixResult) {
	if result.err != nil {
		log.Printf("ip 修复异常: %v", result.err)
	}

	mu.Lock()
	defer mu.Unlock()

	for _, l := range result.listeners {
		allListeners[l.id].busy = false
	}
}

// RegisterListener 注册监听器, 返回监听器 id
func RegisterListener(l Listener) (id string) {
	mu.Lock()
	defer mu.Unlock()

	// 生成唯一 id
	id = randoms.RandomHex(32)
	l.id = id
	l.busy = false

	// 注册
	allListeners[id] = &l
	return
}

// ReportError 报告 ip 异常
func ReportError(id string) {
	mu.Lock()
	defer mu.Unlock()

	// 获取监听器, 判断是否正在使用
	l, ok := allListeners[id]
	if !ok || l.busy {
		return
	}

	// 加入检测
	l.busy = true
	m.appendListener(*l)
}
