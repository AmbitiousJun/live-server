package web

import (
	"log"
	"net/http"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
	"github.com/gin-gonic/gin"
)

// HandleLive 调用处理器处理直播请求
func HandleLive(c *gin.Context) {
	hName := c.Param("handler")
	cName := c.Param("channel")
	if strs.AnyEmpty(hName, cName) {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}

	handler, ok := resolve.GetHandler(hName)
	if !ok {
		c.String(http.StatusBadRequest, "不支持的处理器")
		return
	}

	ua := c.Request.Header.Get("User-Agent")
	clientIp := c.ClientIP()
	log.Printf(colors.ToBlue("Client-IP: %s, User-Agent: %s"), clientIp, ua)

	result, err := handler.Handle(resolve.HandleParams{
		ChName: cName,
		UrlEnv: c.Query("url_env"),
	})
	if err != nil {
		c.String(http.StatusBadRequest, "处理失败: %v", err)
		return
	}

	if result.Type == resolve.ResultRedirect {
		log.Printf(colors.ToGreen("重定向到: %s"), result.Url)
		c.Redirect(http.StatusTemporaryRedirect, result.Url)
		return
	}

	if result.Type == resolve.ResultProxy {
		log.Println(colors.ToGreen("请求被代理"))
		c.Status(result.Code)
		if result.Header != nil {
			for key, values := range result.Header {
				for _, value := range values {
					c.Writer.Header().Add(key, value)
				}
			}
		}
		if result.Body != nil {
			c.Writer.Write(result.Body)
			c.Writer.Flush()
		}
		return
	}
}
