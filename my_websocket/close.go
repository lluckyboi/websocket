package my_websocket

import "log"

func (conn *MyConn) Close() {
	err := conn.conn.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("websocket connection closed")
}
