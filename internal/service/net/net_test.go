package net_test

import (
	"log"
	"testing"

	"github.com/AmbitiousJun/live-server/internal/service/net"
)

func TestIsPrivateIp(t *testing.T) {
	log.Println(net.IsPrivateIp("192.168.0.109"))
	log.Println(net.IsPrivateIp("127.0.0.1"))
	log.Println(net.IsPrivateIp("1.1.1.1"))
}
