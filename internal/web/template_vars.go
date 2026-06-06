package web

import (
	"net/http"
	"time"
)

func tmplURLVars(r *http.Request) Data {
	return Data{
		"Path":   r.URL.Path,
		"URL":    r.URL.String(),
		"Method": r.Method,
		"Now":    time.Now(),
	}
}

func tmplTimeVars() Data {
	return Data{
		"Now": time.Now(),
	}
}
