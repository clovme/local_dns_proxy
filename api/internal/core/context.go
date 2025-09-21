package core

import (
	"github.com/gin-gonic/gin"
)

// Context 自定义gin.Context
type Context struct {
	*gin.Context
	Router routesMap
}

// Get 获取自定义gin.Context中的值
func (r *Context) Get(key string) any {
	if value, exists := r.Context.Get(key); exists {
		return value
	}
	return nil
}

// NewContext 创建自定义gin.Context
//
// 参数:
//   - ctx: gin.Context对象
//
// 返回值:
//   - *Context: 自定义gin.Context对象
//
// 说明:
//   - 创建自定义gin.Context对象，用于自定义路由和中间件。
func NewContext(ctx *gin.Context) *Context {
	return &Context{
		Context: ctx,
		Router:  initRoutesMap(),
	}
}
