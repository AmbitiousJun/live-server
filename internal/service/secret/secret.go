package secret

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/gin-gonic/gin"
)

const (
	DiskFileName = "secret.txt" // 密钥存放位置
)

// localSecret 存储程序本地的访问密钥
var localSecret string

// Init 初始化密钥, 如果磁盘中已经存在密钥 则进行复用
func Init() error {
	if localSecret != "" {
		return nil
	}

	fp, _ := filepath.Abs(filepath.Join(constant.Dir_DataRoot, DiskFileName))
	defer func() {
		if localSecret != "" {
			log.Printf(colors.ToGreen("secret: %s, 可在路径 [%s] 下进行查看和自定义"), localSecret, fp)
		}
	}()

	// 1 优先读取文件中的密钥
	fileBytes, err := os.ReadFile(fp)
	if err == nil && len(fileBytes) > 0 {
		localSecret = strings.TrimSpace(string(fileBytes))
		return nil
	}

	// 2 生成随机密钥
	localSecret = randomSecret()

	// 3 持久化随机密钥
	if err = os.WriteFile(fp, []byte(localSecret), os.ModePerm); err != nil {
		localSecret = ""
		return fmt.Errorf("持久化密钥失败: %v", err)
	}

	return nil
}

// randomSecret 生成一串随机密钥
func randomSecret() string {
	letters := "1234567890abcdefghijklmnopqrstuvwxyz"
	res := strings.Builder{}

	keyLen := 10
	for keyLen > 0 {
		res.WriteByte(letters[rand.IntN(len(letters))])
		keyLen--
	}

	return res.String()
}

// Get 获取程序本地密钥, 必须先调用 Init 进行初始化
func Get() string {
	return localSecret
}

// Need 对 api 进行 secret 校验
func Need(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1 包未经过初始化, 不进行校验
		if localSecret == "" {
			handler(c)
			return
		}

		// 2 校验密钥
		if c.Query("secret") != localSecret {
			c.String(http.StatusForbidden, "禁止访问")
			return
		}
		handler(c)
	}
}
