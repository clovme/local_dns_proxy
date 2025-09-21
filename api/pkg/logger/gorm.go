package logger

import (
	"context"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// GormLogger 自定义 gorm logger，实现 gorm logger.Interface 接口
type GormLogger struct {
	level logger.LogLevel
}

// GetGormLogger 获取 gorm logger
// 返回值：
//   - *gorm.Config: gorm config
//
// 注意：
//   - 日志级别映射：
//   - zerolog.TraceLevel: logger.Info
//   - zerolog.DebugLevel: logger.Info
//   - zerolog.InfoLevel: logger.Info
//   - zerolog.WarnLevel: logger.Warn
//   - zerolog.ErrorLevel: logger.Error
//   - zerolog.FatalLevel: logger.Error
//   - zerolog.PanicLevel: logger.Error
//   - zerolog.Disabled: logger.Silent
//   - zerolog.NoLevel: logger.Silent
func GetGormLogger() *gorm.Config {
	var logLevelMap = map[zerolog.Level]logger.LogLevel{
		zerolog.TraceLevel: logger.Info,
		zerolog.DebugLevel: logger.Info,
		zerolog.InfoLevel:  logger.Info,
		zerolog.WarnLevel:  logger.Warn,
		zerolog.ErrorLevel: logger.Error,
		zerolog.FatalLevel: logger.Error,
		zerolog.PanicLevel: logger.Error,
		zerolog.Disabled:   logger.Silent,
		zerolog.NoLevel:    logger.Silent,
	}

	level, ok := logLevelMap[CurrentCfg.Lvl]
	if !ok {
		level = logger.Info
	}

	return &gorm.Config{
		CreateBatchSize: 1000,
		Logger:          &GormLogger{level: level},
	}
}

// LogMode 日志级别
// 参数：
//   - level: logger.LogLevel
//
// 返回值：
//   - logger.Interface: gorm logger
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.level = level
	return &newlogger
}

// Info 打印 info 日志
// 参数：
//   - ctx: context.Context
//   - s: string
//   - args: ...interface{}
//
// 返回值：
//   - 无
func (l *GormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	if l.level >= logger.Info {
		_log := loggers[DbInfo]
		_log.Info().Msgf(s, args...)
	}
}

// Warn 打印 warn 日志
// 参数：
//   - ctx: context.Context
//   - s: string
//   - args:...interface{}
//
// 返回值：
//   - 无
func (l *GormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	if l.level >= logger.Warn {
		_log := loggers[DbWarn]
		_log.Warn().Msgf(s, args...)
	}
}

// Error 打印 error 日志
// 参数：
//   - ctx: context.Context
//   - s: string
//   - args:...interface{}
//
// 返回值：
//   - 无
func (l *GormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	if l.level >= logger.Error {
		_log := loggers[DbError]
		_log.Error().Msgf(s, args...)
	}
}

// Trace 打印 trace 日志
// 参数：
//   - ctx: context.Context
//   - begin: time.Time
//   - fc: func() (string, int64)
//   - err: error
//
// 返回值：
//   - 无
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.level >= logger.Error:
		_log := loggers[DbError]
		_log.Error().Err(err).Msgf("[%.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)

	case elapsed > 200*time.Millisecond && l.level >= logger.Warn: // 慢查询阈值可以调
		_log := loggers[DbWarn]
		_log.Warn().Msgf("[SLOW SQL >=200ms] [%.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)

	case l.level >= logger.Info:
		_log := loggers[DbInfo]
		_log.Info().Msgf("[%.3fms] [rows:%v] %s", float64(elapsed.Milliseconds()), rows, sql)
	}
}
