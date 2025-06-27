package ipresolver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/jsons"
)

// v6Itellyouip 通过 https://www.itellyouip.com 解析 v6
type v6Itellyouip struct{}

func (r *v6Itellyouip) Resolve(ip string) (string, error) {
	u := "https://www.itellyouip.com/ipapi.php?ip=" + ip
	resp, err := https.Get(u).
		AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36").
		Do()
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

// v6Ipwcn 通过 https://ipw.cn 解析 v6
type v6Ipwcn struct{}

// Resolve 解析 ip
func (r *v6Ipwcn) Resolve(ip string) (string, error) {
	u := fmt.Sprintf("https://rest.ipw.cn/api/aw/v1/ipv6?ip=%s&warning=please-direct-use-please-use-ipplus360.com", ip)
	resp, err := https.Get(u).
		AddHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/137.0.0.0 Safari/537.36").
		Do()
	if err != nil {
		return "", fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("错误响应: %v", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取数据异常: %w", err)
	}

	type DataHolder struct {
		Code string `json:"code"`
		Data struct {
			Country  string `json:"country"`
			Prov     string `json:"prov"`
			City     string `json:"city"`
			District string `json:"district"`
			Isp      string `json:"isp"`
		} `json:"data"`
	}
	var d DataHolder
	if err := json.Unmarshal(bodyBytes, &d); err != nil {
		return "", fmt.Errorf("数据解析失败: %w", err)
	}

	res := d.Data.Country + d.Data.Prov + d.Data.City + d.Data.District
	if d.Data.Isp != "" {
		res += "|" + d.Data.Isp
	}

	return res, nil
}
