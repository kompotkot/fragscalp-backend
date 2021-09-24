package streamglass

import (
	"log"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// GET handler for HTTP
func GetHandler(w http.ResponseWriter, r *http.Request) {
	filename := r.URL.EscapedPath()[1:]
	fmt.Println(filename)
}

// HTTP handler
func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		GetHandler(w, r)
		return
	}
	http.NotFound(w, r)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	/*
		WebSocket connection handler.

		:w: (http.ResponseWriter)
		:r: (*http.Request)
	*/
	// Upgrader will require a Read and Write buffer size
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// CORS checks
	// TODO(kompotkot): Now we just returning true, in future modify it
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// Upgrade this connection to a WebSocket connection
	ws, wsErr := upgrader.Upgrade(w, r, nil)
	if wsErr != nil {
		log.Println(wsErr)
		return
	}

	chWs := make(chan []byte, 1)
	pair := BINANCE_SOLUSDT
	depth := 20
	updateSpeed := 100
	go Extractor(chWs, pair, depth, updateSpeed)

	for {
		select {
		case val := <-chWs:
			wsErr = ws.WriteMessage(websocket.TextMessage, val)
			if wsErr != nil {
				log.Println(wsErr)
			}
		}
	}
}

func Server(host string, port string) {
	uri := fmt.Sprintf("%s:%s", host, port)
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Starting server: %s\n", uri)

	http.HandleFunc("/", httpHandler)
	http.HandleFunc("/ws", wsHandler)
	
	http.ListenAndServe(uri, nil)
}
