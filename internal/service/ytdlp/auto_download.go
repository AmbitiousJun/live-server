package ytdlp

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

const (

	// DefaultReleasePage 默认使用的 yt-dlp 发布页地址
	DefaultReleasePage = "https://ghproxy.cn/https://github.com/yt-dlp/yt-dlp/releases/latest/download"

	// Env_ReleasePage 通过设置环境变量, 覆盖默认的发布页地址
	Env_ReleasePage = "ytdlp_release_page"
)

// arch2ExecNameMap 根据系统的芯片架构, 映射到对应的二进制文件
var arch2ExecNameMap = map[string]string{
	"darwin/amd64":  "yt-dlp_macos",
	"darwin/arm64":  "yt-dlp_macos",
	"windows/386":   "yt-dlp_x86.exe",
	"windows/amd64": "yt-dlp.exe",
	"windows/arm":   "yt-dlp.exe",
	"windows/arm64": "yt-dlp.exe",
	"linux/386":     "yt-dlp_linux",
	"linux/amd64":   "yt-dlp_linux",
	"linux/arm":     "yt-dlp_linux_armv7l",
	"linux/arm64":   "yt-dlp_linux_aarch64",
}

var (
	execOk   = false // 标记二进制文件是否处于就绪状态
	execPath string  // 根据当前系统架构自动生成一个二进制文件地址
)

func init() {
	autoDownloadExec()
}

// autoDownloadExec 自动根据系统架构下载对应版本的 yt-dlp 到数据目录下
//
// 下载失败只会进行日志输出, 不会影响到程序运行
func autoDownloadExec() {
	// 获取系统架构
	gos, garch := runtime.GOOS, runtime.GOARCH

	// 生成二进制文件地址
	execName, ok := arch2ExecNameMap[fmt.Sprintf("%s/%s", gos, garch)]
	if !ok {
		log.Printf("不支持的芯片架构: %s/%s, yt-dlp 相关功能失效", gos, garch)
		return
	}
	execPath = fmt.Sprintf("%s/%s", parentPath, execName)

	// 如果文件不存在, 触发自动下载
	stat, err := os.Stat(execPath)
	if err == nil {
		if stat.IsDir() {
			log.Printf("二进制文件路径被目录占用: %s, 请手动处理后尝试重启服务", execPath)
			return
		}
		execOk = true
		log.Println(colors.ToGreen("yt-dlp 环境检测通过 ✓"))
		return
	}

	log.Println(colors.ToBlue("检测不到 yt-dlp 环境, 即将开始自动下载"))

	if err = os.MkdirAll(parentPath, os.ModePerm); err != nil {
		log.Printf(colors.ToRed("数据目录异常: %s, err: %v"), parentPath, err)
		return
	}

	releasePage, ok := env.Get(Env_ReleasePage)
	if !ok {
		releasePage = DefaultReleasePage
	}
	log.Printf(colors.ToBlue("yt-dlp 下载发布页: %s (可通过环境变量 %s 自定义)"), releasePage, Env_ReleasePage)

	_, resp, err := https.Request(http.MethodGet, releasePage+"/"+execName, nil, nil, true)
	if err != nil {
		log.Printf(colors.ToRed("下载失败: %v"), err)
		return
	}
	defer resp.Body.Close()

	if !https.IsSuccessCode(resp.StatusCode) {
		log.Printf(colors.ToRed("下载失败: %s"), resp.Status)
		return
	}

	execFile, err := os.OpenFile(execPath, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Printf(colors.ToRed("初始化二进制文件路径失败: %s, err: %v"), execPath, err)
		return
	}
	defer execFile.Close()
	io.Copy(execFile, resp.Body)

	// 标记就绪状态
	execOk = true
	log.Printf(colors.ToGreen("yt-dlp 自动下载成功 ✓, 路径: %s"), execPath)
}
