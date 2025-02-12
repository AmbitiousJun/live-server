package web

import (
	"fmt"
	"strconv"
	"time"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/https"
	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		ipAddr := c.GetString(constant.Gin_IpAddrInfoKey)
		if ipAddr != "" {
			ipAddr = "（" + ipAddr + "）"
		}

		// 记录日志
		fmt.Printf("%s %s | %s | %s | %s%s | %s | %s %s\n",
			colors.ToYellow("[ls:"+constant.Version+"]"),
			start.Format("2006-01-02 15:04:05"),

			colorStatusCode(c.Writer.Status()),

			time.Since(start),

			colors.ToBlue(c.ClientIP()),
			colors.ToBlue(ipAddr),

			colors.ToGray(c.GetHeader("User-Agent")),

			colors.ToBlue(c.Request.Method),
			c.Request.RequestURI,
		)
	}
}

// colorStatusCode 将响应码打上颜色标记
func colorStatusCode(code int) string {
	str := strconv.Itoa(code)
	if https.IsSuccessCode(code) || https.IsRedirectCode(code) {
		return colors.ToGreen(str)
	}
	if https.IsErrorCode(code) {
		return colors.ToRed(str)
	}
	return colors.ToBlue(str)
}
