package MyWebsocket

import "log"

func (conn *MyConn) Close() {
	conn.Write([]byte("close"), 8)
	err := conn.Conn.Close()
	if err != nil {
		log.Println(err)
	}
	log.Println("websocket connection closed successfully")
}
