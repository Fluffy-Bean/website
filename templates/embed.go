package templates

import (
	"embed"
)

//go:embed blocks/* pages/* layout.html
var Dir embed.FS
