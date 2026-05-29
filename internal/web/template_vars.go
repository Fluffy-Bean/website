package web

import (
	"net/http"
)

func tmplURLVars(r *http.Request) Data {
	return Data{
		"Path":   r.URL.Path,
		"URL":    r.URL.String(),
		"Method": r.Method,
	}
}
