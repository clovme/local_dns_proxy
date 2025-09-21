package public

import (
	"embed"
)

//go:embed web
var WebFS embed.FS

//go:embed favicon.ico
var Favicon []byte
