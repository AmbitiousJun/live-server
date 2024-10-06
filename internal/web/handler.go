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

	result, err := handler.Handle(resolve.HandleParams{
		ChName: cName,
	})
	if err != nil {
		c.String(http.StatusBadRequest, "处理失败: %v", err)
		return
	}
	log.Printf(colors.ToGreen("重定向到: %s"), result)
	c.Redirect(http.StatusTemporaryRedirect, result)
}
