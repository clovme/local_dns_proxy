package routers

import (
	"gorm.io/gorm"
	"local_dns_proxy/internal/bootstrap/middleware"
	"local_dns_proxy/internal/core"
	"time"
)

// regeditMiddleware 注册中间件
func regeditMiddleware(engine *core.Engine) {
	engine.Engine.Use(middleware.LogMiddleware(2 * time.Second)) // 请求日志，记录全流程
	engine.Engine.Use(middleware.FaviconMiddleware())            // favicon.ico
	middleware.ResourceDirInterception(engine)                   // 资源目录拦截并加载静态资源文件，放在日志中间件后面，其他中间件前面，保证能打印日志，又不走其他中间件
	engine.Use(middleware.RecoveryMiddleware())                  // 抓捕 panic，防止服务崩溃
	engine.Engine.Use(middleware.CorsMiddleware())               // 跨域，处理请求头
}

// regeditTemplate 注册模板
func regeditTemplate(engine *core.Engine) {
	engine.SetHTMLTemplate()
}

// regeditRoutes 注册路由
func regeditRoutes(engine *core.Engine, db *gorm.DB) {
	routers := routeGroup{
		// 接口层
		dnsApi: engine.Group("/api/v1"),

		// 视图层
		dnsView: engine.Group("/"),
	}

	// 注册路由
	routers.register(db)

	// 注册404处理
	middleware.RegisterNoRoute(engine)
}

// Initialization 初始化 web 服务
// 参数：
//   - db: 数据库连接对象
//   - staticDir: 静态文件目录
//
// 返回值：
//   - *gin.Engine: 初始化后的 Gin 引擎
func Initialization(db *gorm.DB) *core.Engine {
	engine := core.New()

	regeditTemplate(engine)
	regeditMiddleware(engine)
	regeditRoutes(engine, db)

	return engine
}
