package log

import (
	"local_dns_proxy/pkg/logger"
	"github.com/rs/zerolog"
)

// Debug 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Debug() *zerolog.Event {
	_log := logger.GetLogger(logger.AppDebug)
	return _log.Debug()
}

// Error 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Error() *zerolog.Event {
	_log := logger.GetLogger(logger.AppError)
	return _log.Error()
}

// Fatal 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Fatal() *zerolog.Event {
	_log := logger.GetLogger(logger.AppFatal)
	return _log.Fatal()
}

// Info 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Info() *zerolog.Event {
	_log := logger.GetLogger(logger.AppInfo)
	return _log.Info()
}

// Panic 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Panic() *zerolog.Event {
	_log := logger.GetLogger(logger.AppPanic)
	return _log.Panic()
}

// Trace 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Trace() *zerolog.Event {
	_log := logger.GetLogger(logger.AppTrace)
	return _log.Trace()
}

// Warn 不同 level 封装
// 参数：
//   - c: gin.Context
//
// 返回值：
//   - *zerolog.Event
func Warn() *zerolog.Event {
	_log := logger.GetLogger(logger.AppWarn)
	return _log.Warn()
}
