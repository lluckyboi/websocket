package websocket

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
	"math/rand"
	"time"
)

// SHA1AndBase64 按照协议进行sha1+base64 进行加密
func SHA1AndBase64(data string) string {
	t := sha1.New()
	io.WriteString(t, data)
	//进行base64编码
	res := base64.StdEncoding.EncodeToString(t.Sum(nil))
	return res
}

// CreatMuskKey 随机生成MuskKey
func CreatMuskKey() []byte {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
		//休眠一纳秒
		time.Sleep(time.Nanosecond)
	}
	return []byte(string(b))
}

func WithWriteTimeout(timeout time.Duration) Option {
	return func(o *ConnOptions) {
		o.WriteTimeOut = timeout
	}
}

func WithPingWait(timeout time.Duration) Option {
	return func(o *ConnOptions) {
		o.PingWait = timeout
	}
}

func WithPongHandler(handler PongHandler) Option {
	return func(options *ConnOptions) {
		options.PongHandler = handler
	}
}

func WithIOLOG(need bool) Option {
	return func(options *ConnOptions) {
		options.IOLog = need
	}
}

func (conn *MyConn) SetIOLog(need bool) {
	conn.Opts.IOLog = need
}

func (conn *MyConn) SetWriteBuffersize(size uint64) {
	conn.WriteBufferSize = size
}

func (conn *MyConn) SetReadBuffersize(size uint64) {
	conn.ReadBufferSize = size
}
