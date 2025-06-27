package net

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/service/net/ipresolver"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

const (

	// ipMaxRetry 最多连续请求 ip 属地接口失败的次数
	ipMaxRetry = 5
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

	ipAddrCache.Put(ip, jsons.FromValue(ipInfo))

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
		return PrivateIpInfoName, nil
	}

	var resolver ipresolver.R
	// v4
	if pip.To4() != nil {
		resolver = ipresolver.V4()
	} else {
		resolver = ipresolver.V6()
	}

	return resolver.Resolve(ip)
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
