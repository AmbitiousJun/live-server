package m3u8

import "strings"

// ParentHeadComments 记录文件头注释
var ParentHeadComments = map[string]struct{}{
	"#EXTM3U": {}, "#EXT-X-VERSION": {}, "#EXT-X-MEDIA-SEQUENCE": {},
	"#EXT-X-TARGETDURATION": {}, "#EXT-X-MEDIA": {}, "#EXT-X-INDEPENDENT-SEGMENTS": {},
	"#EXT-X-STREAM-INF": {}, "#EXT-X-DISCONTINUITY-SEQUENCE": {},
}

// ParentTailComments 记录文件尾注释
var ParentTailComments = map[string]struct{}{
	"#EXT-X-ENDLIST": {},
}

// Info m3u8 信息
type Info struct {
	HeadComments []string // 头注释
	TailComments []string // 尾注释
	TsInfos      []TsInfo // ts 切片信息
}

// TsInfo ts 切片信息
type TsInfo struct {
	ExtInf string
	Url    string
}

// ContentFunc 返回 m3u 文本, 允许修改 ts
func (i *Info) ContentFunc(f func(tsIdx int, tsUrl string) string) string {
	sb := strings.Builder{}

	for _, cmt := range i.HeadComments {
		sb.WriteString(cmt)
		sb.WriteByte('\n')
	}

	for idx, tsInfo := range i.TsInfos {
		sb.WriteString(tsInfo.ExtInf)
		sb.WriteByte('\n')
		sb.WriteString(f(idx, tsInfo.Url))
		sb.WriteByte('\n')
	}

	for _, cmt := range i.TailComments {
		sb.WriteString(cmt)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// Content 返回 m3u 文本
func (i *Info) Content() string {
	return i.ContentFunc(func(tsIdx int, tsUrl string) string { return tsUrl })
}

// UrlInfo m3u url 地址信息
type UrlInfo struct {
	Host    string // 主机地址
	BaseDir string // url 对应资源所在的目录路径
}
