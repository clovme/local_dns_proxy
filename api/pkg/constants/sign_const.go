// 常量定义，常用标志定义, 用于上下文传递，全局使用, 不建议修改
// 用于项目中常使用的一些通用键(签名)，用于设置值和取值

package constants

const (
	ProjectName = "local_dns" // 项目名称
	WebTitle    = "本地DNS代理"

	HttpLogKey = "HTTP_LOG_KEY"
	LimitPage  = "LIMIT_PAGE"

	DNSStop    = "stop"
	DNSRunning = "running"
)
