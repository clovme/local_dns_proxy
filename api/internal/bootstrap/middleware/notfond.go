package middleware

import (
	"local_dns_proxy/internal/core"
	"local_dns_proxy/pkg/enums/code"
	httpLog "local_dns_proxy/pkg/logger/http"
)

// RegisterNoRoute 注册404处理
func RegisterNoRoute(engine *core.Engine) {
	engine.NoRoute(func(c *core.Context) {
		httpLog.Error(c.Context).Msg("请求地址错误")
		c.JsonSafeDesc(code.RequestNotFound, nil)
		c.AbortWithStatus(404)
	})
}
