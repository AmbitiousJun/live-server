package subm3u

import (
	"bufio"
	"errors"
	"regexp"
	"strings"
)

var (
	TvgIdReg      = regexp.MustCompile(`tvg-id="([^"]*)"`)
	TvgNameReg    = regexp.MustCompile(`tvg-name="([^"]*)"`)
	TvgLogoReg    = regexp.MustCompile(`tvg-logo="([^"]*)"`)
	GroupTitleReg = regexp.MustCompile(`group-title="([^"]*)"`)
	CustomNameReg = regexp.MustCompile(`,\s*(.*)$`)
)

// ReadContent 将 m3u8 原始文件整理成 Info 信息
func ReadContent(content string) (map[string]Info, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	// 检测首行
	if scanner.Scan() {
		firstLine := scanner.Text()
		if !strings.HasPrefix(firstLine, "#EXTM3U") {
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

		mapKey := info.TvgName
		if mapKey == "" {
			mapKey = info.CustomName
		}
		if mapKey == "" {
			mapKey = info.TvgId
		}
		if mapKey == "" {
			continue
		}

		if !scanner.Scan() {
			continue
		}
		info.Url = scanner.Text()
		res[mapKey] = info
	}

	return res, nil
}

// readChannelInfo 读取电视台信息
func readChannelInfo(line string) Info {
	res := Info{}

	if TvgIdReg.MatchString(line) {
		res.TvgId = TvgIdReg.FindStringSubmatch(line)[1]
	}

	if TvgNameReg.MatchString(line) {
		res.TvgName = TvgNameReg.FindStringSubmatch(line)[1]
	}

	if TvgLogoReg.MatchString(line) {
		res.TvgLogo = TvgLogoReg.FindStringSubmatch(line)[1]
	}

	if GroupTitleReg.MatchString(line) {
		res.GroupTitle = GroupTitleReg.FindStringSubmatch(line)[1]
	}

	if CustomNameReg.MatchString(line) {
		res.CustomName = CustomNameReg.FindStringSubmatch(line)[1]
	}

	return res
}
