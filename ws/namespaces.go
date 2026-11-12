package ws

import (
	"encoding/json"
	"log"
	"net/http"

	logger "github.com/taigrr/log-socket/v2/log"
)

// NamespacesHandler returns a JSON list of all namespaces that have been used
func NamespacesHandler(w http.ResponseWriter, r *http.Request) {
	namespaces := logger.GetNamespaces()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"namespaces": namespaces,
	}); err != nil {
		// A failed write may have committed the response already.
		log.Printf("log-socket: write namespaces: %v", err)
	}
}
