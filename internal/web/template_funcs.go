package web

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"math/rand/v2"
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
	"hsl_to_rgb":      templHSLToRGB,
	"rand":            templRand,
	"add":             templAdd,
	"sub":             templSub,
	"mul":             templMul,
	"div":             templDiv,
	"mod":             templMod,
}

func templProps(args ...any) map[string]any {
	props := make(map[string]any)

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

// https://github.com/Crazy3lf/colorconv/blob/master/colorconv.go#L97
func templHSLToRGB(h, s, l float64) ([]int, error) {
	if h < 0 || h >= 360 || s < 0 || s > 1 || l < 0 || l > 1 {
		return nil, fmt.Errorf("out of range")
	}

	// When 0 ≤ h < 360, 0 ≤ s ≤ 1 and 0 ≤ l ≤ 1:
	C := (1 - math.Abs((2*l)-1)) * s
	X := C * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := l - (C / 2)
	var Rnot, Gnot, Bnot float64

	switch {
	case 0 <= h && h < 60:
		Rnot, Gnot, Bnot = C, X, 0
	case 60 <= h && h < 120:
		Rnot, Gnot, Bnot = X, C, 0
	case 120 <= h && h < 180:
		Rnot, Gnot, Bnot = 0, C, X
	case 180 <= h && h < 240:
		Rnot, Gnot, Bnot = 0, X, C
	case 240 <= h && h < 300:
		Rnot, Gnot, Bnot = X, 0, C
	case 300 <= h && h < 360:
		Rnot, Gnot, Bnot = C, 0, X
	}

	r := int(math.Round((Rnot + m) * 255))
	g := int(math.Round((Gnot + m) * 255))
	b := int(math.Round((Bnot + m) * 255))

	return []int{r, g, b}, nil
}

func templRand(max float64) float64 {
	return rand.Float64() * max
}

func templAdd(lhs, rhs float64) float64 {
	return lhs + rhs
}

func templSub(lhs, rhs float64) float64 {
	return lhs - rhs
}

func templMul(lhs, rhs float64) float64 {
	return lhs * rhs
}

func templDiv(lhs, rhs float64) float64 {
	return lhs / rhs
}

func templMod(lhs, rhs float64) float64 {
	return math.Mod(lhs, rhs)
}
