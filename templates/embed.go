package templates

import (
	"embed"
)

//go:embed blocks/* icons/* pages/* layout.html
var Dir embed.FS
