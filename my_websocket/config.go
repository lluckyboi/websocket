package my_websocket

import (
	"net"
	"net/http"
	"time"
)

const (
	TextMessage       = 1
	BinaryMessage     = 2
	CloseMessage      = 8
	PingMessage       = 9
	PongMessage       = 10
	WSVersion         = "13"
	MagicString       = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	DefaultTimeOut    = time.Second * 180
	DefaultReadBuffer = 2048
)

type Handler func(conn *MyConn)

type Msg struct {
	Typ     int
	content []byte
}

type MyConn struct {
	conn                            net.Conn
	ReadBufferSize, WriteBufferSize int
	//handler
}

type Upgrader struct {
	// 指定升级 websocket 握手完成的超时时间
	HandshakeTimeout time.Duration

	// 指定 io 操作的缓存大小
	ReadBufferSize, WriteBufferSize int

	// 指定 http 的错误响应函数，如果没有设置 Error 则，会生成 http.Error 的错误响应。
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)
}

type Writer struct {
	idx      int    //记录当前传输位置
	datast   int    //数据开始的下标
	maskst   int    //maskKey开始的下标
	maskKey  []byte //maskKey
	restDate int    //剩余数据大小
	ismain   bool   //是否主片
}
