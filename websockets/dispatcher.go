package websockets

import "log"

// Message wraps the relevant information needed to broadcast a message
type Message struct {
	Channel string      // the channel name to broadcast on
	Data    interface{} // the data to broadcast
}

// Dispatcher maintains the set of active clients and broadcasts messages to the
// clients.
type Dispatcher struct {
	// Broadcase messages to all client.
	broadcast chan *Message

	// Registered clients.
	clients map[string]map[*Client]bool

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewDispatcher creates a new Dispatcher
func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

// Broadcast returns the broadcast channel
func (d *Dispatcher) Broadcast() chan *Message {
	// NOTE: To scale this, a message queue must be used to "broadcast" messages
	// to all running webservers. This function would enqueue messages to the
	// the message queue system, while the dispatcher also reads messages from the queue
	// and sends them to the `broadcast` channel
	return d.broadcast
}

// Run starts the dispatch loop
func (d *Dispatcher) Run() {
	for {
		select {
		case client := <-d.register:
			channel := client.channel
			log.Printf("registered new client to '%s'", channel)
			if _, ok := d.clients[channel]; !ok {
				d.clients[channel] = make(map[*Client]bool)
			}
			d.clients[channel][client] = true
		case client := <-d.unregister:
			channel := client.channel
			log.Printf("unregistered client from '%s'", channel)
			if _, ok := d.clients[channel][client]; ok {
				delete(d.clients[channel], client)
				close(client.send)
			}
		case message := <-d.broadcast:
			channel := message.Channel
			if clients, ok := d.clients[channel]; ok {
				for client := range clients {
					select {
					case client.send <- message.Data:
					default:
						close(client.send)
						delete(d.clients[channel], client)
					}
				}
			}
		}
	}
}
