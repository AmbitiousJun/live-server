package web

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/AmbitiousJun/live-server/internal/constant"
	"github.com/AmbitiousJun/live-server/internal/service/net"
	"github.com/AmbitiousJun/live-server/internal/service/resolve"
	"github.com/AmbitiousJun/live-server/internal/service/whitearea"
	"github.com/AmbitiousJun/live-server/internal/util/colors"
	"github.com/AmbitiousJun/live-server/internal/util/strs"
	"github.com/AmbitiousJun/live-server/internal/util/urls"
	"github.com/gin-gonic/gin"
)

// HandleAddBlackIp 处理黑名单添加事件
func HandleAddBlackIp(c *gin.Context) {
	ip := c.Query("ip")
	if ip = strings.TrimSpace(ip); ip == "" {
		c.String(http.StatusBadRequest, "参数不足")
		return
	}

	if err := net.AddBlackIp(ip); err != nil {
		log.Printf("添加黑名单失败: %v", err)
		c.String(http.StatusInternalServerError, "添加黑名单失败")
		return
	}

	c.String(http.StatusOK, "添加成功")
}

// HandleLive 调用处理器处理直播请求
func HandleLive(c *gin.Context) {
	c.Request.URL.RawQuery = trimDollarSuffix(urls.DecodeURI(c.Request.URL.RawQuery))
	hName := c.Param("handler")
	cName := trimDollarSuffix(urls.DecodeURI(c.Param("channel")))
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
	defer func() {
		log.Printf(colors.ToBlue("Client-IP: %s, User-Agent: %s"), clientIp, ua)
	}()

	if net.IsBlackIp(clientIp) {
		c.String(http.StatusNotFound, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
		return
	}

	if ipInfo, ok := net.GetIpAddrInfo(clientIp); ok {
		clientIp += " (" + ipInfo + ")"
		if !whitearea.Passable(ipInfo) {
			c.String(http.StatusNotFound, "私人服务器, 不对外公开, 望谅解！可前往官方仓库自行部署: https://github.com/AmbitiousJun/live-server")
			return
		}
	}

	result, err := handler.Handle(resolve.HandleParams{
		ChName:   cName,
		UrlEnv:   c.Query("url_env"),
		ProxyM3U: c.Query("proxy_m3u") == "1",
		ProxyTs:  c.Query("proxy_ts") == "1",
		Format:   c.Query("format"),
		ClientIp: c.ClientIP(),
	})
	if err != nil {
		log.Printf(colors.ToRed("解析失败, handler: %s, errMsg: %v"), handler.Name(), err)
		c.String(http.StatusBadRequest, "处理失败: %v", err)
		return
	}

	if result.Type == resolve.ResultRedirect {
		log.Printf(colors.ToGreen("重定向请求: %s"), result.Url)
		c.Redirect(http.StatusFound, result.Url)
		return
	}

	if result.Type == resolve.ResultProxy {
		c.Status(result.Code)
		if result.Header != nil {
			for key, values := range result.Header {
				for _, value := range values {
					c.Writer.Header().Add(key, value)
				}
			}
		}
		c.Writer.Header().Set("Cache-Control", "no-cache,no-store,must-revalidate")
		c.Writer.Header().Set("Expires", "0")
		c.Writer.Header().Set("Pragma", "no-cache")
		if result.Body != nil {
			c.Writer.Write(result.Body)
			c.Writer.Flush()
		}
		return
	}
}

// HandleHelpDoc 输出帮助文档
func HandleHelpDoc(c *gin.Context) {
	content := strings.ReplaceAll(resolve.HelpDoc(), "\n", "<br/>")
	tplBytes, _ := base64.StdEncoding.DecodeString(constant.HelpDocHtmlTemplate)
	result := strings.ReplaceAll(string(tplBytes), "${docContent}", content)
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, result)
}

// ToFengAuthPage 跳转到凤凰秀授权地址
func ToFengAuthPage(c *gin.Context) {
	bytes, _ := base64.StdEncoding.DecodeString(constant.FengAuthHtml)
	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	c.Writer.Write(bytes)
	c.Writer.Flush()
}

// ToConfigPage 跳转到配置页
func ToConfigPage(c *gin.Context) {
	bytes, _ := base64.StdEncoding.DecodeString(constant.ConfigPageHtml)
	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)
	c.Writer.Write(bytes)
	c.Writer.Flush()
}
