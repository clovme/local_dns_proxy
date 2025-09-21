package core

import (
	"fmt"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"html/template"
	"local_dns_proxy/pkg/logger/log"
	"local_dns_proxy/pkg/utils/file"
	"local_dns_proxy/public"
	"os"
	"sort"
)

// Engine 自定义gin.Engine
type Engine struct {
	Engine *gin.Engine
	RouterGroup
}

// SetHTMLTemplate 设置模版
//
// 参数:
//   - funcMap: 模版函数
//
// 返回值:
//   - *template.Template: 模版对象
func (engine *Engine) SetHTMLTemplate() {
	tmpl := template.Must(template.New("").ParseFS(public.WebFS, "web/*.html"))
	engine.Engine.SetHTMLTemplate(tmpl)
}

// Use 注册中间件
//
// 参数:
//   - middleware: 中间件函数列表
//
// 说明:
//   - 注册中间件函数，用于在请求处理过程中进行预处理或后处理。
func (engine *Engine) Use(middleware ...HandlerFunc) {
	middlewareList := make([]gin.HandlerFunc, 0)
	for _, m := range middleware {
		middlewareList = append(middlewareList, wrapHandler(m))
	}
	engine.Engine.Use(middlewareList...)
}

// NoRoute 注册404路由
//
// 参数:
//   - handler: 404路由处理函数
//
// 说明:
//   - 注册404路由处理函数，当请求的路由不存在时调用。
func (engine *Engine) NoRoute(handler HandlerFunc) {
	engine.Engine.NoRoute(wrapHandler(handler))
}

// Group 注册路由组
//
// 参数:
//   - relativePath: 路由组的相对路径
//   - handlers: 路由组的中间件函数列表
//
// 返回值:
//   - *RouterGroup: 路由组对象
//
// 说明:
//   - 注册路由组，用于组织和管理相关的路由。
func (engine *Engine) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return groupFunc(engine.RouterGroup, relativePath, handlers...)
}

// Routes 获取所有路由信息
//
// 返回值:
//   - []RoutesInfo: 路由信息列表
//
// 说明:
//   - 获取所有注册的路由信息，包括路由路径、请求方法、路由名称、路由类型和路由描述。
func (engine *Engine) Routes() []RoutesInfo {
	engine.checkDuplicateRoutes()
	// 预设方法顺序
	methodSort := []string{"GET", "POST", "PUT", "DELETE"}
	methodSet := make(map[string]struct{}, len(methodSort))
	for _, m := range methodSort {
		methodSet[m] = struct{}{}
	}

	// 分类存储
	methodRoutes := make(map[string][]RoutesInfo)

	// 收集所有 key 并排序
	keys := make([]string, 0, len(routesInfo))
	for k := range routesInfo {
		keys = append(keys, k)
	}
	sort.Strings(keys) // 字典序

	// 遍历 routesInfo 分类到 methodRoutes，同时动态扩展 methodSort
	for _, k := range keys {
		info := routesInfo[k]
		methodRoutes[info.Method] = append(methodRoutes[info.Method], info)
		if _, ok := methodSet[info.Method]; !ok {
			methodSort = append(methodSort, info.Method)
			methodSet[info.Method] = struct{}{}
		}
	}

	// 拼接最终结果
	result := make([]RoutesInfo, 0, len(routesInfo))
	for _, m := range methodSort {
		result = append(result, methodRoutes[m]...)
	}

	return result
}

// Run 启动HTTPS服务
//
// 参数:
//   - host: 主机名
//   - port: 端口号
//
// 返回值:
//   - error: 错误信息
func (engine *Engine) Run(ip string, port int) {
	engine.checkDuplicateRoutes()

	// 打印路由信息
	for i, route := range engine.Routes() {
		method := fmt.Sprintf("[%s]", route.Method)
		hostAddr := fmt.Sprintf("%s:%d", ip, port)
		log.Info().Msgf("%03d %-8s http://%s%-40s%-10s%-20s%s", i+1, method, hostAddr, route.Path, "-->", route.Name, route.Description)
	}

	absPath, err := file.GetFileAbsPath(".")
	if err != nil {
		return
	}
	log.Info().Msgf("程序所在路径 %s", absPath)

	addr := fmt.Sprintf("%s:%d", ip, port) // ⚠️ 用IP，不要用Domain
	if err := engine.Engine.Run(addr); err != nil {
		log.Error().Err(err).Msg("服务启动失败")
		os.Exit(-1)
	}
}

func (engine *Engine) checkDuplicateRoutes() {
	routesName := make(map[string]RoutesInfo, len(routesInfo))

	for _, route := range routesInfo {
		if _, ok := routesName[route.Name]; !ok {
			routesName[route.Name] = route
		} else {
			panic(fmt.Sprintf("路由名称重复: %s\n   %+v\n   %+v\n", route.Name, route, routesName[route.Name]))
		}
	}
}

// New 创建自定义gin.Engine
//
// 参数:
//   - opts: gin.Engine的选项函数列表
//
// 返回值:
//   - *Engine: 自定义gin.Engine对象
//
// 说明:
//   - 创建自定义gin.Engine对象，用于自定义路由和中间件。
func New(opts ...gin.OptionFunc) *Engine {
	routesInfo = make(map[string]RoutesInfo)

	// 创建 gin web 实例
	engine := gin.New(opts...)
	// 注册全局 gzip
	engine.Use(gzip.Gzip(gzip.DefaultCompression))

	return &Engine{
		Engine: engine,
		RouterGroup: RouterGroup{
			RouterGroup: &engine.RouterGroup,
		},
	}
}
