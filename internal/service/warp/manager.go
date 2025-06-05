package warp

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"golang.org/x/sync/errgroup"
)

// fixResult 记录检测结果
type fixResult struct {

	// err 检测异常反馈
	err error

	// listeners 检测时用到的监听器列表
	listeners []Listener
}

// manager 管理监听器 刷新 ip
type manager struct {

	// mu 线程安全同步锁
	mu sync.Mutex

	// listeners 监听器列表
	listeners []Listener

	// resultC 存放检测结果
	resultC chan fixResult

	// working 当前管理器是否处于工作状态
	working bool

	// maxTryPerTime 每次修复最大尝试次数
	maxTryPerTime int

	// httpCli http 客户端
	httpCli *http.Client
}

// newManager 初始化管理器对象
func newManager(resultChanSize, maxTryPerTime int) *manager {
	return &manager{
		resultC:       make(chan fixResult, resultChanSize),
		maxTryPerTime: maxTryPerTime,
		httpCli: &http.Client{
			Transport: &http.Transport{
				DisableKeepAlives:     true,
				TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
				Dial:                  (&net.Dialer{Timeout: time.Second * 10}).Dial,
				ResponseHeaderTimeout: time.Second * 10,
			},
		},
	}
}

// checkIP 检查当前 ip 是否对于所有监听器是否可用
func (m *manager) checkIP() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	g := errgroup.Group{}
	for _, l := range m.listeners {
		l := l
		g.Go(l.CheckIP)
	}
	return g.Wait()
}

// refreshIP 调用脚本 刷新 ip
func (m *manager) refreshIP(execPath string, needV6 bool) error {
	cmd := exec.Command("bash", execPath)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color")
	stdin, _ := cmd.StdinPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("运行脚本失败: %v", err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "1\n")
		// v4 单栈 / v4.v6 双栈
		ipOption := "1\n"
		if needV6 {
			ipOption = "3\n"
		}
		io.WriteString(stdin, ipOption)
	}()

	errBytes, _ := io.ReadAll(stderr)
	if len(errBytes) > 0 {
		return fmt.Errorf("脚本执行失败: %v", string(errBytes))
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("脚本执行失败: %v", err)
	}
	return nil
}

// printCurIP 输出当前的 ip 信息
func (m *manager) printCurIP(needV6 bool) {
	var v4, v6 string
	g := errgroup.Group{}

	inner := func(url, tag string, destPtr *string) func() error {
		return func() error {
			header := http.Header{"User-Agent": []string{"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36"}}
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			req.Header = header
			resp, err := m.httpCli.Do(req)
			if err != nil {
				return fmt.Errorf("请求 %s 地址失败: %v", tag, err)
			}
			defer resp.Body.Close()
			if !https.IsSuccessCode(resp.StatusCode) {
				return fmt.Errorf("请求 %s 地址失败, 请求响应: %s", tag, resp.Status)
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("请求 %s 地址失败, 读取响应失败: %v", tag, err)
			}
			*destPtr = strings.TrimSpace(string(bodyBytes))
			return nil
		}
	}

	g.Go(inner("https://ipinfo.io/ip", "v4", &v4))
	if needV6 {
		g.Go(inner("https://v6.ipinfo.io/ip", "v6", &v6))
	}
	if err := g.Wait(); err != nil {
		log.Printf(colors.ToYellow("获取最新 ip 异常: %v"), err)
	}
	log.Printf(colors.ToGreen("最新 ip 信息 => v4: [%s], v6: [%s]"), v4, v6)
}

// doFix 无限重试, 刷新 ip, 直至所有监听器的检验通过
func (m *manager) doFix(execPath string) {
	defer m.stopFix(nil)
	for curTime := 0; curTime <= m.maxTryPerTime; curTime++ {
		// 1 验证当前所有监听器是否可用
		err := m.checkIP()
		if err == nil {
			log.Println(colors.ToGreen("warp ip 修复完成, 当前可用"))
			return
		}
		log.Printf(colors.ToYellow("warp ip 不可用, 开始进行自动修复, err: %v"), err)

		// 2 随机判断是否获取 needV6 地址
		needV6 := rand.Float64() >= 0.5

		// 3 执行脚本, 刷新 ip
		if err := m.refreshIP(execPath, needV6); err != nil {
			log.Printf(colors.ToRed("warp ip 刷新失败: %v"), err)
			continue
		}

		// 4 输出 v4 v6 信息
		time.Sleep(time.Second * 10)
		m.printCurIP(needV6)
	}
	log.Println(colors.ToPurple("warp ip 刷新重试次数已达上限"))
}

// fix 验证脚本路径后开始异步刷新 ip
func (m *manager) fix() {
	// 获取 warp 执行路径环境变量
	execPath, ok := env.Get(constant.Warp_ExecPathEnvKey)
	if !ok {
		m.stopFix(nil)
		return
	}

	// 检测路径是否是一个可执行文件
	s, err := os.Stat(execPath)
	if err != nil || !s.Mode().IsRegular() || s.Mode().Perm()&0111 == 0 {
		m.stopFix(fmt.Errorf("warp 修复失败, 路径无效: %s", execPath))
		return
	}

	go m.doFix(execPath)
}

// stopFix 停止修复, 恢复状态
func (m *manager) stopFix(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.working = false
	result := fixResult{
		err:       err,
		listeners: m.listeners,
	}
	m.listeners = nil
	m.resultC <- result
}

// appendListener 添加验证监听器
//
// 当管理器正在工作时, 加入当轮校验
// 当管理器为非工作状态时, 开启新一轮校验
// 监听器会在校验完成后自动移除
func (m *manager) appendListener(l Listener) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.listeners = append(m.listeners, l)

	if m.working {
		return
	}
	m.working = true
	m.fix()
}
