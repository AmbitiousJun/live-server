package web

import (
	"fmt"
	"log"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/resolve/handler"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/gin-gonic/gin"
)

// Listen 在指定端口上开启服务
func Listen(port int) error {
	r := gin.Default()
	r.GET("/handler/:handler/ch/:channel", HandleLive)
	r.HEAD("/handler/:handler/ch/:channel", HandleLive)

	r.GET("/black_ip", HandleAddBlackIp)
	r.GET("/env", env.StoreEnv)
	r.GET("/help", HandleHelpDoc)

	// 利用服务器流量代理切片
	r.GET("/proxy_ts", resolve.ProxyTs)

	handler.Init()

	log.Printf(colors.ToYellow("在端口【%d】上开启 http 服务..."), port)
	log.Printf(colors.ToYellow("查看帮助文档请到浏览器访问: :%d/help"), port)
	return r.Run(fmt.Sprintf(":%d", port))
}
