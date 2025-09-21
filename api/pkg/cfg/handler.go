package cfg

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/pkg/utils"
	"local_dns_proxy/pkg/utils/file"
	"os"
	"strings"
)

var (
	C    *Config
	Path string
)

func init() {
	ifaces, err := utils.GetNetworkInterfaces()
	if err != nil {
		ifaces[0].Name = "以太网"
	}
	C = &Config{
		Server: Server{
			Port:       6500,
			DNSRunning: constants.DNSStop,
			Iface:      ifaces[0].Name,
		},
		Logger: Logger{
			Level:      zerolog.InfoLevel.String(),
			MaxSize:    50,
			MaxAge:     7,
			MaxBackups: 5,
		},
	}

	// ini 覆盖
	_path, err := file.GetFileAbsPath(".", fmt.Sprintf("%s_conf.ini", constants.ProjectName))
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	} else {
		Path = _path
	}

	if _file, err := ini.Load(_path); err == nil {
		_ = _file.MapTo(C)
	}

	if C.Logger.Level == "no" {
		C.Logger.Level = ""
	}
	C.Logger.Level = strings.ToLower(C.Logger.Level)
}

// SaveToIni 保存配置到 ini 文件
func SaveToIni() {
	_file := ini.Empty()
	err := _file.ReflectFrom(C)
	if err != nil {
		log.Fatal().Err(err).Msg("配置保存，序列化成ini失败")
	}

	if _file.SaveTo(Path) != nil {
		log.Fatal().Err(err).Msg("配置文件保存失败")
	}
}
