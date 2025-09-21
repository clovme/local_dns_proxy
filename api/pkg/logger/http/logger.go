package http

import (
	"bytes"
	"local_dns_proxy/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io/ioutil"
)

// 核心字段统一封装
// 注意：
//   - 日志字段顺序：Method TraceID Status ClientIP UserAgent Path RequestURI Referer ContentLength
//   - 日志字段内容：Method=请求方法，TraceID=请求跟踪ID，Status=响应状态码，ClientIP=客户端IP，UserAgent=用户代理，Path=请求路径，RequestURI=请求URI，Referer=请求来源，ContentLength=响应内容长度
//   - 日志字段类型：Method=string，TraceID=string，Status=int，ClientIP=string，UserAgent=string，Path=string，RequestURI=string，Referer=string，ContentLength=int
//
// 参数：
//   - _log: zerolog.Event
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func _field(_log *zerolog.Event, c *gin.Context) *zerolog.Event {
	return _log.Str("Method", c.Request.Method).
		Str("TraceID", c.GetHeader("X-Trace-Id")).
		Int("Status", c.Writer.Status()).
		Str("X-Requested-With", c.GetHeader("X-Requested-With")).
		Str("ClientIP", c.ClientIP()).
		Str("UserAgent", c.Request.UserAgent()).
		Str("Path", c.Request.URL.Path).
		Str("RequestURI", c.Request.RequestURI).
		Str("Referer", c.Request.Referer()).
		Str("RemoteAddr", c.Request.RemoteAddr).
		Int("Length", c.Writer.Size())
}

// debug 专用扩展字段
// 注意：
//   - 日志字段顺序：Body Header Query Form MultipartForm RemoteAddr TLS Response
//   - 日志字段内容：Body=请求体，Header=请求头，Query=请求参数，Form=表单参数，MultipartForm=多部分表单参数，RemoteAddr=远程地址，TLS=TLS信息，Response=响应信息
//   - 日志字段类型：Body=[]byte，Header=map[string][]string，Query=map[string][]string，Form=map[string][]string，MultipartForm=*multipart.Form，RemoteAddr=string，TLS=tls.ConnectionState，Response=gin.ResponseWriter
//
// 参数：
//   - _log: zerolog.Event
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func _addDebugFields(_log *zerolog.Event, c *gin.Context) *zerolog.Event {
	// Safe 读取 body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		// 读取后要重置 body，防止后续读不到
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return _log.
		Bytes("Body", bodyBytes).
		Interface("Header", c.Request.Header).
		Interface("Query", c.Request.URL.Query()).
		Interface("PostForm", c.Request.PostForm).
		Interface("Form", c.Request.Form).
		Interface("Body", c.Request.Body).
		Interface("TransferEncoding", c.Request.TransferEncoding).
		Interface("MultipartForm", c.Request.MultipartForm).
		Interface("TLS", c.Request.TLS).
		Interface("Response", c.Writer)
}

// Debug 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Debug(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpDebug)
	return _addDebugFields(_field(_log.Debug(), c), c)
}

// Info 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Info(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpInfo)
	return _field(_log.Info(), c)
}

// Warn 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Warn(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpWarn)
	return _field(_log.Warn(), c)
}

// Error 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Error(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpError)
	return _field(_log.Error(), c).Interface("Errors", c.Errors)
}

// Trace 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Trace(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpTrace)
	return _field(_log.Trace(), c)
}

// Panic 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Panic(c *gin.Context) *zerolog.Event {
	_log := logger.GetLogger(logger.HttpTrace)
	return _field(_log.Panic(), c)
}

// Log 仅用于正常请求日志，异常和慢请求请显式调用 Error / Warn
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Log(c *gin.Context) *zerolog.Event {
	switch logger.CurrentCfg.Lvl {
	case zerolog.DebugLevel:
		return Debug(c)
	default:
		return Info(c)
	}
}
