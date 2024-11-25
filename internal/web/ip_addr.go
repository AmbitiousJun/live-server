package web

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

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

var (

	// ipAddrCacheMap 缓存已经查询过属地的 ip
	ipAddrCacheMap = sync.Map{}

	// ipPreCacheChan 新的 ip 查询请求放入该通道中
	ipPreCacheChan = make(chan string, 100)
)

func init() {
	go func() {
		tryNum := 0
		for ip := range ipPreCacheChan {
			// 缓存过的不重复请求
			if _, ok := ipAddrCacheMap.Load(ip); ok {
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

			ipAddrCacheMap.Store(ip, ipInfo)
			log.Printf(colors.ToGreen("获取 ip 属地信息成功, %s => %s"), ip, ipInfo)
			tryNum = 0
		}
	}()
}

// resolveIpAddr 解析 ip 地址信息
func resolveIpAddr(ip string) (string, error) {
	ip = strings.TrimSpace(ip)

	// 解析
	pip := net.ParseIP(ip)
	if pip == nil {
		return "", errors.New("解析 ip 失败")
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

	if isp, ok := resJson.Attr("data").Attr("isp").String(); ok {
		sb.WriteString("|")
		sb.WriteString(isp)
	}

	return sb.String(), nil
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
