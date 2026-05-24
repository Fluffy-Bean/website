package blog

import (
	"bytes"
	"time"
)

type Blog struct {
	Slug        string
	Title       string
	PublishedAt time.Time
	Data        bytes.Buffer
}
