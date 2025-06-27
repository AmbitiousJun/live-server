package ipresolver

type R interface {

	// Resolve 解析 ip
	Resolve(string) (string, error)
}

// V4 获取 v4 解析器
func V4() R {
	return &v4Ipchaxun{}
}

// V6 获取 v6 解析器
func V6() R {
	return &v6Ipwcn{}
}
