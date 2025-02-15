package m3u8

import (
	"bufio"
	"errors"
	"log"
	"net/url"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/util/colors"
)

// ExtractUrl 提取 m3u 地址信息
func ExtractUrl(m3uLink string) (info UrlInfo) {
	u, err := url.Parse(m3uLink)
	if err != nil {
		log.Printf(colors.ToRed("解析 m3u 地址失败: %v, m3u: %s"), err, m3uLink)
		return
	}
	info.Host = u.Scheme + "://" + u.Host

	// 切掉 query 参数
	qStartIdx := strings.Index(m3uLink, "?")
	if qStartIdx != -1 {
		m3uLink = m3uLink[:qStartIdx]
	}

	// 找到最后一个 /
	lastSlashIdx := strings.LastIndex(m3uLink, "/")
	if lastSlashIdx == -1 {
		log.Printf(colors.ToRed("url 地址缺少 / 分隔符: %s"), m3uLink)
		return
	}
	info.BaseDir = strings.TrimPrefix(m3uLink[:lastSlashIdx], info.Host)
	return
}

// ReadContent 将 m3u8 原始文件整理成 Info 信息
//
// 有的 m3u 切片地址是相对路径, 需要手动拼接前缀
func ReadContent(urlInfo UrlInfo, content string) (Info, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	res := Info{}
	urlInfo.BaseDir = strings.TrimSuffix(urlInfo.BaseDir, "/")

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
		secondLine := scanner.Text()
		if !strings.HasPrefix(secondLine, "http") {
			if strings.HasPrefix(secondLine, "/") {
				// 绝对路径
				secondLine = urlInfo.Host + secondLine
			} else {
				// 相对路径
				secondLine = urlInfo.Host + urlInfo.BaseDir + "/" + secondLine
			}
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
