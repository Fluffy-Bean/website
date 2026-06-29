package blog

import (
	"bytes"
	"time"
)

var oldBlogTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

type Blog struct {
	Slug        string
	Title       string
	PublishedAt time.Time
	Data        bytes.Buffer
}

func (b Blog) IsOldBlog() bool {
	return b.PublishedAt.Before(oldBlogTime)
}
