package websockets

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// upgrader upgrades the request to WS
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handlers handles websocket connections and provides an interface to dispatch events
type Handlers struct {
	dispatcher *Dispatcher
	serveWs    func(*Dispatcher, http.ResponseWriter, *http.Request, string)
}

// NewHandlers creates and returns a pointer to a new instance of `Handler`
func NewHandlers(dispatcher *Dispatcher) *Handlers {
	return &Handlers{
		dispatcher: dispatcher,
		serveWs:    serveWs,
	}
}

// WebsocketHandler handles all of the websocket connections
func (h *Handlers) WebsocketHandler(c *gin.Context) {
	channel := c.Param("channel")
	// TODO: check if channel(project) exists
	h.serveWs(h.dispatcher, c.Writer, c.Request, channel)
}

// serveWs handles websocket requests from the peer.
func serveWs(dispatcher *Dispatcher, w http.ResponseWriter, r *http.Request, channel string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	client := NewClient(dispatcher, conn, channel)
	dispatcher.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WriteDispatch()
	go client.ReadDispatch()
}
