package version

import "strings"

var (
	BuildTime = "0000-00-00_00:00:00"
	Version   = "v0.0.0"
)

func init() {
	BuildTime = strings.ReplaceAll(BuildTime, "_", " ")
}
