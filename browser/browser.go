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
	wsResource := websocketScheme(r) + r.Host + r.URL.Path
	wsResource = strings.TrimSuffix(wsResource, "/") + "/ws"
	homeTemplate.Execute(w, wsResource)
}

func websocketScheme(r *http.Request) string {
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		protocol := strings.ToLower(strings.TrimSpace(strings.Split(forwardedProto, ",")[0]))
		if protocol == "https" {
			return "wss://"
		}
		if protocol == "http" {
			return "ws://"
		}
	}
	if r.TLS != nil {
		return "wss://"
	}
	return "ws://"
}

var homeTemplate = template.Must(template.New("").Parse(webpage))
