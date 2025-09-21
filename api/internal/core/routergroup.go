package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// RouterGroup 路由组
type RouterGroup struct {
	*gin.RouterGroup
	uriPrefix string
}

type RoutesInfo struct {
	Path        string // 路由路径(/regedit.html)唯一标识，不能重复
	Name        string // 路由名称(regeditApi)唯一标识，不能重复
	Method      string // 请求方法(GET)
	Group       string // 路由分组(noAuthView)
	Description string // 路由描述(注册页面)
}

// HandlerFunc 路由处理函数
type HandlerFunc func(*Context)

var routesInfo = map[string]RoutesInfo{}

// wrapHandler 路由处理函数包装
func wrapHandler(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlerFunc(NewContext(c))
	}
}

func groupFunc(group RouterGroup, relativePath string, handlers ...HandlerFunc) *RouterGroup {
	handlerList := make([]gin.HandlerFunc, 0)
	for _, h := range handlers {
		handlerList = append(handlerList, wrapHandler(h))
	}
	uriPrefix := group.uriPrefix + relativePath
	if strings.HasSuffix(relativePath, "/") {
		uriPrefix = strings.TrimSuffix(uriPrefix, "/")
	}
	return &RouterGroup{
		RouterGroup: group.RouterGroup.Group(relativePath, handlerList...),
		uriPrefix:   uriPrefix,
	}
}

// handle 路由处理函数注册
func (group *RouterGroup) handle(httpMethod, relativePath string, handler HandlerFunc, groupName, name, description string) {
	key := fmt.Sprintf("%s:%s%s", httpMethod, group.uriPrefix, relativePath)
	routesInfo[key] = RoutesInfo{
		Method:      httpMethod,
		Path:        group.uriPrefix + relativePath,
		Name:        name,
		Group:       groupName,
		Description: description,
	}
	group.RouterGroup.Handle(httpMethod, relativePath, wrapHandler(handler))
}

// POST 请求 router.Handle("POST", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) POST(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodPost, relativePath, handler, groupName, name, description)
}

// GET 请求 router.Handle("GET", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) GET(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodGet, relativePath, handler, groupName, name, description)
}

// DELETE 请求 router.Handle("DELETE", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) DELETE(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodDelete, relativePath, handler, groupName, name, description)
}

// PATCH 请求 router.Handle("PATCH", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) PATCH(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodPatch, relativePath, handler, groupName, name, description)
}

// PUT 请求 router.Handle("PUT", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) PUT(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodPut, relativePath, handler, groupName, name, description)
}

// OPTIONS 请求 router.Handle("OPTIONS", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) OPTIONS(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodOptions, relativePath, handler, groupName, name, description)
}

// HEAD 请求 router.Handle("HEAD", path, handler, groupName, name, description).
//
// 参数:
//   - relativePath: 路由路径(/regedit.html)唯一标识，不能重复
//   - handler: 路由处理函数
//   - groupName: 路由分组(noAuthView)
//   - typ: 路由类型(view)
//   - name: 路由名称(regeditApi)唯一标识，不能重复
//   - description: 路由描述(注册页面)
func (group *RouterGroup) HEAD(relativePath string, handler HandlerFunc, groupName, name, description string) {
	group.handle(http.MethodHead, relativePath, handler, groupName, name, description)
}

// Group creates a new router group. You should add all the routes that have common middlewares or the same path prefix.
// For example, all the routes that use a common middleware for authorization could be grouped.
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return groupFunc(*group, relativePath, handlers...)
}
