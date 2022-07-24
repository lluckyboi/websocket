package my_websocket

import "log"

func (conn *MyConn) pingHandler() {
	//刷新心跳
	conn.PingTimeOut.Add(conn.Opts.PingWait)
	//发送pong
	p := make([]byte, 2)
	p[0] = 1<<7 + 10
	_, err := conn.conn.Write(p)
	if err != nil {
		log.Println(err)
	}
}

func (conn *MyConn) pongHandler() {
	conn.Opts.PongHandler()
}
