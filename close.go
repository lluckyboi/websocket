package websocket

import "log"

func (conn *MyConn) Close() {
	err := conn.Conn.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("websocket connection closed successfully")
}
