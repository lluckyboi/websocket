package MyWebsocket

import "log"

func (conn *MyConn) pingHandler() {
	//刷新心跳
	err := conn.Conn.SetDeadline(conn.PingTimeOut())
	if err != nil {
		log.Println("refresh ping err:", err)
	}
	//发送pong
	p := make([]byte, 2)
	p[0] = 1<<7 + 10
	_, err = conn.Conn.Write(p)
	if err != nil {
		log.Println("send pong:", err)
	}
}

func (conn *MyConn) pongHandler() {
	conn.Opts.PongHandler()
}
