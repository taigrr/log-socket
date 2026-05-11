package browser

import (
	_ "embed"
	"html/template"
	"net/http"
	"strings"
)

//go:embed viewer.html
var webpage string

func LogSocketViewHandler(w http.ResponseWriter, r *http.Request) {
	wsResource := websocketScheme(r) + websocketHost(r) + r.URL.Path
	wsResource = strings.TrimSuffix(wsResource, "/") + "/ws"
	homeTemplate.Execute(w, wsResource)
}

func websocketScheme(r *http.Request) string {
	switch strings.ToLower(forwardedHeaderValue(r, "X-Forwarded-Proto")) {
	case "https":
		return "wss://"
	case "http":
		return "ws://"
	}
	if r.TLS != nil {
		return "wss://"
	}
	return "ws://"
}

func websocketHost(r *http.Request) string {
	if host := forwardedHeaderValue(r, "X-Forwarded-Host"); host != "" {
		return host
	}
	return r.Host
}

func forwardedHeaderValue(r *http.Request, key string) string {
	return strings.TrimSpace(strings.Split(r.Header.Get(key), ",")[0])
}

var homeTemplate = template.Must(template.New("").Parse(webpage))
