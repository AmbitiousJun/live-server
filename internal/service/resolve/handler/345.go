package handler

import (
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
	"github.com/AmbitiousJun/live-server/internal/util/base64s"
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

// iptv345Handler 345 直播处理器
type iptv345Handler struct {
	cacher *resolve.Cacher[iptv345Params]

	// serverHost 服务器地址
	serverHost string

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

	// channels 支持的频道列表
	channels map[string]iptv345Params
}

func init() {
	ih := new(iptv345Handler)
	ih.initCacher()
	ih.serverHost = base64s.MustDecodeString("aHR0cHM6Ly9pcHR2MzQ1LmNvbQ==")
	ih.playUriSeg = "{{,}}"

	ih.decodeTokenScript = base64s.MustDecodeString(`KGZ1bmN0aW9uKCkgewoKICAgIHN0cmluZyA9ICIke3N0cmluZ30iOwoKICAgIGZ1bmN0aW9uIGRlY29kZShkYXRhKSB7CiAgICAgICAgdmFyIGtleVN0ciA9ICJBQkNERUZHSElKS0xNTk9QUVJTVFVWV1hZWmFiY2RlZmdoaWprbG1ub3BxcnN0dXZ3eHl6MDEyMzQ1Njc4OSsvPSI7CiAgICAgICAgdmFyIG8xLCBvMiwgbzMsIGgxLCBoMiwgaDMsIGg0LCBiaXRzLCBpID0gMCwgYWMgPSAwLCBkZWMgPSAiIiwgdG1wX2FyciA9IFtdOwogICAgICAgIGlmICghZGF0YSkgewogICAgICAgICAgICByZXR1cm4gZGF0YQogICAgICAgIH0KICAgICAgICBkYXRhICs9ICIiOwogICAgICAgIGRvIHsKICAgICAgICAgICAgaDEgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDIgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDMgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgaDQgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgICAgICAgYml0cyA9IGgxIDw8IDE4IHwgaDIgPDwgMTIgfCBoMyA8PCA2IHwgaDQ7CiAgICAgICAgICAgIG8xID0gYml0cyA+PiAxNiAmIDI1NTsKICAgICAgICAgICAgbzIgPSBiaXRzID4+IDggJiAyNTU7CiAgICAgICAgICAgIG8zID0gYml0cyAmIDI1NTsKICAgICAgICAgICAgaWYgKGgzID09IDY0KSB7CiAgICAgICAgICAgICAgICB0bXBfYXJyW2FjKytdID0gU3RyaW5nLmZyb21DaGFyQ29kZShvMSkKICAgICAgICAgICAgfSBlbHNlIHsKICAgICAgICAgICAgICAgIGlmIChoNCA9PSA2NCkgewogICAgICAgICAgICAgICAgICAgIHRtcF9hcnJbYWMrK10gPSBTdHJpbmcuZnJvbUNoYXJDb2RlKG8xLCBvMikKICAgICAgICAgICAgICAgIH0gZWxzZSB7CiAgICAgICAgICAgICAgICAgICAgdG1wX2FyclthYysrXSA9IFN0cmluZy5mcm9tQ2hhckNvZGUobzEsIG8yLCBvMykKICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgfQogICAgICAgIH0gd2hpbGUgKGkgPCBkYXRhLmxlbmd0aCk7CiAgICAgICAgZGVjID0gdG1wX2Fyci5qb2luKCIiKTsKICAgICAgICByZXR1cm4gZGVjCiAgICB9CgogICAgZnVuY3Rpb24gZXh0cmFjdFRva2VuTmFtZXMoeGFjID0gJycpIHsKICAgICAgICBjb25zdCB1cmlOYW1lUmVzID0geGFjLm1hdGNoKC9zdGFydFBsYXllclwoKFteXCldKilcKS9pKTsKICAgICAgICBpZiAoIXVyaU5hbWVSZXMgfHwgdXJpTmFtZVJlcy5sZW5ndGggPCAyKSB7CiAgICAgICAgICAgIHJldHVybiB7IGVycjogYHhhYzogdXJpTmFtZSDmiKrlj5blvILluLg6ICR7dXJpTmFtZVJlc31gIH07CiAgICAgICAgfQogICAgICAgIGNvbnN0IHVyaSA9IHVyaU5hbWVSZXNbMV07CgogICAgICAgIGNvbnN0IGhrZW5SZXMgPSB4YWMubWF0Y2gobmV3IFJlZ0V4cChgJHt1cml9ID0gW15cXChdKlxcKCR7dXJpfSwoW15cXCldKilcXClgLCAnaScpKTsKICAgICAgICBpZiAoIWhrZW5SZXMgfHwgaGtlblJlcy5sZW5ndGggPCAyKSB7CiAgICAgICAgICAgIHJldHVybiB7IGVycjogYHhhYzogaGtlbiDmiKrlj5blvILluLg6ICR7aGtlblJlc31gIH07CiAgICAgICAgfQoKICAgICAgICBjb25zdCB0b2tlblJlcyA9IHhhYy5tYXRjaChuZXcgUmVnRXhwKGAke3VyaX0gPSAke3VyaX0ucmVwbGFjZVxcKCJ0b2tlbj0iXFwrKFteLF0qKSwgInRva2VuPSJcXCsoW15cXCldKilcXClgLCAnaScpKTsKICAgICAgICBpZiAoIXRva2VuUmVzIHx8IHRva2VuUmVzLmxlbmd0aCA8IDMpIHsKICAgICAgICAgICAgcmV0dXJuIHsgZXJyOiBgeGFjOiB0b2tlbiDmiKrlj5blvILluLg6ICR7dG9rZW5SZXN9YCB9OwogICAgICAgIH0KCiAgICAgICAgcmV0dXJuIHsgaGtlbjogaGtlblJlc1sxXSwgaGtlbnM6IHRva2VuUmVzWzFdLCB0b2tlbjogdG9rZW5SZXNbMl0gfTsKICAgIH0KCiAgICBmdW5jdGlvbiBleHRyYWN0SGtlbkV4dHJhS2V5KHhhYyA9ICcnKSB7CiAgICAgICAgY29uc3Qga2V5UmVzID0geGFjLm1hdGNoKC9rZXkgPSBrZXlcKyIoW14iXSopIi9pKTsKICAgICAgICBpZiAoIWtleVJlcyB8fCBrZXlSZXMubGVuZ3RoIDwgMikgewogICAgICAgICAgICByZXR1cm4geyBlcnI6IGBoa2VuIGV4dHJhIGtleSDmiKrlj5blvILluLg6ICR7a2V5UmVzfWAgfTsKICAgICAgICB9CiAgICAgICAgcmV0dXJuIHsgdmFsdWU6IGtleVJlc1sxXSB9OwogICAgfQoKICAgIHZhciBrZXkgPSAiaXB0di5jb20iOwogICAgc3RyaW5nID0gZGVjb2RlKHN0cmluZyk7CiAgICBsZW4gPSBrZXkubGVuZ3RoOwogICAgY29kZSA9ICIiOwogICAgZm9yIChpID0gMDsgaSA8IHN0cmluZy5sZW5ndGg7IGkrKykgewogICAgICAgIGsgPSBpICUgbGVuOwogICAgICAgIGNvZGUgKz0gU3RyaW5nLmZyb21DaGFyQ29kZShzdHJpbmcuY2hhckNvZGVBdChpKSBeIGtleS5jaGFyQ29kZUF0KGspKQogICAgfQogICAgdmFyIHhhYyA9IGRlY29kZShjb2RlKTsKICAgIHhhYyA9IHVuZXNjYXBlKHhhYyk7CiAgICAKICAgIGNvbnN0IG1hdGNoZXMgPSB4YWMubWF0Y2goLzxzY3JpcHRcYltePl0qPiguKj8pPFwvc2NyaXB0Pi9naSk7CiAgICBsZXQgc2NyaXB0RWxtID0gbnVsbDsKICAgIGZvciAoY29uc3QgbWF0Y2ggb2YgbWF0Y2hlcykgewogICAgICAgIGlmIChtYXRjaC5pbmNsdWRlcygnam9pbicpKSB7CiAgICAgICAgICAgIHNjcmlwdEVsbSA9IG1hdGNoOwogICAgICAgICAgICBicmVhazsKICAgICAgICB9CiAgICB9CiAgICBpZiAoIXNjcmlwdEVsbSkgewogICAgICAgIHJldHVybiB7IGVycjogYHhhYyDmiKrlj5blvILluLg6ICR7bWF0Y2hlc31gIH07CiAgICB9CiAgICBpZiAoZmFsc2UpIHsKICAgICAgICByZXR1cm4geyBlcnI6IHhhYyB9OwogICAgfQoKICAgIGNvbnN0IHRva2VuTmFtZXMgPSBleHRyYWN0VG9rZW5OYW1lcyh4YWMpIHx8IHt9OwogICAgaWYgKHRva2VuTmFtZXMuZXJyKSB7CiAgICAgICAgcmV0dXJuIHsgZXJyOiBg6Kej5p6QIHRva2VuIOWPmOmHj+WQjeensOWksei0pTogJHt0b2tlbk5hbWVzLmVycn1gIH07CiAgICB9CgogICAgY29uc3QgZXh0cmFLZXkgPSBleHRyYWN0SGtlbkV4dHJhS2V5KHhhYykgfHwge307CiAgICBpZiAoZXh0cmFLZXkuZXJyKSB7CiAgICAgICAgcmV0dXJuIHsgZXJyOiBg6Kej5p6QIGV4dHJhIGtleSDlpLHotKU6ICR7ZXh0cmFLZXkuZXJyfWAgfTsKICAgIH0KCiAgICBjb25zdCBzY3JpcHQgPSBzY3JpcHRFbG0ucmVwbGFjZUFsbCgvPFwvP3NjcmlwdD4vZywgJycpOwogICAgZXZhbChzY3JpcHQpOwogICAgcmV0dXJuIHsgaGtlbjogZXZhbCh0b2tlbk5hbWVzWydoa2VuJ10pLCBoa2VuczogZXZhbCh0b2tlbk5hbWVzWydoa2VucyddKSwgdG9rZW46IGV2YWwodG9rZW5OYW1lc1sndG9rZW4nXSksIGhrZW5FeHRyYUtleTogZXh0cmFLZXkudmFsdWUgfTsKfSkoKQ==`)
	ih.decodePlayUriScript = base64s.MustDecodeString(`KGZ1bmN0aW9uICgpIHsKICBjb25zdCBoa2VuID0gIiR7aGtlbn0iOwogIGNvbnN0IGhrZW5zID0gIiR7aGtlbnN9IjsKICBjb25zdCB0b2tlbiA9ICIke3Rva2VufSI7CiAgY29uc3QgaGtlbkV4dHJhS2V5ID0gIiR7aGtlbkV4dHJhS2V5fSI7CiAgbGV0IHVyaSA9ICIke3VyaX0iOwoKICBmdW5jdGlvbiBiZGVjb2RlKGRhdGEpIHsKICAgIHZhciBrZXlTdHIgPQogICAgICAiQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVphYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ejAxMjM0NTY3ODkrLz0iOwogICAgdmFyIGExLAogICAgICBhMiwKICAgICAgYTMsCiAgICAgIGgxLAogICAgICBoMiwKICAgICAgaDMsCiAgICAgIGg0LAogICAgICBiaXRzLAogICAgICBpID0gMCwKICAgICAgYWMgPSAwLAogICAgICBkZWMgPSAiIiwKICAgICAgdG1wX2FyciA9IFtdOwogICAgaWYgKCFkYXRhKSB7CiAgICAgIHJldHVybiBkYXRhOwogICAgfQogICAgZGF0YSArPSAiIjsKICAgIGRvIHsKICAgICAgaDEgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgaDIgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgaDMgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgaDQgPSBrZXlTdHIuaW5kZXhPZihkYXRhLmNoYXJBdChpKyspKTsKICAgICAgYml0cyA9IChoMSA8PCAxOCkgfCAoaDIgPDwgMTIpIHwgKGgzIDw8IDYpIHwgaDQ7CiAgICAgIGExID0gKGJpdHMgPj4gMTYpICYgMHhmZjsKICAgICAgYTIgPSAoYml0cyA+PiA4KSAmIDB4ZmY7CiAgICAgIGEzID0gYml0cyAmIDB4ZmY7CiAgICAgIGlmIChoMyA9PSA2NCkgewogICAgICAgIHRtcF9hcnJbYWMrK10gPSBTdHJpbmcuZnJvbUNoYXJDb2RlKGExKTsKICAgICAgfSBlbHNlIGlmIChoNCA9PSA2NCkgewogICAgICAgIHRtcF9hcnJbYWMrK10gPSBTdHJpbmcuZnJvbUNoYXJDb2RlKGExLCBhMik7CiAgICAgIH0gZWxzZSB7CiAgICAgICAgdG1wX2FyclthYysrXSA9IFN0cmluZy5mcm9tQ2hhckNvZGUoYTEsIGEyLCBhMyk7CiAgICAgIH0KICAgIH0gd2hpbGUgKGkgPCBkYXRhLmxlbmd0aCk7CiAgICBkZWMgPSB0bXBfYXJyLmpvaW4oIiIpOwogICAgcmV0dXJuIGRlYzsKICB9CgogIGZ1bmN0aW9uIGJkZWNvZGViKHN0ciwga2V5KSB7CiAgICBzdHJpbmcgPSBiZGVjb2RlKHN0cik7CiAgICBrZXkgKz0gaGtlbkV4dHJhS2V5OwogICAgbGVuID0ga2V5Lmxlbmd0aDsKICAgIGNvZGUgPSAiIjsKICAgIGZvciAoaSA9IDA7IGkgPCBzdHJpbmcubGVuZ3RoOyBpKyspIHsKICAgICAgayA9IGkgJSBsZW47CiAgICAgIGNvZGUgKz0gU3RyaW5nLmZyb21DaGFyQ29kZShzdHJpbmcuY2hhckNvZGVBdChpKSBeIGtleS5jaGFyQ29kZUF0KGspKTsKICAgIH0KICAgIHN0cmEgPSBiZGVjb2RlKGNvZGUpOwogICAgcmV0dXJuIHN0cmE7CiAgfQoKICB1cmkgPSB1cmkuc3BsaXQoIiIpLnJldmVyc2UoKS5qb2luKCIiKTsKICB1cmkgPSBiZGVjb2RlYih1cmksIGhrZW4pOwogIHVyaSA9IHVyaS5yZXBsYWNlKCJ0b2tlbj0xMjMiLCAidG9rZW49IiArIHRva2VuKTsKICB1cmkgPSB1cmkucmVwbGFjZSgidG9rZW49IiArIGhrZW5zLCAidG9rZW49IiArIHRva2VuKTsKICB1cmkgPSB1cmkucmVwbGFjZShoa2VuLCAiIik7CiAgcmV0dXJuIHVyaTsKfSkoKTsK`)

	ih.requestHeader = make(http.Header)
	ih.requestHeader.Set("Referer", ih.serverHost+"/")
	ih.requestHeader.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")
	ih.requestHeader.Set("Origin", "null")

	ih.sessionCli = https.NewCacheClient(50, time.Minute*30)
	ih.cliMaxSessionNum = 3
	// 定时发送心跳包, 维护 session
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			ih.sendSessionHeartBeatPkg()
		}
	}()

	ih.channels = map[string]iptv345Params{
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

		"hunws":    {tid: "ws", id: "1", comment: "湖南卫视"},
		"jsws":     {tid: "ws", id: "2", comment: "江苏卫视"},
		"zjws":     {tid: "ws", id: "3", comment: "浙江卫视"},
		"dfws":     {tid: "ws", id: "4", comment: "东方卫视"},
		"bjws":     {tid: "ws", id: "5", comment: "北京卫视"},
		"szws":     {tid: "ws", id: "6", comment: "深圳卫视"},
		"gdws":     {tid: "ws", id: "7", comment: "广东卫视"},
		"ahws":     {tid: "ws", id: "8", comment: "安徽卫视"},
		"dnws":     {tid: "ws", id: "9", comment: "东南卫视"},
		"hebws":    {tid: "ws", id: "10", comment: "河北卫视"},
		"hljws":    {tid: "ws", id: "11", comment: "黑龙江卫视"},
		"hubws":    {tid: "ws", id: "12", comment: "湖北卫视"},
		"jxws":     {tid: "ws", id: "13", comment: "江西卫视"},
		"lnws":     {tid: "ws", id: "14", comment: "辽宁卫视"},
		"hainws":   {tid: "ws", id: "15", comment: "海南卫视"},
		"sdws":     {tid: "ws", id: "16", comment: "山东卫视"},
		"scws":     {tid: "ws", id: "17", comment: "四川卫视"},
		"tjws":     {tid: "ws", id: "18", comment: "天津卫视"},
		"cqws":     {tid: "ws", id: "19", comment: "重庆卫视"},
		"gzws":     {tid: "ws", id: "20", comment: "贵州卫视"},
		"jlws":     {tid: "ws", id: "21", comment: "吉林卫视"},
		"gxws":     {tid: "ws", id: "22", comment: "广西卫视"},
		"henws":    {tid: "ws", id: "23", comment: "河南卫视"},
		"gsws":     {tid: "ws", id: "24", comment: "甘肃卫视"},
		"qhws":     {tid: "ws", id: "25", comment: "青海卫视"},
		"ynws":     {tid: "ws", id: "26", comment: "云南卫视"},
		"nmgws":    {tid: "ws", id: "27", comment: "内蒙古卫视"},
		"shan1xws": {tid: "ws", id: "28", comment: "山西卫视"},
		"shan3xws": {tid: "ws", id: "29", comment: "陕西卫视"},
		"btws":     {tid: "ws", id: "30", comment: "兵团卫视"},
		"xjws":     {tid: "ws", id: "31", comment: "新疆卫视"},
		"xzws":     {tid: "ws", id: "32", comment: "西藏卫视"},
		"nxws":     {tid: "ws", id: "33", comment: "宁夏卫视"},
		"ybws":     {tid: "ws", id: "34", comment: "延边卫视"},
		"kbws":     {tid: "ws", id: "35", comment: "康巴卫视"},
		"dwqws":    {tid: "ws", id: "36", comment: "大湾区卫视"},
		"gdzj":     {tid: "ws", id: "37", comment: "广东珠江"},
		"xmws":     {tid: "ws", id: "38", comment: "厦门卫视"},
		"adws":     {tid: "ws", id: "39", comment: "安多卫视"},
		"nlws":     {tid: "ws", id: "40", comment: "农林卫视"},
		"ssws":     {tid: "ws", id: "41", comment: "三沙卫视"},

		"fhzw":    {tid: "gt", id: "1", comment: "凤凰中文"},
		"fhzx":    {tid: "gt", id: "2", comment: "凤凰资讯"},
		"fhxg":    {tid: "gt", id: "3", comment: "凤凰香港"},
		"tvbfc":   {tid: "gt", id: "5", comment: "TVB翡翠"},
		"tvbwxxw": {tid: "gt", id: "6", comment: "TVB无线新闻"},
		"tvbplus": {tid: "gt", id: "8", comment: "TVB Plus"},
		"tvbmz":   {tid: "gt", id: "9", comment: "TVB明珠"},
		"xgyxxw":  {tid: "gt", id: "11", comment: "香港有线新闻"},
		"xgyx18":  {tid: "gt", id: "12", comment: "香港有线18"},
		"nowbg":   {tid: "gt", id: "15", comment: "NOW爆谷"},
		"nowzb":   {tid: "gt", id: "16", comment: "NOW直播"},
		"nowxw":   {tid: "gt", id: "17", comment: "NOW新闻"},
		"nowcj":   {tid: "gt", id: "18", comment: "NOW财经"},
		"ztxw":    {tid: "gt", id: "43", comment: "中天新闻"},
		"dsxw":    {tid: "gt", id: "44", comment: "东森新闻"},
		"dscjxw":  {tid: "gt", id: "45", comment: "东森财经新闻"},
		"zsxw":    {tid: "gt", id: "46", comment: "中视新闻"},
		"tsxw":    {tid: "gt", id: "47", comment: "台视新闻"},
		"slxw":    {tid: "gt", id: "48", comment: "三立新闻"},
		"msxw":    {tid: "gt", id: "50", comment: "民视新闻"},
		"hsxw":    {tid: "gt", id: "51", comment: "华视新闻"},
		"jdsxw":   {tid: "gt", id: "52", comment: "镜电视新闻"},
		"hyxw":    {tid: "gt", id: "53", comment: "寰宇新闻"},
		"tvbsxw":  {tid: "gt", id: "54", comment: "TVBS新闻"},
	}

	resolve.RegisterHandler(ih)
}

// Handle 处理直播, 返回一个用于重定向的远程地址
func (ih *iptv345Handler) Handle(params resolve.HandleParams) (resolve.HandleResult, error) {
	reqParams, ok := ih.channels[params.ChName]
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
	sb.WriteString("\n2. 如果想通过直连播放（重定向）, 请使用可自定请求头的播放器（酷9、天光云影等）设置 Referer 头: [" + ih.serverHost + "/] 后才可正常播放")
	sb.WriteString("\n3. 如果是切片代理模式, 则无需指定请求头")
	sb.WriteString("\n4. 连续播放一段时间可能会发生断流, 需要重新加载")

	chs := []string{}
	for k, v := range ih.channels {
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

// SupportCustomHeaders 是否支持自定义请求头
// 如果返回 true, 会自动在帮助文档中加入标记
func (ih *iptv345Handler) SupportCustomHeaders() bool {
	return false
}

// Enabled 标记处理器是否是启用状态
func (ih *iptv345Handler) Enabled() bool {
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
			hkenExtraKey := secrets["hkenExtraKey"].(string)

			rawPlayUris := make([]string, len(playUris))
			for i, playUri := range playUris {
				res, err := ih.recoverPlayUri(hken, hkens, token, hkenExtraKey, playUri)
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

// readPageContent 读取网页源代码
func (ih *iptv345Handler) readPageContent(url string, header http.Header) (string, error) {
	res, err := https.Get(url).Header(header).Do()
	if err != nil {
		return "", fmt.Errorf("请求页面失败: %v", err)
	}
	defer res.Body.Close()
	pageBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("读取页面数据失败: %v", err)
	}
	return string(pageBytes), nil
}

// fetchOriginData 获取网页中的 string 混淆代码变量值
func (ih *iptv345Handler) fetchOriginData(tid, id string) (str2Decode string, playUris []string, fe error) {
	header := make(http.Header)
	header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36")

	// 1 解析出带有 token 的频道页面地址
	listPageUrl := fmt.Sprintf("%s/?tid=%s", ih.serverHost, tid)
	pageContent, err := ih.readPageContent(listPageUrl, header)
	if err != nil {
		fe = fmt.Errorf("读取频道列表页面失败: %v", err)
		return
	}

	linkGetter := regexp.MustCompile(fmt.Sprintf(`<a href="([^"]*&tid=%s&id=%s)"`, tid, id))
	if !linkGetter.MatchString(pageContent) {
		fe = fmt.Errorf("获取不到频道链接")
		return
	}
	chLink := linkGetter.FindStringSubmatch(pageContent)[1]

	// 2 读取频道页面原代码
	pageUrl := ih.serverHost + "/" + chLink
	pageCode, err := ih.readPageContent(pageUrl, header)
	if err != nil {
		fe = fmt.Errorf("读取频道页面失败: %v", err)
		return
	}

	// 3 解析出 string 变量
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

	// 4 解析原始的播放地址
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
func (ih *iptv345Handler) recoverPlayUri(hken, hkens, token, hkenExtraKey, originPlayUri string) (string, error) {
	resolvePlayUriCode := ih.decodePlayUriScript
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${hken}", hken, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${hkens}", hkens, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${token}", token, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${hkenExtraKey}", hkenExtraKey, 1)
	resolvePlayUriCode = strings.Replace(resolvePlayUriCode, "${uri}", originPlayUri, 1)

	vm := goja.New()
	res, err := vm.RunString(resolvePlayUriCode)
	if err != nil {
		return "", fmt.Errorf("执行 js 异常: %v", err)
	}
	return strings.Replace(res.String(), "type=.flv", "type=.m3u8", 1), nil
}

// sendSessionHeartBeatPkg 发送 session 心跳包, 防止 session 被关闭
func (ih *iptv345Handler) sendSessionHeartBeatPkg() {
	urls := ih.sessionCli.GetAllCacheUrls()
	for _, url := range urls {
		resp, err := https.Get(url).Header(ih.requestHeader.Clone()).Do()
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
