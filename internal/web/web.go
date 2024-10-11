package web

import (
	"fmt"
	"log"

	"github.com/AmbitiousJun/live-server/internal/service/env"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/gin-gonic/gin"
)

// Listen 在指定端口上开启服务
func Listen(port int) error {
	r := gin.Default()
	r.GET("/handler/:handler/ch/:channel", HandleLive)
	r.GET("/env", env.StoreEnv)
	log.Printf(colors.ToBlue("在端口【%d】上开启 http 服务..."), port)
	return r.Run(fmt.Sprintf(":%d", port))
}
