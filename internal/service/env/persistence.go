package env

import (
	"log"
	"os"
	"path/filepath"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (

	// diskFilePath 持久化文件保存目录
	diskFilePath = constant.Dir_DataRoot

	// diskFileName 持久化的文件名称
	diskFileName = "env.json"
)

// keyPair 环境变量键值对
type keyPair struct {
	key   string
	value string
}

// diskPreStoreChan 存储预持久化的环境变量
var diskPreStoreChan = make(chan keyPair, 1000)

func init() {
	envJson := readEnvFromDisk()
	go loopWritingDisk(envJson)
}

// loopWritingDisk 将内存中最新的环境变量持久化到磁盘中,
// 需要运行在独立的 goroutine 上
func loopWritingDisk(envJson *jsons.Item) {
	if envJson == nil || envJson.Type() != jsons.JsonTypeObj {
		envJson = jsons.NewEmptyObj()
	}
	fp := filepath.Join(diskFilePath, diskFileName)
	for pair := range diskPreStoreChan {
		envJson.Put(pair.key, jsons.NewByVal(pair.value))
		if err := os.WriteFile(fp, []byte(envJson.String()), os.ModePerm); err != nil {
			log.Printf(colors.ToRed("环境变量持久化失败: %v"), err)
		}
	}
}

// readEnvFromDisk 将磁盘中的环境变量读取到内存中
func readEnvFromDisk() *jsons.Item {
	// 检查文件是否存在
	if err := os.MkdirAll(diskFilePath, os.ModeDir|os.ModePerm); err != nil {
		log.Printf(colors.ToYellow("初始化 data 目录失败: %v"), err)
		return nil
	}
	fp := filepath.Join(diskFilePath, diskFileName)
	stat, err := os.Stat(fp)
	if err == nil && stat.IsDir() {
		log.Printf(colors.ToYellow("环境变量路径被占用: %s"), fp)
		return nil
	}

	// 读取文件内容
	fileBytes, err := os.ReadFile(fp)
	if err != nil && !os.IsNotExist(err) {
		log.Printf(colors.ToYellow("读取磁盘环境变量出错: %v"), err)
		return nil
	}

	// 序列化 json
	envJson, err := jsons.New(string(fileBytes))
	if err != nil {
		log.Printf(colors.ToYellow("读取磁盘环境变量出错: 非标准的 JSON 格式文件: %v"), err)
		return nil
	}

	// 设置到内存
	envJson.RangeObj(func(key string, valueI *jsons.Item) error {
		value, ok := valueI.Ti().String()
		if !ok {
			// 必须是 string 类型的值才会被读取
			return nil
		}
		globalEnv.Store(key, value)
		log.Printf(colors.ToGreen("恢复环境变量: %s"), key)
		return nil
	})
	return envJson
}
