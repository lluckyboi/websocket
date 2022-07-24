package my_websocket

import (
	"errors"
	"net/http"
	"time"
)

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (conn *MyConn, err error) {
	conn = &MyConn{}
	//设置默认值
	if u.ReadBufferSize == 0 {
		u.ReadBufferSize = DefaultReadBuffer
	}
	if u.WriteBufferSize == 0 {
		u.WriteBufferSize = DefaultReadBuffer
	}
	if u.HandshakeTimeout == time.Duration(0) {
		u.HandshakeTimeout = DefaultTimeOut
	}

	conn.ReadBufferSize = u.ReadBufferSize
	conn.WriteBufferSize = u.WriteBufferSize

	//检查方法
	if r.Method != http.MethodGet {
		return &MyConn{}, errors.New("the method is not get")
	}
	//检查请求头
	if r.Header.Get("Connection") != "Upgrade" {
		return &MyConn{}, errors.New("the connection is not Upgrade")
	}
	if r.Header.Get("Upgrade") != "websocket" {
		return &MyConn{}, errors.New("the connection is not websocket")
	}
	if r.Header.Get("Sec-Websocket-Version") != WSVersion {
		return &MyConn{}, errors.New("the version is not 13")
	}

	//从http.ResponseWriter重新拿到conn 出错就返回
	hijcakcer, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return &MyConn{}, errors.New("upgrade conn err:get conn")
	}
	conn.conn, _, err = hijcakcer.Hijack()

	//拿到浏览器生成的密钥 并与Websocket的Magic String拼接
	wskey := append([]byte(r.Header.Get("Sec-Websocket-Key")), []byte(MagicString)...)
	respAccept := SHA1AndBase64(string(wskey))

	//回复报文 超时返回
	resp := []byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + respAccept + "\r\n\r\n")
	err = conn.conn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	if err != nil {
		return nil, err
	}
	_, err = conn.conn.Write(resp)
	if err != nil {
		return nil, err
	}
	return
}
