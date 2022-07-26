package websocket

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

func Dail(host, port, Url string) (*MyConn, error) {
	u, err := url.Parse(Url)
	if err != nil {
		return nil, err
	}

	req := &http.Request{
		Method:     http.MethodGet,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	var SWKey []byte
	SWKey = append(SWKey, CreatMuskKey()...)
	SWKey = append(SWKey, CreatMuskKey()...)
	SWKey = append(SWKey, CreatMuskKey()...)
	SWKey = append(SWKey, CreatMuskKey()...)

	req.Header["Upgrade"] = []string{"websocket"}
	req.Header["Connection"] = []string{"Upgrade"}
	req.Header["Sec-WebSocket-Key"] = []string{string(SWKey)}
	req.Header["Sec-WebSocket-Version"] = []string{"13"}
	c := &http.Client{}
	res, err := c.Do(req)
	log.Println(res)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		fmt.Println("net.Dail err", err)
		return nil, err
	}

	MyConn := &MyConn{
		Conn:            conn,
		ReadBufferSize:  DefaultReadBuffer,
		WriteBufferSize: DefaultWriteBuffer,
		PingTimeOut: func() time.Time {
			return time.Now().Add(DefaultPingWait)
		},
	}
	return MyConn, nil
}
