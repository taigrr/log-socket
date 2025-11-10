package ws

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	logger "github.com/taigrr/log-socket/v2/log"
)

var upgrader = websocket.Upgrader{} // use default options

func SetUpgrader(u websocket.Upgrader) {
	upgrader = u
}

func LogSocketHandler(w http.ResponseWriter, r *http.Request) {
	// Get namespaces from query parameter, comma-separated
	// Empty or missing means all namespaces
	namespacesParam := r.URL.Query().Get("namespaces")
	var namespaces []string
	if namespacesParam != "" {
		namespaces = strings.Split(namespacesParam, ",")
	}
	
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("upgrade:", err)
		return
	}
	defer c.Close()
	lc := logger.CreateClient(namespaces...)
	defer lc.Destroy()
	lc.SetLogLevel(logger.LTrace)
	logger.Info("Websocket client attached.")
	for {
		logEvent := lc.Get()
		logJSON, _ := json.Marshal(logEvent)
		err = c.WriteMessage(websocket.TextMessage, logJSON)
		if err != nil {
			logger.Warn("write:", err)
			break
		}
	}
}
