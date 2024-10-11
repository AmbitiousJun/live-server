package m3u8

// Info m3u 格式电视直播的电视台信息封装
type Info struct {
	TvgName    string // 电视台名称
	TvgId      string // 电视台 Id
	TvgLogo    string // 台标地址
	GroupTitle string // 分类
	CustomName string // 自定义的电台名称
	Url        string // 直播源地址
}
