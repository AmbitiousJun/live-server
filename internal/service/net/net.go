package net

import "net"

// privateBlocks 局域网网段
var privateBlocks = []string{
	"10.0.0.0/8",     // 私有 IPv4
	"172.16.0.0/12",  // 私有 IPv4
	"192.168.0.0/16", // 私有 IPv4
	"127.0.0.0/8",    // 本地回环 IPv4
	"::1/128",        // 本地回环 IPv6
	"fc00::/7",       // 私有 IPv6
	"fe80::/10",      // 链路本地 IPv6
}

// IsPrivateIp 判断一个 ip
func IsPrivateIp(ip string) bool {
	pip := net.ParseIP(ip)
	if pip == nil {
		return false
	}

	for _, block := range privateBlocks {
		_, subnet, _ := net.ParseCIDR(block)
		if subnet.Contains(pip) {
			return true
		}
	}

	return false
}
