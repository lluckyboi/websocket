package my_websocket

func (conn *MyConn) Close() {
	//发送关闭数据帧
	p := make([]byte, 2)
	p[0] = 136
	conn.conn.Write(p)
}
