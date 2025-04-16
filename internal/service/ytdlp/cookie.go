package ytdlp

import (
	"log"
	"os"
	"path/filepath"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

var (

	// CookieFilePath cookie 文件存放路径
	CookieFilePath = filepath.Join(parentPath, "ck.txt")

	// ckExist 程序初始化时验证文件是否存在
	ckExist = false
)

func init() {
	if stat, err := os.Stat(CookieFilePath); err == nil && !stat.IsDir() {
		ckExist = true
		return
	}
	log.Printf(colors.ToYellow("未检测到 cookie 文件: %s, 可能会导致 yt-dlp 解析失败"), CookieFilePath)
}
