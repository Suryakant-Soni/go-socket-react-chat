package main

import (
	"go-socket-react-chat/pkg/websocket"
	"log"
	"net/http"
)

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	// whenever somebody hits the ws endpoint we need to upgrade the connection first from http to ws protocol
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		log.Print("Error", err)
	}
	// create the client for the person who has hit the url
	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}
	// add this client to the pool
	pool.Register <- client
	// read the message coming from this new client within this url call
	client.Read()
}

func setupRoutes() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("ws handler triggered")
		serveWS(pool, w, r)
	})
}
func main() {

	setupRoutes()
	log.Println("llisten and serve at 9000")
	http.ListenAndServe(":9000", nil)
}
