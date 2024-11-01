package m3u8

import (
	"bufio"
	"errors"
	"strings"
)

// ExtractPrefix 提取 m3u 地址请求前缀
func ExtractPrefix(url string) string {
	// 切掉 query 参数
	qStartIdx := strings.Index(url, "?")
	if qStartIdx != -1 {
		url = url[:qStartIdx]
	}

	// 找到最后一个 /
	lastSlashIdx := strings.LastIndex(url, "/")
	if lastSlashIdx == -1 {
		// 不规范的 url 地址
		return url
	}
	return url[:lastSlashIdx]
}

// ReadContent 将 m3u8 原始文件整理成 Info 信息
//
// 有的 m3u 切片地址是相对路径, 需要手动拼接前缀
func ReadContent(prefix, content string) (Info, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	res := Info{}
	prefix = strings.TrimSuffix(prefix, "/")

	// 检测首行
	if scanner.Scan() {
		firstLine := scanner.Text()
		if !strings.HasPrefix(firstLine, "#EXTM3U") {
			return Info{}, errors.New("不是标准的 m3u8 文本")
		}
		res.HeadComments = append(res.HeadComments, firstLine)
	}

	// 遍历文本, 整理信息
	for scanner.Scan() {
		firstLine := scanner.Text()

		if strings.HasPrefix(firstLine, "#") {
			cmtName := extractCommentName(firstLine)
			// 头注释
			if _, ok := ParentHeadComments[cmtName]; ok {
				res.HeadComments = append(res.HeadComments, firstLine)
				continue
			}
			// 尾注释
			if _, ok := ParentTailComments[cmtName]; ok {
				res.TailComments = append(res.TailComments, firstLine)
				continue
			}
		}

		// 除了头尾注释, 只识别切片注释
		if !strings.HasPrefix(firstLine, "#EXTINF") {
			continue
		}
		// 切片注释下一行必须有内容
		if !scanner.Scan() {
			continue
		}

		// 为切片地址补充前缀
		secondLine := strings.TrimPrefix(scanner.Text(), prefix)
		if !strings.HasPrefix(secondLine, "http") {
			if !strings.HasPrefix(secondLine, "/") {
				secondLine = "/" + secondLine
			}
			secondLine = prefix + secondLine
		}

		tsInfo := TsInfo{
			ExtInf: firstLine,
			Url:    secondLine,
		}
		res.TsInfos = append(res.TsInfos, tsInfo)
	}

	return res, nil
}

// extractCommentName 提取一行文本中的注释名称
//
// 注释名称和注释内容用 : 符号分隔开
// 若找不到 : 符号, 表明一整行就是注释名称
func extractCommentName(line string) string {
	idx := strings.Index(line, ":")
	if idx == -1 {
		return line
	}
	return line[:idx]
}
