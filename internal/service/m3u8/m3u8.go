package m3u8

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
)

// ReadContent 将 m3u8 原始文件整理成 Info 信息
func ReadContent(content string) (map[string]Info, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	// 检测首行
	if scanner.Scan() {
		firstLine := scanner.Text()
		if firstLine != "#EXTM3U" {
			return nil, errors.New("不是标准的 m3u8 文本")
		}
	}

	// 遍历文本, 每两行合成一个信息
	res := make(map[string]Info)
	for scanner.Scan() {
		firstLine := scanner.Text()
		if !strings.HasPrefix(firstLine, "#EXTINF") {
			continue
		}
		info := readChannelInfo(firstLine)
		if info.TvgName == "" {
			continue
		}

		if !scanner.Scan() {
			continue
		}
		info.Url = scanner.Text()
		res[info.TvgName] = info
	}

	return res, nil
}

// readChannelInfo 读取电视台信息
func readChannelInfo(line string) Info {
	// 格式化 info 信息
	line = strings.ReplaceAll(line, ",", " ")
	line = strings.ReplaceAll(line, `"`, "")
	reg := regexp.MustCompile(`\s+`)
	line = reg.ReplaceAllString(line, " ")

	// 根据空格分割
	segs := strings.Split(line, " ")
	res := Info{}
	for i, seg := range segs {
		// 最后一个分段固定为频道自定义别名
		if i == len(segs)-1 {
			res.CustomName = seg
			continue
		}

		// 根据 = 分割
		kvs := strings.Split(seg, "=")
		if len(kvs) != 2 {
			continue
		}

		switch kvs[0] {
		case "tvg-name":
			res.TvgName = kvs[1]
		case "tvg-id":
			res.TvgId = kvs[1]
		case "tvg-logo":
			res.TvgLogo = kvs[1]
		case "group-title":
			res.GroupTitle = kvs[1]
		}
	}
	return res
}
