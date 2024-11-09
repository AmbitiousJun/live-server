package main

import (
	"flag"
	"log"

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

	if err := web.Listen(*port); err != nil {
		log.Fatal(err)
	}
}
