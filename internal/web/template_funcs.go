package web

import (
	"bytes"
	"html/template"

	"github.com/yuin/goldmark"
)

var templateFuncs = template.FuncMap{
	"props":    templProps,
	"noescape": templNoEscape,
	"markdown": templMarkdown,
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
