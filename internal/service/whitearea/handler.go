package whitearea

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetHandler 设置白名单处理器
func SetHandler(c *gin.Context) {
	Set(c.Query("area"))
	c.String(http.StatusOK, "set 操作成功")
}

// DelHandler 移除白名单处理器
func DelHandler(c *gin.Context) {
	Del(c.Query("area"))
	c.String(http.StatusOK, "del 操作成功")
}
