package boot

import (
	"fmt"
	"local_dns_proxy/pkg/cfg"
	"local_dns_proxy/pkg/logger"
	"local_dns_proxy/pkg/utils/file"
	"os"
)

// InitializationLogger 初始化日志
func InitializationLogger(c cfg.Logger) {
	path, err := file.GetFileAbsPath("logs")
	if err != nil {
		fmt.Println("获取日志目录失败:", err)
		os.Exit(-1)
	}
	// 初始化一次
	logger.InitLogger(logger.LoggerConfig{
		Dir:        path,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   true,
		Level:      c.Level,
		FormatJSON: false, // true=结构化；false=文本
	})
}
