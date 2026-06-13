package data

import (
	"embed"
)

//go:embed blogs/* art.json
var Dir embed.FS
