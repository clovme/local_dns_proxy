package copyright

import (
	"fmt"
	"local_dns_proxy/pkg/constants"
	"local_dns_proxy/version"
	"os"
	"runtime"
)

// _copyright 版权信息
type _copyright struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
	Pid       int    `json:"pid"`
	BuildTime string `json:"buildTime"`
}

func NewCopyright() _copyright {
	return _copyright{
		Name:      constants.WebTitle,
		Version:   version.Version,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		Pid:       os.Getpid(),
		BuildTime: version.BuildTime,
	}
}
