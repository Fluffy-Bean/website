package web

import (
	"bytes"
	"fmt"
	"html/template"
	"path"
	"time"

	"github.com/yuin/goldmark"
)

var templateFuncs = template.FuncMap{
	"props":           templProps,
	"noescape":        templNoEscape,
	"markdown":        templMarkdown,
	"format_time":     templFormatTime,
	"format_duration": templFormatDuration,
	"path_join":       templPathJoin,
	"format_filesize": templFormatFileSize,
}

func templProps(args ...interface{}) map[string]interface{} {
	props := make(map[string]interface{})

	if len(args)%2 != 0 {
		panic("args must have even number of parameters")
	}

	for i := 0; i < len(args); i += 2 {
		props[args[i].(string)] = args[i+1]
	}

	return props
}

func templNoEscape(s string) template.HTML {
	return template.HTML(s)
}

func templMarkdown(s string) template.HTML {
	var buff bytes.Buffer

	err := goldmark.Convert([]byte(s), &buff)
	if err != nil {
		panic(err)
	}

	return template.HTML(buff.String())
}

func templFormatTime(t time.Time, format string) string {
	switch format {
	case "UnixDate":
		return t.Format(time.UnixDate)
	case "RFC822":
		return t.Format(time.RFC822)
	case "RFC3339":
		return t.Format(time.RFC3339)
	case "Kitchen":
		return t.Format(time.Kitchen)
	case "Stamp":
		return t.Format(time.Stamp)
	case "DateTime":
		return t.Format(time.DateTime)
	case "DateOnly":
		return t.Format(time.DateOnly)
	case "TimeOnly":
		return t.Format(time.TimeOnly)
	default:
		return t.Format(format)
	}
}

func templFormatDuration(d time.Duration) string {
	var format string

	if v := int(d.Hours()); v > 0 {
		d -= time.Duration(v) * time.Hour
		format += fmt.Sprintf("%dh ", v)
	}

	if v := int(d.Minutes()); v > 0 {
		d -= time.Duration(v) * time.Minute
		format += fmt.Sprintf("%dm ", v)
	}

	format += fmt.Sprintf("%ds", int(d.Seconds()))

	return format
}

func templPathJoin(elem ...string) string {
	return path.Join(elem...)
}

func templFormatFileSize(size int64) string {
	threshold := float64(1024)
	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}

	s := float64(size)
	i := 0

	for s >= threshold && i < len(units) {
		s /= threshold
		i += 1
	}

	if i > 1 {
		return fmt.Sprintf("%.2f%s", s, units[i])
	}

	return fmt.Sprintf("%.0f%s", s, units[i])
}
