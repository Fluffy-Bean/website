package static

import (
	"embed"
)

//go:embed css/* images/*
var Dir embed.FS
