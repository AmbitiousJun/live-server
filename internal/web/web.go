package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/secret"
	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/gin-gonic/gin"
)

// Listen 在指定端口上开启服务
func Listen(port int) error {
	r := gin.Default()
	r.GET("/handler/:handler/ch/:channel", HandleLive)
	r.HEAD("/handler/:handler/ch/:channel", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	r.GET("/black_ip", secret.Need(HandleAddBlackIp))
	r.GET("/env", secret.Need(env.StoreEnv))
	r.GET("/help", HandleHelpDoc)

	// 地域白名单操作接口
	war := r.Group("/white_area")
	war.GET("/set", secret.Need(whitearea.SetHandler))
	war.GET("/del", secret.Need(whitearea.DelHandler))

	// 利用服务器流量代理切片
	r.GET("/proxy_ts", resolve.ProxyTs)

	// 凤凰秀授权页
	r.GET("/feng/auth", secret.Need(ToFengAuthPage))
	// 配置页
	r.GET("/config", secret.Need(ToConfigPage))

	log.Printf(colors.ToYellow("在端口【%d】上开启 http 服务..."), port)
	log.Printf(colors.ToYellow("查看帮助文档请到浏览器访问: :%d/help"), port)
	return r.Run(fmt.Sprintf(":%d", port))
}
