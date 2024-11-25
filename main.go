package main

import (
	"flag"
	"log"

	"github.com/AmbitiousJun/live-server/internal/service/resolve/handler"
	"github.com/AmbitiousJun/live-server/internal/service/secret"
	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
	"github.com/AmbitiousJun/live-server/internal/service/ytdlp"
	"github.com/AmbitiousJun/live-server/internal/web"
	"github.com/gin-gonic/gin"
)

func main() {
	port := flag.Int("p", 5666, "自定义程序运行的端口号")
	prod := flag.Bool("prod", false, "是否开启线上模式")
	flag.Parse()

	if *prod {
		gin.SetMode(gin.ReleaseMode)
	}

	handler.Init()
	ytdlp.Init()
	if err := secret.Init(); err != nil {
		log.Panicf("初始化程序密钥失败: %v", err)
	}
	if err := whitearea.Init(); err != nil {
		log.Panicf("初始化地域白名单失败: %v", err)
	}

	if err := web.Listen(*port); err != nil {
		log.Panic(err)
	}
}
