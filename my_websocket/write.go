package my_websocket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func (conn *MyConn) wsWrite(msg []byte, opcode int) error {
	//初始化
	ts := &Writer{
		idx:      0,
		datast:   0,
		maskKey:  make([]byte, 4),
		restDate: uint64(len(msg)),
		ismain:   true,
	}
	p := make([]byte, conn.WriteBufferSize)

	//剩余数据大于缓冲 扩容 协议头最多四个字节(不会扩容八个字节)
	if ts.restDate > conn.WriteBufferSize-4 {
		m := make([]byte, DefaultWriteBuffer)
		p = append(p, m...)
	}

	for {
		//如果大于DefaultWriteBuffer 即扩容两个字节也不够 分片传输
		if ts.restDate >= DefaultWriteBuffer {
			//  0 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text 如果末尾0表示当前为分片
			//  0 0 0 0 0 0 0 2 表示无扩展协议 传输分片 类型binary 如果末尾0表示当前为分片

			//是否主片
			if ts.ismain {
				p[0] = byte(opcode)
			} else {
				p[0] = 0
			}

			//长度扩展2个字节
			p[1] = 126
			ts.datast = 4

			var i uint64 = 0
			for ; i < conn.WriteBufferSize; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
			}

			//设置长度
			p[3] = byte(i % 128)
			p[2] = byte((i - (i % 128)) >> 8)

			//记录传输位置
			ts.idx += i
			//更新剩余数据大小
			ts.restDate -= i

			//写入
			ts.ismain = false

			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p[:i+ts.datast+1])
			if err != nil {
				return err
			}
			//是否打印数据帧
			if conn.Opts.IOLog {
				log.Printf("write p %d Bytes:%b", i, p[:i+ts.datast+1])
			}
			fmt.Println("剩余数据大于缓冲 分片传输中，剩余大小:", ts.restDate, "Byte")
		} else { //剩余数据小于缓冲 一次发完
			// 1 0 0 0 0 0 0 2 表示无扩展协议 传输不分片 类型binary
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			// 1 0 0 0 0 0 0 0 表示当前为最后一片 类型扩展数据
			//是否主片
			if ts.ismain {
				p[0] = 128 + byte(opcode)
			} else {
				p[0] = 128
			}
			//服务器发给客户端不应该掩码
			//处理Payload len
			if ts.restDate < 125 {
				p[1] = byte(ts.restDate)
				if p[1]%2 == 1 {
					p[1]++
				}
				ts.datast = 2
			} else {
				p[1] = 126
				ts.datast = 4
				i := ts.restDate - ts.datast - 1
				p[3] = byte(i % 128)
				p[2] = byte((i - (i % 128)) >> 8)
			}
			var i uint64 = 0
			for i = 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}

			//是否打印数据帧
			if conn.Opts.IOLog {
				log.Printf("write p %d Bytes:%b", i, p[:i+ts.datast+1])
			}
			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p[:i+ts.datast+1])
			if err != nil {
				return err
			}
			if opcode == 1 {
				log.Println("send:", string(p[ts.datast:uint64(len(msg))+ts.datast]))
			} else {
				log.Println("send file success")
			}
			break
		}
	}
	return nil
}

func (conn *MyConn) WriteJSON(v interface{}, opts ...Option) error {
	//可选参数 设置读写时间
	op := ConnOptions{
		WriteTimeOut: time.Second,
	}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.WriteTimeOut = op.WriteTimeOut

	//序列化
	msg, err := json.Marshal(v)
	if err != nil {
		return err
	}

	return conn.wsWrite(msg, 1)
}

func (conn *MyConn) WriteString(s string, opts ...Option) error {
	op := ConnOptions{
		WriteTimeOut: time.Second,
	}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.WriteTimeOut = op.WriteTimeOut

	msg := []byte(s)
	return conn.wsWrite(msg, 1)
}

func (conn *MyConn) WriteImageJPG(filePath string, opts ...Option) error {
	buff := new(bytes.Buffer)
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("open file err=", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println("file closed err:", err)
		}
	}(file)

	_, _ = io.Copy(buff, file)
	//log.Println(buff.Bytes())

	msg := buff.Bytes()

	//可选参数 设置读写时间
	op := ConnOptions{
		WriteTimeOut: time.Second,
	}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.WriteTimeOut = op.WriteTimeOut

	return conn.wsWrite(msg, 2)
}
