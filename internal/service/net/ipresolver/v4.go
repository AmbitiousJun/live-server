package ipresolver

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/util/https"
)

// v4Ipchaxun 调用 https://ipchaxun.com 解析 v4 地址
type v4Ipchaxun struct{}

// Resolve
func (r *v4Ipchaxun) Resolve(ip string) (string, error) {
	url := "https://ipchaxun.com/" + ip
	resp, err := https.Get(url).
		AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36").
		Do()
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
	ipAsnReg := regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)<a href="[^"]*" target="_blank" rel="nofollow">(.*)</a>(.*)</span>`)
	ipAsn1Reg := regexp.MustCompile(`<span class="name">归属地：</span><span class="value">(.*)</span>`)
	ipProviderReg := regexp.MustCompile(`<label><span class="name">运营商：</span><span class="value">(.*)</span></label>`)

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
