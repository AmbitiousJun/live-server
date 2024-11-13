package web

import (
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
)

const (

	// ipLookupUrl 查询 ip 属地信息的请求地址
	ipLookupUrl = `https://ipchaxun.com/${ip}`

	// ipMaxRetry 最多连续请求 ip 属地接口失败的次数
	ipMaxRetry = 5
)

var (
	ipAsnReg      = regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)<a href="[^"]*" target="_blank" rel="nofollow">(.*)</a>(.*)</span>`)
	ipProviderReg = regexp.MustCompile(`<label><span class="name">运营商：</span><span class="value">(.*)</span></label>`)
)

var (

	// ipAddrCacheMap 缓存已经查询过属地的 ip
	ipAddrCacheMap = sync.Map{}

	// ipPreCacheChan 新的 ip 查询请求放入该通道中
	ipPreCacheChan = make(chan string, 100)
)

func init() {
	go func() {
		tryNum := 0
		headerWithUA := make(http.Header)
		headerWithUA.Set("User-Agent", "libmpv")
		for ip := range ipPreCacheChan {
			// 缓存过的不重复请求
			if _, ok := ipAddrCacheMap.Load(ip); ok {
				continue
			}
			// 只查询 ipv4 地址的信息
			pip := net.ParseIP(ip)
			if pip == nil || pip.To4() == nil {
				continue
			}
			// 已经到达最大重试次数
			tryNum++
			if tryNum > ipMaxRetry {
				continue
			}

			// 请求 ip 属地
			url := strings.ReplaceAll(ipLookupUrl, "${ip}", ip)
			_, resp, err := https.Request(http.MethodGet, url, headerWithUA, nil, true)
			if err != nil {
				log.Printf(colors.ToRed("获取 ip 属地失败: %v"), err)
				continue
			}
			bodyBytes, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				log.Printf(colors.ToRed("获取 ip 属地失败: %v"), err)
				continue
			}
			body := string(bodyBytes)

			// 利用正则表达式匹配出目标信息
			if !ipAsnReg.MatchString(body) {
				log.Printf(colors.ToRed("解析 ip 属地失败, 远程原始响应: %s"), body)
				continue
			}
			sb := strings.Builder{}
			asns := ipAsnReg.FindStringSubmatch(body)
			sb.WriteString(asns[1])
			sb.WriteString(asns[2])
			sb.WriteString(asns[3])
			if ipProviderReg.MatchString(body) {
				sb.WriteString("|")
				sb.WriteString(ipProviderReg.FindStringSubmatch(body)[1])
			}

			ipAddrCacheMap.Store(ip, sb.String())
			log.Printf(colors.ToGreen("获取 ip 属地信息成功, %s => %s"), ip, sb.String())
			tryNum = 0
		}
	}()
}

// GetIpAddrInfo 获取 ip 属地信息
func GetIpAddrInfo(ip string) (string, bool) {
	if info, ok := ipAddrCacheMap.Load(ip); ok {
		return info.(string), true
	}

	select {
	case ipPreCacheChan <- ip:
	default:
	}

	return "", false
}
