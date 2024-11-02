package web

import (
	"log"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
	"github.com/gin-gonic/gin"
)

// HandleAddBlackIp 处理黑名单添加事件
func HandleAddBlackIp(c *gin.Context) {
	ip := c.Query("ip")
	if ip = strings.TrimSpace(ip); ip == "" {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}

	if err := AddBlackIp(ip); err != nil {
		log.Printf("添加黑名单失败: %v", err)
		c.String(http.StatusInternalServerError, "添加黑名单失败")
		return
	}

	c.String(http.StatusOK, "添加成功")
}

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

	if IsBlackIp(clientIp) {
		c.String(http.StatusForbidden, "私人服务器, 不对外公开, 望谅解！")
		return
	}

	result, err := handler.Handle(resolve.HandleParams{
		ChName:   cName,
		UrlEnv:   c.Query("url_env"),
		ProxyM3U: c.Query("proxy_m3u") == "1",
		ProxyTs:  c.Query("proxy_ts") == "1",
		Format:   c.Query("format"),
	})
	if err != nil {
		log.Printf(colors.ToRed("解析失败, handler: %s, errMsg: %v"), handler.Name(), err)
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

// HandleHelpDoc 输出帮助文档
func HandleHelpDoc(c *gin.Context) {
	c.String(http.StatusOK, resolve.HelpDoc())
}
