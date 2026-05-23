package blog

import (
	"bytes"
)

type Blog struct {
	Title string
	Data  bytes.Buffer
}
