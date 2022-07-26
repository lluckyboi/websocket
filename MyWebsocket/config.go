package MyWebsocket

import (
	"net"
	"net/http"
	"time"
)

const (
	TextMessage        = 1
	BinaryMessage      = 2
	FileImageMessage   = 3
	CloseMessage       = 8
	PingMessage        = 9
	PongMessage        = 10
	WSVersion          = "13"
	MagicString        = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
	DefaultTimeOut     = time.Second * 180
	DefaultReadBuffer  = 65535 + 125 + 1
	DefaultWriteBuffer = 65535 + 125 + 1
	DefaultPingWait    = 30 * time.Second
)

type Msg struct {
	Typ     int
	content []byte
}

type MyConn struct {
	Conn                            net.Conn
	ReadBufferSize, WriteBufferSize uint64
	PingTimeOut                     func() time.Time
	Opts                            ConnOptions
}

type Upgrader struct {
	// 指定升级 websocket 握手完成的超时时间
	HandshakeTimeout time.Duration

	// 指定 io 操作的缓存大小
	ReadBufferSize, WriteBufferSize uint64

	// 指定 http 的错误响应函数，如果没有设置 Error 则，会生成 http.Error 的错误响应。
	Error func(w http.ResponseWriter, r *http.Request, status int, reason error)

	//自定义检查源函数 默认返回true
	CheckOrigin func(r *http.Request) bool
}

type Writer struct {
	idx      uint64 //记录当前传输位置
	datast   uint64 //数据开始的下标
	maskst   uint64 //maskKey开始的下标
	maskKey  []byte //maskKey
	restDate uint64 //剩余数据大小
	ismain   bool   //是否主片
}

type ConnOptions struct {
	WriteTimeOut time.Duration
	PingWait     time.Duration
	PongHandler  PongHandler
	IOLog        bool
}

type Option func(*ConnOptions)

type PongHandler func()
