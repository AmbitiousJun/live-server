package bloom

import (
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/toniphan21/go-bf"
)

// bloomFileName 布隆文件名称
const bloomFileName = "black_ip_bloom.dat"

// FileStorage 通过本地磁盘文件持久化布隆过滤器
type FileStorage struct {
	capacity uint32   // 码位个数
	fileDir  string   // 文件存储路径
	bigInt   *big.Int // 布隆过滤器在内存中的表示
	mu       sync.RWMutex
}

// Set 标记指定比特位
func (fs *FileStorage) Set(index uint32) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	filePath := filepath.Join(fs.fileDir, bloomFileName)
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf(colors.ToRed("布隆过滤器无法生效, 打开本地文件失败: %s, err: %v"), filePath, err)
		return
	}

	// 标记比特位
	if fs.bigInt == nil {
		fs.bigInt = new(big.Int)
		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Printf(colors.ToRed("布隆过滤器无法生效, 读取本地文件失败: %s, err: %v"), filePath, err)
		} else {
			fs.bigInt.SetBytes(bytes)
		}
	}
	fs.bigInt.SetBit(fs.bigInt, int(index), 1)

	// 持久化
	if _, err = file.Write(fs.bigInt.Bytes()); err != nil {
		log.Printf(colors.ToRed("布隆过滤器持久化失败: %s, err: %v"), filePath, err)
	}
}

// Get 获取指定比特位的标记状态
func (fs *FileStorage) Get(index uint32) bool {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	if fs.bigInt != nil {
		return fs.bigInt.Bit(int(index)) == 1
	}

	filePath := filepath.Join(fs.fileDir, bloomFileName)
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Printf(colors.ToYellow("布隆过滤器可能失效, 读取本地文件失败: %s, err: %v"), filePath, err)
		}
		return false
	}

	fs.bigInt = new(big.Int)
	fs.bigInt.SetBytes(bytes)

	return fs.bigInt.Bit(int(index)) == 1
}

// Capacity 获取过滤器容量
func (fs *FileStorage) Capacity() uint32 {
	return fs.capacity
}

// FileStorageFactory 初始化文件存储的工厂
type FileStorageFactory struct {
	FileDir string // 文件持久化目录
}

// Make 创建文件存储过滤器
func (fsf *FileStorageFactory) Make(capacity uint32) (bf.Storage, error) {
	// 检查目录是否存在
	stat, err := os.Stat(fsf.FileDir)
	if err != nil {
		if err := os.MkdirAll(fsf.FileDir, os.ModePerm); err != nil {
			return nil, fmt.Errorf("初始化存储目录失败, dir: %s, err: %v", fsf.FileDir, err)
		}
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("路径被占用: %s", fsf.FileDir)
	}

	return &FileStorage{capacity: capacity, fileDir: fsf.FileDir}, nil
}
