package main

import (
	"log"
	"websocket/MyWebsocket"
)

func main() {
	conn, err := MyWebsocket.Dail("127.0.0.1", "9924", "http://127.0.0.1:9924/ws")
	if err != nil {
		log.Println(err)
	}
	err = conn.Write([]byte("ws"), 1)
	if err != nil {
		log.Println(err)
	}
	_, ms, err := conn.ReadMsg()
	if err != nil {
		log.Println(err)
	}
	log.Println(ms)
}
