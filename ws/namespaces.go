package ws

import (
	"encoding/json"
	"net/http"

	logger "github.com/taigrr/log-socket/v2/log"
)

// NamespacesHandler returns a JSON list of all namespaces that have been used
func NamespacesHandler(w http.ResponseWriter, r *http.Request) {
	namespaces := logger.GetNamespaces()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"namespaces": namespaces,
	})
}
