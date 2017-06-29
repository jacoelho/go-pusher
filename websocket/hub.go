package websocket

import (
	"log"
	"net/http"

	ws "github.com/gorilla/websocket"
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type hub struct {
	// Registered Clients.
	Clients map[*client]bool

	// Inbound messages from the Clients.
	broadcast chan []byte

	// Register requests from the Clients.
	register chan *client

	// Unregister requests from Clients.
	unregister chan *client

	// Logger
	log *log.Logger
}

func NewHub(l *log.Logger) *hub {
	return &hub{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		Clients:    make(map[*client]bool),
		log:        l,
	}
}

func (h *hub) RunWithEmiter(ch <-chan []byte) {
	for {
		select {
		case client := <-h.register:
			h.Clients[client] = true

		case client := <-h.unregister:
			delete(h.Clients, client)
			close(client.send)

		case message := <-ch:
			for client := range h.Clients {
				select {
				case client.send <- message:
				default:
					// send failed, close
					h.log.Println("failed send")
					close(client.send)
					delete(h.Clients, client)
				}
			}
		}
	}
}

func (h *hub) ServeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.log.Println(err)
		return

	}
	client := &client{
		hub:  h,
		conn: conn,
		send: make(chan []byte, 256),
		log:  h.log,
	}

	h.register <- client

	go client.Write()
	go client.Read()
}
