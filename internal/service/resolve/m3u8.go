package resolve

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/AmbitiousJun/live-server/internal/service/subm3u"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
)

// CommonM3U8 封装 m3u8 通用解析逻辑
type CommonM3U8 struct{}

// ResolveSub 解析 m3u8 订阅地址, 返回频道信息
func (cm *CommonM3U8) ResolveSub(client *https.CacheClient, subAddr string, headers http.Header) (map[string][]subm3u.Info, error) {
	// 参数校验
	if client == nil || strs.AnyEmpty(subAddr) {
		return nil, errors.New("参数不足")
	}

	// 请求 m3u 文本
	_, resp, err := client.Request(http.MethodGet, subAddr, headers, nil, true)
	if err != nil {
		return nil, fmt.Errorf("请求远程地址失败: %v", err)
	}
	defer resp.Body.Close()

	// 解析文本内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		client.RemoveUrlCache(subAddr)
		return nil, fmt.Errorf("读取远程响应失败: %v", err)
	}
	infos, err := subm3u.ReadContent(string(bodyBytes))
	if err != nil {
		client.RemoveUrlCache(subAddr)
		return nil, fmt.Errorf("解析远程响应失败: %v, 原始响应: %s", err, string(bodyBytes))
	}
	return infos, nil
}

// MatchChannel 匹配频道指定频道信息
func (cm *CommonM3U8) MatchChannel(infos map[string][]subm3u.Info, chName, format string) (subm3u.Info, error) {
	// 参数校验
	if infos == nil {
		return subm3u.Info{}, errors.New("参数不足")
	}

	destInfos, ok := infos[chName]
	if !ok || len(destInfos) == 0 {
		return subm3u.Info{}, fmt.Errorf("匹配频道名称失败: %s, 请检查远程地址是否有效", chName)
	}

	// 根据 format 参数筛选出目标频道
	resInfo := destInfos[0]
	formatIdx, err := strconv.Atoi(format)
	if strs.AllNotEmpty(format) && err != nil {
		log.Printf(colors.ToYellow("format 索引解析失败: %v, 有效范围: [1, %d], 本次请求使用默认值"), err, len(destInfos))
	}
	if err == nil {
		formatIdx--
		if formatIdx < 0 || formatIdx >= len(destInfos) {
			return subm3u.Info{}, fmt.Errorf("format 索引传递错误, 有效范围: [1, %d]", len(destInfos))
		}
		resInfo = destInfos[formatIdx]
	}
	return resInfo, nil
}

// ChannelSlice 获取频道信息列表中所有可用的频道
func (cm *CommonM3U8) ChannelSlice(infos map[string][]subm3u.Info) []string {
	channels := make([]string, 0, len(infos))
	for chName := range infos {
		channels = append(channels, chName)
	}
	return channels
}

// M3U8Result 根据处理器参数返回 m3u 地址的处理结果
func M3U8Result(url string, params HandleParams) (HandleResult, error) {
	// 如果无需代理, 直接重定向
	if !params.ProxyM3U {
		return HandleResult{Type: ResultRedirect, Url: url}, nil
	}

	content, err := ProxyM3U(url, params.Headers, params.ProxyTs, params.TsProxyMode, params.ClientHost)
	if err != nil {
		return HandleResult{}, fmt.Errorf("代理 m3u 失败: %v", err)
	}

	respHeader := make(http.Header)
	respHeader.Set("Content-Type", "application/vnd.apple.mpegurl; charset=utf-8")
	return HandleResult{
		Type:   ResultProxy,
		Code:   http.StatusOK,
		Body:   []byte(content),
		Header: respHeader,
	}, nil
}
