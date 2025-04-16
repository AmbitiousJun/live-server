package net

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (

	// ipLookupUrlV4 查询 ipv4 属地信息的请求地址
	ipLookupUrlV4 = `https://ipchaxun.com/${ip}`

	// ipLookupUrlV6 查询 ipv6 属地信息的请求地址
	ipLookupUrlV6 = `https://www.itellyouip.com/ipapi.php?ip=${ip}`

	// ipMaxRetry 最多连续请求 ip 属地接口失败的次数
	ipMaxRetry = 5
)

var (
	ipAsnReg      = regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)<a href="[^"]*" target="_blank" rel="nofollow">(.*)</a>(.*)</span>`)
	ipAsn1Reg     = regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)</span>`)
	ipProviderReg = regexp.MustCompile(`<label><span class="name">运营商：</span><span class="value">(.*)</span></label>`)
)

const (
	IpAddrDiskFileName = "ip_addr.json" // ip 属地信息持久化的文件名
)

var (

	// ipAddrCache 缓存已经查询过属地的 ip
	ipAddrCache *jsons.Item

	// ipAddrCacheOpMutex ip 属地缓存并发控制
	ipAddrCacheOpMutex = sync.RWMutex{}

	// ipPreCacheChan 新的 ip 查询请求放入该通道中
	ipPreCacheChan = make(chan string, 100)
)

func init() {
	readIpAddrFromDisk()
	go func() {
		tryNum := 0
		for ip := range ipPreCacheChan {
			// 缓存过的不重复请求
			if _, ok := getIpAddrCache(ip); ok {
				continue
			}

			// 已经到达最大重试次数
			tryNum++
			if tryNum > ipMaxRetry {
				// 10 分钟后再尝试重新开放 ip 地址获取
				go func() {
					ticker := time.NewTicker(time.Minute * 10)
					defer ticker.Stop()
					<-ticker.C
					tryNum = 0
				}()
				continue
			}

			// 请求 ip 属地
			ipInfo, err := resolveIpAddr(ip)
			if err != nil {
				log.Printf(colors.ToRed("获取 ip 属地失败: %v"), err)
				continue
			}

			setIpAddrCache(ip, ipInfo)
			log.Printf(colors.ToGreen("获取 ip 属地信息成功, %s => %s"), ip, ipInfo)
			tryNum = 0
		}
	}()
}

// getIpAddrCache 从缓存中获取 ip 属地信息
func getIpAddrCache(ip string) (string, bool) {
	ipAddrCacheOpMutex.RLock()
	defer ipAddrCacheOpMutex.RUnlock()

	ip = strings.TrimSpace(ip)
	if cache, ok := ipAddrCache.Attr(ip).String(); ok {
		return cache, true
	}

	return "", false
}

// setIpAddrCache 设置 ip 属地信息到缓存中并持久化到磁盘
func setIpAddrCache(ip, ipInfo string) {
	ipAddrCacheOpMutex.Lock()
	defer ipAddrCacheOpMutex.Unlock()

	ip, ipInfo = strings.TrimSpace(ip), strings.TrimSpace(ipInfo)
	if ip == "" || ipInfo == "" {
		return
	}

	ipAddrCache.Put(ip, jsons.NewByVal(ipInfo))

	fp := filepath.Join(constant.Dir_DataRoot, IpAddrDiskFileName)
	if err := os.WriteFile(fp, []byte(ipAddrCache.String()), os.ModePerm); err != nil {
		log.Printf(colors.ToRed("ip addr 持久化失败: %v, path: %s"), err, fp)
	}
}

// readIpAddrFromDisk 从磁盘中读取 ip 属地信息
func readIpAddrFromDisk() {
	defer func() {
		// 读取失败时, 实例化一个空的 json 对象
		if ipAddrCache == nil {
			ipAddrCache = jsons.NewEmptyObj()
		}
	}()

	if err := os.MkdirAll(constant.Dir_DataRoot, os.ModePerm); err != nil {
		log.Printf(colors.ToRed("检测数据目录失败: %v"), err)
		return
	}

	// 读取文件内容
	fp, _ := filepath.Abs(filepath.Join(constant.Dir_DataRoot, IpAddrDiskFileName))
	fileBytes, err := os.ReadFile(fp)
	if err != nil {
		if !strings.Contains(err.Error(), "no such file or directory") {
			log.Printf(colors.ToRed("读取文件异常: %v, path: %s"), err, fp)
		}
		return
	}

	// 转换成 json
	fileJson, err := jsons.New(string(fileBytes))
	if err != nil {
		log.Printf(colors.ToRed("json 转换失败: %v, path: %s"), err, fp)
		return
	}

	// 校验后赋值到全局变量
	if fileJson.Type() != jsons.JsonTypeObj {
		log.Printf(colors.ToRed("json 类型异常: %s, path: %s"), fileJson.Type(), fp)
		return
	}

	ipAddrCache = fileJson
	log.Printf(colors.ToGreen("成功加载 ip 属地信息: %s"), fp)
}

// resolveIpAddr 解析 ip 地址信息
func resolveIpAddr(ip string) (string, error) {
	ip = strings.TrimSpace(ip)

	// 解析
	pip := net.ParseIP(ip)
	if pip == nil {
		return "", errors.New("解析 ip 失败")
	}

	// 优先判断局域网 ip
	if IsPrivateIp(ip) {
		return "局域网", nil
	}

	// v4
	if pip.To4() != nil {
		return resolveIpv4Addr(ip)
	}

	return resolveIpv6Addr(ip)
}

// resolveIpv4Addr 解析 ipv4 的属地信息
func resolveIpv4Addr(ip string) (string, error) {
	headerWithUA := make(http.Header)
	headerWithUA.Set("User-Agent", "libmpv")
	url := strings.ReplaceAll(ipLookupUrlV4, "${ip}", ip)
	_, resp, err := https.Request(http.MethodGet, url, headerWithUA, nil, true)
	if err != nil {
		return "", err
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	body := string(bodyBytes)

	// 利用正则表达式匹配出目标信息
	if !ipAsnReg.MatchString(body) && !ipAsn1Reg.MatchString(body) {
		return "", fmt.Errorf("解析 ip 属地信息失败: %s", ip)
	}
	sb := strings.Builder{}
	if ipAsnReg.MatchString(body) {
		asns := ipAsnReg.FindStringSubmatch(body)
		sb.WriteString(asns[1])
		sb.WriteString(asns[2])
		sb.WriteString(asns[3])
	} else {
		sb.WriteString(ipAsn1Reg.FindStringSubmatch(body)[1])
	}

	// 补充运营商信息
	if ipProviderReg.MatchString(body) {
		sb.WriteString("|")
		sb.WriteString(ipProviderReg.FindStringSubmatch(body)[1])
	}
	return sb.String(), nil
}

// resolveIpv6Addr 解析 ipv6 的属地信息
func resolveIpv6Addr(ip string) (string, error) {
	u := strings.ReplaceAll(ipLookupUrlV6, "${ip}", ip)
	_, resp, err := https.Request(http.MethodGet, u, nil, nil, true)
	if err != nil {
		return "", fmt.Errorf("请求远程失败: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应体失败: %v", err)
	}

	body := string(bodyBytes)
	resJson, err := jsons.New(body)
	if err != nil {
		return "", fmt.Errorf("解析响应体失败, 原始响应: %s", body)
	}

	if code, ok := resJson.Attr("code").Int(); !ok || code != http.StatusOK {
		return "", fmt.Errorf("获取到非预期的响应: %s", body)
	}

	sb := strings.Builder{}
	local, ok := resJson.Attr("data").Attr("local").String()
	if !ok {
		return "", fmt.Errorf("获取到非预期的响应: %s", body)
	}
	sb.WriteString(local)

	if isp, ok := resJson.Attr("data").Attr("isp").String(); ok && isp != "" {
		sb.WriteString("|")
		sb.WriteString(isp)
	}

	return sb.String(), nil
}

// GetIpAddrInfo 获取 ip 属地信息
func GetIpAddrInfo(ip string) (string, bool) {
	if info, ok := getIpAddrCache(ip); ok {
		return info, true
	}

	select {
	case ipPreCacheChan <- ip:
	default:
	}

	return "", false
}
