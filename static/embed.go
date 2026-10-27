package static

import (
	"embed"
)

//go:embed css/* images/* generated/*
var Dir embed.FS
