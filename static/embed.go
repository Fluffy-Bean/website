package static

import (
	"embed"
)

//go:embed css/* images/* fonts/* generated/*
var Dir embed.FS
