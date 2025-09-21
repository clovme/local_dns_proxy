package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// LoggerConfig 日志配置
//
//	type LoggerConfig struct {
//			Dir        string // 日志目录，如 ./logs
//			MaxSize    int    // MB
//			MaxBackups int    // 个数
//			MaxAge     int    // 天
//			Compress   bool	  // 是否压缩
//			FormatJSON bool   // 是否 JSON 格式
//			Level      string // 最低输出级别，如 debug、info、error
//			Lvl        zerolog.Level
//	}
type LoggerConfig struct {
	Dir        string // 日志目录，如 ./logs
	MaxSize    int    // MB
	MaxBackups int    // 个数
	MaxAge     int    // 天
	Compress   bool   // 是否压缩
	FormatJSON bool   // 是否 JSON 格式
	Level      string // 最低输出级别，如 debug、info、error
	Lvl        zerolog.Level
}

var (
	consoleWriter io.Writer                 // 控制台输出
	loggers       map[string]zerolog.Logger // 日志记录器
	mu            sync.RWMutex              // 读写锁
	initialized   bool                      // 初始化标志
	CurrentCfg    *LoggerConfig             // 当前配置
)

var loc, _ = time.LoadLocation("Asia/Shanghai") // 设置为上海时间

// 格式化时间
func formatTimestamp(i interface{}) string {
	return time.Now().In(loc).Format("[2006-01-02 15:04:05]")
}

// 格式化日志级别
func formatLevel(i interface{}) string {
	return fmt.Sprintf("[%s]", i)
}

// InitLogger 初始化日志
// 参数：
//   - cfg: LoggerConfig
//
// 返回值：
//   - error: 初始化错误
func InitLogger(cfg LoggerConfig) {
	CurrentCfg = &cfg
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return
	}

	lvl, err := zerolog.ParseLevel(cfg.Level)
	if err != nil {
		lvl = zerolog.InfoLevel
		fmt.Printf("[日志初始化] 日志级别[%s]无效，使用默认级别[%s]\n", cfg.Level, zerolog.InfoLevel.String())
	}

	cfg.Lvl = lvl

	zerolog.LevelFieldMarshalFunc = func(l zerolog.Level) string {
		return strings.ToUpper(l.String())
	}
	zerolog.TimeFieldFormat = "[2006-01-02 15:04:05]"
	if cfg.FormatJSON {
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	}

	zerolog.TimestampFunc = func() time.Time {
		return time.Now().In(loc)
	}

	// 控制台输出
	consoleWriter = zerolog.ConsoleWriter{
		Out:             os.Stdout,
		NoColor:         false,
		FormatTimestamp: formatTimestamp,
		FormatLevel:     formatLevel,
	}

	loggers = make(map[string]zerolog.Logger)

	levels := []zerolog.Level{
		zerolog.DebugLevel,
		zerolog.InfoLevel,
		zerolog.WarnLevel,
		zerolog.ErrorLevel,
		zerolog.FatalLevel,
		zerolog.PanicLevel,
		zerolog.Disabled,
		zerolog.TraceLevel,
	}

	for _, level := range levels {
		for _, t := range []string{"app", "http", "db"} {
			fileName := filepath.Join(cfg.Dir, level.String(), fmt.Sprintf("%s.log", t))
			fileWriter := &lumberjack.Logger{
				Filename:   fileName,
				MaxSize:    cfg.MaxSize,
				MaxBackups: cfg.MaxBackups,
				MaxAge:     cfg.MaxAge,
				Compress:   cfg.Compress,
			}

			var writer io.Writer

			if cfg.FormatJSON {
				// JSON 格式：直接复用 zerolog 默认 JSON 格式输出
				writer = io.MultiWriter(consoleWriter, fileWriter)
			} else {
				// ConsoleWriter 格式：也要给 file 和 WebSocket 同样的格式化输出
				textWriter := zerolog.ConsoleWriter{
					Out:             fileWriter,
					NoColor:         true,
					FormatTimestamp: formatTimestamp,
					FormatLevel:     formatLevel,
				}
				writer = io.MultiWriter(consoleWriter, textWriter)
			}

			_callerSkipFrames := 2
			if skip, ok := callerSkipFrames[t]; ok {
				_callerSkipFrames = skip
			}
			if cfg.Level != "info" && cfg.Level != "warn" {
				loggers[fmt.Sprintf("%s_%s", t, level.String())] = zerolog.New(writer).Level(lvl).With().CallerWithSkipFrameCount(_callerSkipFrames).Timestamp().Logger()
			} else {
				loggers[fmt.Sprintf("%s_%s", t, level.String())] = zerolog.New(writer).Level(lvl).With().Timestamp().Logger()
			}
		}
	}

	initialized = true
}

// GetLogger 获取指定文件名的日志记录器
// 参数：
//   - filename: 日志文件名
//
// 返回值：
//   - zerolog.Logger: 日志记录器
func GetLogger(filename string) zerolog.Logger {
	return loggers[filename]
}
