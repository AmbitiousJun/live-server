package handler

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
	"github.com/dop251/goja"
)

// iptv345Params 获取 345 直播地址的请求参数
type iptv345Params struct {
	tid     string // 分类 id
	id      string // 频道 id
	comment string // 备注
}

// iptv345ChMap 将调用方传递的频道名称映射为实际的请求参数
var iptv345ChMap = map[string]iptv345Params{
	"cctv1":  {tid: "ys", id: "1"},
	"cctv2":  {tid: "ys", id: "2"},
	"cctv3":  {tid: "ys", id: "3"},
	"cctv4":  {tid: "ys", id: "4"},
	"cctv5":  {tid: "ys", id: "5"},
	"cctv5p": {tid: "ys", id: "6"},
	"cctv6":  {tid: "ys", id: "7"},
	"cctv7":  {tid: "ys", id: "8"},
	"cctv8":  {tid: "ys", id: "9"},
	"cctv9":  {tid: "ys", id: "10"},
	"cctv10": {tid: "ys", id: "11"},
	"cctv11": {tid: "ys", id: "12"},
	"cctv12": {tid: "ys", id: "13"},
	"cctv13": {tid: "ys", id: "14"},
	"cctv14": {tid: "ys", id: "15"},
	"cctv15": {tid: "ys", id: "16"},
	"cctv16": {tid: "ys", id: "17"},
	"cctv17": {tid: "ys", id: "18"},

	"fhzw": {tid: "gt", id: "1", comment: "凤凰中文"},
	"fhzx": {tid: "gt", id: "2", comment: "凤凰资讯"},
	"fhxg": {tid: "gt", id: "3", comment: "凤凰香港"},
	"ztxw": {tid: "gt", id: "43", comment: "中天新闻"},
	"dsxw": {tid: "gt", id: "44", comment: "东森新闻"},
	"hyxw": {tid: "gt", id: "53", comment: "寰宇新闻"},
}

// iptv345Handler 345 直播处理器
type iptv345Handler struct {
	cacher *resolve.Cacher[iptv345Params]

	// sessionCli 缓存指定 session 的 m3u8
	sessionCli *https.CacheClient

	// cliMaxSessionNum 最多维护的 session 个数
	cliMaxSessionNum int

	// playUriSeg 每个频道可以解析出多条线路, 每条线路之间用这个分隔符进行分割
	playUriSeg string

	// decodeTokenScript 解析 token 脚本
	decodeTokenScript string

	// decodePlayUriScript 解析播放地址脚本
	decodePlayUriScript string

	// requestHeader 请求需要携带的请求头
	requestHeader http.Header
}

func init() {
	ih := new(iptv345Handler)
	ih.initCacher()
	ih.playUriSeg = "{{,}}"

	bytes, _ := base64.StdEncoding.DecodeString(`KGZ1bmN0aW9uKCkgewoKICAgIHN0cmluZyA9ICIke3N0cmluZ30iOwoKICAgIGZ1bmN0aW9uIGRlY29kZShkYXRhKSB7CiAgICAgICAgdmFyIGtleVN0ciA9ICJBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSsvPSI7CiAgICAgICAgdmFyIG8xLCBvMiwgbzMsIGgxLCBoMiwgaDMsIGg0LCBiaXRzLCBpID0gMCwgYWMgPSAwLCBkZWMgPSAiIiwgdG1wX2FyciA9IFtdOwogICAgICAgIGlmICghZGF0YSkgewogICAgICAgICAgICByZXR1cm4gZGF0YQogICAgICAgIH0KICAgICAgICBkYXRhICs9ICIiOwogICAgICAgIGRvIHsKICAgICAgICAgICAgaDEgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDIgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDMgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDQgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgYml0cyA9IGgxIDw8IDE4IHwgaDIgPDwgMTIgfCBoMyA8PCA2IHwgaDQ7CiAgICAgICAgICAgIG8xID0gYml0cyA+PiAxNiAmIDI1NTsKICAgICAgICAgICAgbzIgPSBiaXRzID4+IDggJiAyNTU7CiAgICAgICAgICAgIG8zID0gYml0cyAmIDI1NTsKICAgICAgICAgICAgaWYgKGgzID09IDY0KSB7CiAgICAgICAgICAgICAgICB0bXBfYXJyW2FjKytdID0gU3RyaW5nLmZyb21DaGFyQ29kZShvMSkKICAgICAgICAgICAgfSBlbHNlIHsKICAgICAgICAgICAgICAgIGlmIChoNCA9PSA2NCkgewogICAgICAgICAgICAgICAgICAgIHRtcF9hcnJbYWMrK10gPSBTdHJpbmcuZnJvbUNoYXJDb2RlKG8xLCBvMikKICAgICAgICAgICAgICAgIH0gZWxzZSB7CiAgICAgICAgICAgICAgICAgICAgdG1wX2FyclthYysrXSA9IFN0cmluZy5mcm9tQ2hhckNvZGUobzEsIG8yLCBvMykKICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgfQogICAgICAgIH0gd2hpbGUgKGkgPCBkYXRhLmxlbmd0aCk7CiAgICAgICAgZGVjID0gdG1wX2Fyci5qb2luKCIiKTsKICAgICAgICByZXR1cm4gZGVjCiAgICB9CiAgICB2YXIga2V5ID0gImlwdHYuY29tIjsKICAgIHN0cmluZyA9IGRlY29kZShzdHJpbmcpOwogICAgbGVuID0ga2V5Lmxlbmd0aDsKICAgIGNvZGUgPSAiIjsKICAgIGZvciAoaSA9IDA7IGkgPCBzdHJpbmcubGVuZ3RoOyBpKyspIHsKICAgICAgICBrID0gaSAlIGxlbjsKICAgICAgICBjb2RlICs9IFN0cmluZy5mcm9tQ2hhckNvZGUoc3RyaW5nLmNoYXJDb2RlQXQoaSkgXiBrZXkuY2hhckNvZGVBdChrKSkKICAgIH0KICAgIHZhciB4YWMgPSBkZWNvZGUoY29kZSk7CiAgICB4YWMgPSB1bmVzY2FwZSh4YWMpOwogICAgCiAgICBjb25zdCBtYXRjaGVzID0geGFjLm1hdGNoKC88c2NyaXB0XGJbXj5dKj4oLio/KTxcL3NjcmlwdD4vZ2kpOwogICAgaWYgKG1hdGNoZXMubGVuZ3RoIDwgMykgewogICAgICAgIHJldHVybiB7IGVycjogYHhhYyDmiKrlj5blvILluLg6ICR7bWF0Y2hlc31gIH07CiAgICB9CiAgICBjb25zdCBzY3JpcHQgPSAobWF0Y2hlc1syXSB8fCAnJykucmVwbGFjZUFsbCgvPFwvP3NjcmlwdD4vZywgJycpCiAgICBldmFsKHNjcmlwdCk7CiAgICByZXR1cm4geyBoa2VuLCBoa2VucywgdG9rZW4gfTsKfSkoKQ==`)
	ih.decodeTokenScript = string(bytes)

	bytes, _ = base64.StdEncoding.DecodeString(`KGZ1bmN0aW9uKCkgewoKICAgIGNvbnN0IGhrZW4gPSAiJHtoa2VufSI7CiAgICBjb25zdCBoa2VucyA9ICIke2hrZW5zfSI7CiAgICBjb25zdCB0b2tlbiA9ICIke3Rva2VufSI7CiAgICBsZXQgdXJpID0gIiR7dXJpfSI7CgogICAgZnVuY3Rpb24gYmRlY29kZShkYXRhKSB7CiAgICAgICAgdmFyIGtleVN0ciA9ICJBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSsvPSI7CiAgICAgICAgdmFyIGExLCBhMiwgYTMsIGgxLCBoMiwgaDMsIGg0LCBiaXRzLCBpID0gMCwKICAgICAgICBhYyA9IDAsCiAgICAgICAgZGVjID0gIiIsCiAgICAgICAgdG1wX2FyciA9IFtdOwogICAgICAgIGlmICghZGF0YSkgewogICAgICAgICAgICByZXR1cm4gZGF0YTsKICAgICAgICB9CiAgICAgICAgZGF0YSArPSAnJzsKICAgICAgICBkbyB7CiAgICAgICAgICAgIGgxID0ga2V5U3RyLmluZGV4T2YoZGF0YS5jaGFyQXQoaSsrKSk7CiAgICAgICAgICAgIGgyID0ga2V5U3RyLmluZGV4T2YoZGF0YS5jaGFyQXQoaSsrKSk7CiAgICAgICAgICAgIGgzID0ga2V5U3RyLmluZGV4T2YoZGF0YS5jaGFyQXQoaSsrKSk7CiAgICAgICAgICAgIGg0ID0ga2V5U3RyLmluZGV4T2YoZGF0YS5jaGFyQXQoaSsrKSk7CiAgICAgICAgICAgIGJpdHMgPSBoMSA8PCAxOCB8IGgyIDw8IDEyIHwgaDMgPDwgNiB8IGg0OwogICAgICAgICAgICBhMSA9IGJpdHMgPj4gMTYgJiAweGZmOwogICAgICAgICAgICBhMiA9IGJpdHMgPj4gOCAmIDB4ZmY7CiAgICAgICAgICAgIGEzID0gYml0cyAmIDB4ZmY7CiAgICAgICAgICAgIGlmIChoMyA9PSA2NCkgewogICAgICAgICAgICAgICAgdG1wX2FyclthYysrXSA9IFN0cmluZy5mcm9tQ2hhckNvZGUoYTEpOwogICAgICAgICAgICB9IGVsc2UgaWYgKGg0ID09IDY0KSB7CiAgICAgICAgICAgICAgICB0bXBfYXJyW2FjKytdID0gU3RyaW5nLmZyb21DaGFyQ29kZShhMSwgYTIpOwogICAgICAgICAgICB9IGVsc2UgewogICAgICAgICAgICAgICAgdG1wX2FyclthYysrXSA9IFN0cmluZy5mcm9tQ2hhckNvZGUoYTEsIGEyLCBhMyk7CiAgICAgICAgICAgIH0KICAgICAgICB9IHdoaWxlIChpIDwgZGF0YS5sZW5ndGgpOwogICAgICAgIGRlYyA9IHRtcF9hcnIuam9pbignJyk7CiAgICAgICAgcmV0dXJuIGRlYzsKICAgIH0KCiAgICBmdW5jdGlvbiBiZGVjb2RlYihzdHIsa2V5KSB7CiAgICAgICAgc3RyaW5nID0gYmRlY29kZShzdHIpOwogICAgICAgIGxlbiA9IGtleS5sZW5ndGg7CiAgICAgICAgY29kZSA9ICIiOwogICAgICAgIGZvciAoaSA9IDA7IGkgPCBzdHJpbmcubGVuZ3RoOyBpKyspIHsKICAgICAgICBrID0gaSAlIGxlbjsKICAgICAgICBjb2RlICs9IFN0cmluZy5mcm9tQ2hhckNvZGUoc3RyaW5nLmNoYXJDb2RlQXQoaSkgXiBrZXkuY2hhckNvZGVBdChrKSk7CiAgICAgICAgfQogICAgICAgIHN0cmEgPSBiZGVjb2RlKGNvZGUpOwogICAgICAgIHJldHVybiBzdHJhOwogICAgfQoKICAgIHVyaSA9IHVyaS5zcGxpdCgiIikucmV2ZXJzZSgpLmpvaW4oIiIpOwogICAgdXJpID0gYmRlY29kZWIodXJpLGhrZW4pOwogICAgdXJpID0gdXJpLnJlcGxhY2UoInRva2VuPTEyMyIsICJ0b2tlbj0iK3Rva2VuKTsKICAgIHVyaSA9IHVyaS5yZXBsYWNlKCJ0b2tlbj0iK2hrZW5zLCAidG9rZW49Iit0b2tlbik7CiAgICB1cmkgPSB1cmkucmVwbGFjZShoa2VuLCAiIik7CiAgICByZXR1cm4gdXJpOwp9KSgp`)
	ih.decodePlayUriScript = string(bytes)

	ih.requestHeader = make(http.Header)
	ih.requestHeader.Set("Referer", "https://iptv345.com/")

	ih.sessionCli = https.NewCacheClient(50, time.Minute*30)
	ih.cliMaxSessionNum = 3
	// 定时发送心跳包, 维护 session
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			ih.sendSessionHeartBeatPkg()
		}
	}()

	resolve.RegisterHandler(ih)
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (ih *iptv345Handler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	reqParams, ok := iptv345ChMap[params.ChName]
	if !ok {
		return resolve.HandleResult{}, fmt.Errorf("不支持的频道: [%s]", params.ChName)
	}

	allPlayUri, err := ih.cacher.Request(reqParams)
	if err != nil {
		return resolve.HandleResult{}, fmt.Errorf("解析频道地址失败: %s, ch: [%s]", err, params.ChName)
	}
	playUris := strings.Split(allPlayUri, ih.playUriSeg)
	if len(playUris) == 0 {
		return resolve.HandleResult{}, errors.New("频道地址为空")
	}

	// 根据用户请求的格式返回对应的线路
	idx := 0
	wantIdx, err := strconv.Atoi(params.Format)
	if err == nil {
		wantIdx--
		if wantIdx < 0 || wantIdx >= len(playUris) {
			return resolve.HandleResult{}, fmt.Errorf("format 指定错误, 可配置值: [%d ~ %d]", 1, len(playUris))
		}
		idx = wantIdx
	}

	sessionM3U8, resp, err := ih.sessionCli.Request(http.MethodGet, playUris[idx], ih.requestHeader.Clone(), nil, true)
	if err != nil || !https.IsSuccessCode(resp.StatusCode) {
		return resolve.HandleResult{}, fmt.Errorf("获取 m3u8 session 失败: %v", err)
	}
	resp.Body.Close()
	ih.removeRedundantSessionIfMaximum(sessionM3U8)

	params.Headers = ih.requestHeader.Clone()
	return resolve.M3U8Result(sessionM3U8, params)
}

// Name 处理器名称
func (ih *iptv345Handler) Name() string {
	return "345"
}

// HelpDoc 处理器说明文档
func (ih *iptv345Handler) HelpDoc() string {
	sb := strings.Builder{}
	sb.WriteString("\n1. 该处理器仅限自用, 不适合分享")
	sb.WriteString("\n2. 如果想通过直连播放（重定向）, 请使用可自定请求头的播放器（酷9、天光云影等）设置 Referer 头: [https://iptv345.com/] 后才可正常播放")
	sb.WriteString("\n3. 如果是切片代理模式, 则无需指定请求头")
	sb.WriteString("\n4. 连续播放一段时间可能会发生断流, 需要重新加载")

	chs := []string{}
	for k, v := range iptv345ChMap {
		cur := k
		if v.comment != "" {
			cur += "（" + v.comment + "）"
		}
		chs = append(chs, cur)
	}
	sb.WriteString("\n5. 支持频道: " + strings.Join(chs, "、"))

	sb.WriteString("\n6. 可通过 query 参数 [format] 切换线路, 可配置值 [1 ~ ?], 默认请求线路 1, 每个频道的最多线路不同, 可自行测试, 超出最大线路时会直接报错提示")

	return sb.String()
}

// SupportProxy 是否支持 m3u 代理
//
// 如果返回 true, 会自动在帮助文档中加入标记
func (ih *iptv345Handler) SupportM3UProxy() bool {
	return true
}

func (ih *iptv345Handler) initCacher() {
	ih.cacher = resolve.NewCacher(
		resolve.WithCacheTimeout[iptv345Params](time.Hour),
		resolve.WithRemoveInterval[iptv345Params](time.Minute*10),
		resolve.WithUpdateInterval[iptv345Params](time.Hour+time.Minute*30),

		resolve.WithUpdateComplete[iptv345Params](func(success, fail, remove int) {
			log.Printf(colors.ToGreen("345 缓存更新完成, 成功: %d, 失败: %d, 移除: %d"), success, fail, remove)
		}),

		resolve.WithCalcCacheKey(func(p iptv345Params) string {
			return p.tid + ":" + p.id
		}),

		resolve.WithRecoverCacheKey(func(s string) (iptv345Params, bool) {
			splits := strings.Split(s, ":")
			if len(splits) != 2 {
				return iptv345Params{}, false
			}
			return iptv345Params{tid: splits[0], id: splits[1]}, true
		}),

		resolve.WithFetchValue(func(p iptv345Params) (string, error) {
			str2Decode, playUris, err := ih.fetchOriginData(p.tid, p.id)
			if err != nil {
				return "", fmt.Errorf("解析网页源代码异常: %v", err)
			}

			secrets, err := ih.decodeStringSecrets(str2Decode)
			if err != nil {
				return "", fmt.Errorf("解析 token 异常: %v", err)
			}

			hken := secrets["hken"].(string)
			hkens := secrets["hkens"].(string)
			token := secrets["token"].(string)

			rawPlayUris := make([]string, len(playUris))
			for i, playUri := range playUris {
				res, err := ih.recoverPlayUri(hken, hkens, token, playUri)
				if err != nil {
					log.Printf(colors.ToYellow("恢复频道播放地址失败: %v"), err)
					continue
				}
				rawPlayUris[i] = res
			}
			return strings.Join(rawPlayUris, ih.playUriSeg), nil
		}),
	)
}

// fetchOriginData 获取网页中的 string 混淆代码变量值
func (ih *iptv345Handler) fetchOriginData(tid, id string) (str2Decode string, playUris []string, fe error) {
	pageUrl := fmt.Sprintf("https://iptv345.com/?act=play&tid=%s&id=%s", tid, id)
	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	_, res, err := https.Request(http.MethodGet, pageUrl, header, nil, false)
	if err != nil {
		fe = fmt.Errorf("请求页面失败: %v", err)
		return
	}
	defer res.Body.Close()
	pageBytes, err := io.ReadAll(res.Body)
	if err != nil {
		fe = fmt.Errorf("读取页面数据失败: %v", err)
		return
	}
	pageCode := string(pageBytes)

	// 2 解析出 string 变量
	varNameGetter := regexp.MustCompile(`var string=([^;]+);`)
	if !varNameGetter.MatchString(pageCode) {
		fe = fmt.Errorf("未找到 string 的混淆变量名称, 源代码: %s", pageCode)
		return
	}
	varName := varNameGetter.FindStringSubmatch(pageCode)[1]

	varValueGetter := regexp.MustCompile(fmt.Sprintf(`var %s="([^"]+)";`, varName))
	if !varValueGetter.MatchString(pageCode) {
		fe = fmt.Errorf("未找到 string 变量值, 源代码: %s", pageCode)
		return
	}
	stringValue := varValueGetter.FindStringSubmatch(pageCode)[1]

	varReverseGetter := regexp.MustCompile(fmt.Sprintf(`%s = %s\.split\(""\)\.reverse\(\)\.join\(""\);`, varName, varName))
	if varReverseGetter.MatchString(pageCode) {
		stringValue = strs.ReverseString(stringValue)
	}

	// 3 解析原始的播放地址
	var uris []string
	uriGetter := regexp.MustCompile(`<option value="([^"]+)">[^<]*<\/option>`)
	if !uriGetter.MatchString(pageCode) {
		fe = fmt.Errorf("获取不到原始播放地址, 源代码: %s", pageCode)
		return
	}
	matches := uriGetter.FindAllStringSubmatch(pageCode, -1)
	if len(matches) == 0 {
		fe = fmt.Errorf("获取不到原始播放地址, 源代码: %s", pageCode)
		return
	}
	for _, match := range matches {
		uris = append(uris, match[1])
	}

	str2Decode = stringValue
	playUris = uris
	return
}

// decodeStringSecrets 解码 string 混淆代码变量, 得到几个关键密钥 hken, hkens, token
func (ih *iptv345Handler) decodeStringSecrets(encodeStr string) (map[string]interface{}, error) {
	resolveStringCode := strings.Replace(ih.decodeTokenScript, "${string}", encodeStr, 1)

	// 执行 js 代码, 解混淆
	vm := goja.New()
	res, err := vm.RunString(resolveStringCode)
	if err != nil {
		return nil, fmt.Errorf("执行解析 string 的 js 代码失败: %v", err)
	}
	secrets, ok := res.Export().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("解析 string 代码失败: 非 map 类型, 原始响应: %v", res)
	}

	// 4 如果包含 err 字段, 代表出异常
	if val, ok := secrets["err"]; ok {
		return nil, fmt.Errorf("%v", val)
	}

	return secrets, nil
}

// recoverPlayUri 恢复播放地址明文
func (ih *iptv345Handler) recoverPlayUri(hken, hkens, token, originPlayUri string) (string, error) {
	resolvePlayUriCode := ih.decodePlayUriScript
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${hken}", hken, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${hkens}", hkens, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${token}", token, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${uri}", originPlayUri, 1)

	vm := goja.New()
	res, err := vm.RunString(resolvePlayUriCode)
	if err != nil {
		return "", fmt.Errorf("执行 js 异常: %v", err)
	}
	return res.String(), nil
}

// sendSessionHeartBeatPkg 发送 session 心跳包, 防止 session 被关闭
func (ih *iptv345Handler) sendSessionHeartBeatPkg() {
	urls := ih.sessionCli.GetAllCacheUrls()
	for _, url := range urls {
		_, resp, err := https.Request(http.MethodGet, url, ih.requestHeader.Clone(), nil, true)
		if err != nil {
			log.Printf(colors.ToYellow("发送 345 心跳包异常: %v"), err)
			continue
		}
		resp.Body.Close()
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusNotFound {
			log.Println(colors.ToYellow("345 session 过期, 进行移除"))
			ih.sessionCli.RemoveUrlCache(url)
		}
	}
}

// removeRedundantSessionIfMaximum 如果内存中的会话满了, 就淘汰其他多余的
func (ih *iptv345Handler) removeRedundantSessionIfMaximum(exclude string) {
	urls := ih.sessionCli.GetAllCacheUrls()
	removeCnt := len(urls) - ih.cliMaxSessionNum
	if removeCnt <= 0 {
		return
	}
	for _, url := range urls {
		if url == exclude {
			continue
		}
		ih.sessionCli.RemoveUrlCache(url)
		removeCnt--
		if removeCnt <= 0 {
			return
		}
	}
}
