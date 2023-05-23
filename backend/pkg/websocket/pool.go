package websocket

import "log"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// to start working/functioning of the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Println("no of clients within the pool", pool.Clients)
			for client, _ := range pool.Clients {
				log.Println("client details", client)
				client.Conn.WriteJSON(Message{Type: 1, Body: "New User joined..."})
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.Println("size of connection pool", len(pool.Clients))
			for client := range pool.Clients {
				client.Conn.WriteJSON(Message{Type: 1, Body: "User disconnected"})
			}
		case message := <-pool.Broadcast:
			log.Println("message getting broadcasted", message)
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Println("errir", err.Error())
					return
				}
			}
		}
	}
}
