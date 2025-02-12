package env

import (
	"net/http"

	"github.com/AmbitiousJun/live-server/internal/util/strs"
	"github.com/gin-gonic/gin"
)

// StoreEnv 设置环境变量
func StoreEnv(c *gin.Context) {
	key := c.PostForm("key")
	value := c.PostForm("value")
	if strs.AnyEmpty(key, value) {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}
	Set(key, value)
	c.String(http.StatusOK, "设置成功")
}
