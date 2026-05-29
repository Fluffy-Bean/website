package data

import (
	"embed"
)

//go:embed assets/* blogs/* art.json
var Dir embed.FS
