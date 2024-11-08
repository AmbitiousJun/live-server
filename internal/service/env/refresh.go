package env

import (
	"errors"
	"log"
	"time"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

// ErrRemoveAndStop 当定时器触发时, 返回该错误表示环境变量已废弃, 需移除
//
// 删除之后定时器停止
var ErrRemoveAndStop = errors.New("remove current env")

// EnvChanger 环境变量更新函数
type EnvChanger func(curVal string) (string, error)

// SetAutoRefresh 设置环境变量定时自动更新
//
// 当定时器触发时, 如果内存中没有指定的 key, 不会触发刷新
func SetAutoRefresh(key string, changer EnvChanger, sched time.Duration) {
	ticker := time.NewTicker(sched)
	go func() {
		for range ticker.C {
			curVal, ok := Get(key)
			if !ok {
				continue
			}
			newVal, err := changer(curVal)

			if err == ErrRemoveAndStop {
				Remove(key)
				log.Printf(colors.ToGray("环境变量已被删除: %s"), key)
				break
			}

			if err != nil {
				log.Printf(colors.ToYellow("环境变量刷新失败, key: %s, err: %v"), key, err)
				continue
			}

			Set(key, newVal)
			log.Printf(colors.ToGreen("环境变量自动刷新, key: %s, value: %s"), key, newVal)
		}
	}()
}
