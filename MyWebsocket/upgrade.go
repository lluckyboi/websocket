package MyWebsocket

import (
	"errors"
	"net/http"
	"time"
)

func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request, opts ...Option) (conn *MyConn, err error) {
	conn = &MyConn{}
	op := ConnOptions{}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.PingWait = op.PingWait
	conn.Opts.IOLog = op.IOLog
	conn.Opts.PongHandler = op.PongHandler

	if conn.Opts.PingWait == time.Duration(0) {
		conn.Opts.PingWait = DefaultPingWait
	}
	if conn.Opts.PongHandler == nil {
		conn.Opts.PongHandler = func() {
			return
		}
	}

	conn.PingTimeOut = func() time.Time {
		return time.Now().Add(conn.Opts.PingWait)
	}

	//设置默认值
	u.ReadBufferSize = DefaultReadBuffer

	u.WriteBufferSize = DefaultWriteBuffer

	u.HandshakeTimeout = DefaultTimeOut

	if u.CheckOrigin == nil {
		u.CheckOrigin = func(r *http.Request) bool {
			return true
		}
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
	//check-origin
	if !u.CheckOrigin(r) {
		return &MyConn{}, errors.New("check origin false")
	}

	//从http.ResponseWriter重新拿到conn 出错就返回
	hijcakcer, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, http.StatusText(500), 500)
		return &MyConn{}, errors.New("upgrade Conn err:get Conn")
	}
	conn.Conn, _, err = hijcakcer.Hijack()

	//拿到浏览器生成的密钥 并与Websocket的Magic String拼接
	wskey := append([]byte(r.Header.Get("Sec-Websocket-Key")), []byte(MagicString)...)
	respAccept := SHA1AndBase64(string(wskey))

	//回复报文 超时返回
	resp := []byte("HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\nConnection: Upgrade\r\nSec-WebSocket-Accept: " + respAccept + "\r\n\r\n")
	err = conn.Conn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	if err != nil {
		return nil, err
	}
	_, err = conn.Conn.Write(resp)
	if err != nil {
		return nil, err
	}
	return
}
