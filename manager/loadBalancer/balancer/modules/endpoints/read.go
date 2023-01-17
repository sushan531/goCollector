package endpoints

import (
	"github.com/gorilla/websocket"
	"github.com/pebbe/zmq4"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

type Handler struct {
	*zmq4.Socket
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, message, err := wsConn.ReadMessage()
		if err != nil {
			return
		}
		_, _ = h.SendMessage(message)

	}
}

func LogHandler(server *zmq4.Socket) Handler {
	var handler Handler
	handler.Socket = server
	return handler
}
