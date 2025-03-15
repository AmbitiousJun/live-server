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

// FindEnv 查询环境变量
func FindEnv(c *gin.Context) {
	key := c.Query("key")
	if strs.AnyEmpty(key) {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}
	if value, ok := Get(key); ok {
		c.String(http.StatusOK, value)
		return
	}
	c.String(http.StatusNotFound, "环境变量不存在")
}

// DeleteEnv 删除环境变量
func DeleteEnv(c *gin.Context) {
	key := c.PostForm("key")
	if strs.AnyEmpty(key) {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}
	Remove(key)
	c.String(http.StatusOK, "删除成功")
}
