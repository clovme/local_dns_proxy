package middleware

import (
	"github.com/gin-gonic/gin"
	"io/fs"
	"local_dns_proxy/internal/core"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/pkg/enums/code"
	"local_dns_proxy/public"
	"net/http"
	"strings"
)

// FaviconMiddleware 加载 favicon.ico
func FaviconMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/favicon.ico" {
			c.Data(200, "image/x-icon", public.Favicon)
			c.Set(constants.HttpLogKey, "favicon")
			c.Abort()
			return
		}
		c.Next()
	}
}

// ResourceDirInterception 静态路由拦截
func ResourceDirInterception(engine *core.Engine) {
	// 静态资源路由拦截中间件处理，中间件部分
	engine.Use(func(c *core.Context) {
		if strings.HasSuffix(c.FullPath(), "*filepath") {
			c.Set(constants.HttpLogKey, "静态资源")
		}
		if strings.EqualFold(c.Request.URL.Path, "/assets/") {
			c.JsonSafeDesc(code.RequestNotFound, nil)
			c.AbortWithStatus(404)
		}
	})

	// 读取嵌入二进制的静态资源目录
	staticFS, _ := fs.Sub(public.WebFS, "web/assets")
	engine.Engine.StaticFS("/assets", http.FS(staticFS))
}
