package websockets

import "log"

// Dispatcher maintains the set of active clients and broadcasts messages to the
// clients.
type Dispatcher struct {
	// Broadcase messages to all client.
	Broadcast chan interface{}

	// Registered clients.
	clients map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewDispatcher creates a new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		Broadcast:  make(chan interface{}),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Register returns the register channel
func (d *Dispatcher) Register() chan *Client {
	return d.register
}

// Run starts the dispatch loop
func (d *Dispatcher) Run() {
	for {
		select {
		case client := <-d.register:
			log.Printf("registered new client")
			d.clients[client] = true
		case client := <-d.unregister:
			log.Printf("unregistered client")
			if _, ok := d.clients[client]; ok {
				delete(d.clients, client)
				close(client.send)
			}
		case message := <-d.Broadcast:
			for client := range d.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(d.clients, client)
				}
			}
		}
	}
}
