package cfg

type Server struct {
	Port       int    `ini:"port" json:"port" yaml:"port" comment:"Web 服务监听端口"`
	Iface      string `ini:"iface" json:"iface" yaml:"iface" comment:"当前网卡名称"`
	DNSRunning string `ini:"-" json:"-" yaml:"-" comment:"是否启用 DNS 代理服务"`
}

type Logger struct {
	Level      string `ini:"level" json:"level" yaml:"level" comment:"数据库日志级别 info > warn > error > silent  silent 不记录任何日志，相当于disabled\n; 系统日志级别   trace > debug > info > warn > error > fatal > panic > no > disabled\n; trace\t\t细粒度最高，最大量日志\n; debug\t\t调试日志\n; info\t\t常规运行状态日志\n; warn\t\t警告，非致命异常\n; error\t\t错误日志，功能异常\n; fatal\t\t致命错误，程序即将终止\n; panic\t\t更严重，触发 panic 行为\n; no\t\t没有级别，适合特殊用途\n; disabled\t禁止所有日志"`
	MaxSize    int    `ini:"max_size" json:"max_size" yaml:"max_size" comment:"单个日志文件最大尺寸，单位为 MB，超过该大小将触发日志切割"`
	MaxAge     int    `ini:"max_age" json:"max_age" yaml:"max_age" comment:"日志文件最大保存天数，超过该天数的日志文件将被删除"`
	MaxBackups int    `ini:"max_backups" json:"max_backups" yaml:"max_backups" comment:"保留旧日志文件的最大数量，超过时自动删除最早的日志"`
}

type Config struct {
	Server Server `ini:"Server" json:"Server" yaml:"Server"`
	Logger Logger `ini:"Logger" json:"Logger" yaml:"Logger"`
}
