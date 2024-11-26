package whitearea

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (
	DiskFileName  = "white_area.json" // 磁盘文件名
	LevelSperator = "/"               // 用于分隔不同的地域层级
)

// whiteAreas 分级存储的地域白名单对象
var whiteAreas *jsons.Item

// whiteAreasOpMutex 并发操作支持
var whiteAreasOpMutex sync.RWMutex

// Init 初始化地域白名单模块
func Init() error {
	whiteAreasOpMutex.Lock()
	defer whiteAreasOpMutex.Unlock()

	fp := filepath.Join(constant.Dir_DataRoot, DiskFileName)
	fileBytes, err := os.ReadFile(fp)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			// 文件存在, 但读取时出现其他的错误
			return fmt.Errorf("读取地域白名单时出错: %v", err)
		}
		// 如果文件不存在, 默认就是一个空对象
		fileBytes = []byte("{}")
	}

	whiteAreas, err = jsons.New(string(fileBytes))
	if err != nil {
		return fmt.Errorf("读取地域白名单时出错: %v", err)
	}
	if whiteAreas.Type() != jsons.JsonTypeObj {
		return fmt.Errorf("错误类型: %s, 地域白名单数据格式必须为对象", whiteAreas.Type())
	}

	log.Printf(colors.ToGreen("成功加载地域白名单信息: %s"), fp)
	return nil
}

// Passable 判断一个地域是否可以通过白名单
func Passable(area string) bool {
	whiteAreasOpMutex.RLock()
	defer whiteAreasOpMutex.RUnlock()

	// 未设置白名单, 默认全放行
	if whiteAreas == nil || whiteAreas.Empty() {
		return true
	}

	// 逐层判断是否放行
	passFlag := false
	curLayer := whiteAreas
	var checkEnd = errors.New("check end")
	var err error
	for {
		continueFind := false
		err = curLayer.RangeObj(func(key string, value *jsons.Item) error {
			if !strings.Contains(area, key) {
				return nil
			}

			// 没有进一步设置子节点, 放行
			if value == nil || value.Type() != jsons.JsonTypeObj || value.Empty() {
				passFlag = true
				return checkEnd
			}

			// 往深一层继续查找
			curLayer = value
			continueFind = true
			return jsons.ErrBreakRange
		})

		if err == checkEnd || !continueFind {
			break
		}
	}

	return passFlag
}

// Set 设置白名单
func Set(area string) {
	area = strings.TrimSpace(area)
	if area == "" || whiteAreas == nil {
		return
	}
	whiteAreasOpMutex.Lock()
	defer whiteAreasOpMutex.Unlock()

	// 逐层设置白名单
	levels := strings.Split(area, LevelSperator)
	curLevel := whiteAreas
	for idx, level := range levels {
		// 到达底层
		if idx == len(levels)-1 {
			curLevel.Put(level, jsons.NewEmptyObj())
			break
		}

		// 成功获取到有效的下一层
		nextLevel, ok := curLevel.Attr(level).Done()
		if ok && nextLevel != nil && nextLevel.Type() == jsons.JsonTypeObj {
			curLevel = nextLevel
			continue
		}

		// 初始化下一层级
		nextLevel = jsons.NewEmptyObj()
		curLevel.Put(level, nextLevel)
		curLevel = nextLevel
	}

	// 持久化
	fp := filepath.Join(constant.Dir_DataRoot, DiskFileName)
	if err := os.MkdirAll(constant.Dir_DataRoot, os.ModePerm); err != nil {
		log.Printf(colors.ToRed("地域白名单数据目录创建失败: %v"), err)
		return
	}
	if err := os.WriteFile(fp, []byte(whiteAreas.String()), os.ModePerm); err != nil {
		log.Printf(colors.ToRed("地域白名单持久化失败: %v"), err)
	}
}

// Del 删除白名单
func Del(area string) {
	area = strings.TrimSpace(area)
	if area == "" || whiteAreas == nil {
		return
	}
	whiteAreasOpMutex.Lock()
	defer whiteAreasOpMutex.Unlock()

	// 遍历删除底层的白名单
	levels := strings.Split(area, LevelSperator)
	curLevel := whiteAreas
	for idx, level := range levels {
		// 到达底层
		if idx == len(levels)-1 {
			curLevel.DelKey(level)
			break
		}

		// 成功获取到有效的下一层
		nextLevel, ok := curLevel.Attr(level).Done()
		if ok && nextLevel != nil && nextLevel.Type() == jsons.JsonTypeObj {
			curLevel = nextLevel
			continue
		}

		// 链路中断, 放弃删除操作
		return
	}

	// 持久化
	fp := filepath.Join(constant.Dir_DataRoot, DiskFileName)
	if err := os.MkdirAll(constant.Dir_DataRoot, os.ModePerm); err != nil {
		log.Printf(colors.ToRed("地域白名单数据目录创建失败: %v"), err)
		return
	}
	if err := os.WriteFile(fp, []byte(whiteAreas.String()), os.ModePerm); err != nil {
		log.Printf(colors.ToRed("地域白名单持久化失败: %v"), err)
	}

}
