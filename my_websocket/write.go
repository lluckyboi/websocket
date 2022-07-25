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

func (conn *MyConn) Write(msg []byte, opcode int) error {
	//初始化
	ts := &Writer{
		idx:      0,
		datast:   0,
		maskKey:  make([]byte, 4),
		restDate: len(msg),
		ismain:   true,
	}
	p := make([]byte, conn.WriteBufferSize)
	//剩余数据大于缓冲 扩容 协议头最多占用14字节
	if ts.restDate >= conn.WriteBufferSize-14 {
		m := make([]byte, DefaultWriteBuffer)
		p = append(p, m...)
	}
	for {
		//如果还是不够 分片传输
		if ts.restDate >= conn.WriteBufferSize+DefaultWriteBuffer-14 {
			fmt.Println("剩余数据大于缓冲 分片传输中，剩余大小:", ts.restDate, "B")
			//  0 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text 如果末尾0表示当前为分片
			//  0 0 0 0 0 0 0 2 表示无扩展协议 传输分片 类型binary 如果末尾0表示当前为分片
			if ts.ismain {
				p[0] = byte(opcode)
			} else {
				p[0] = 0
			}

			//长度扩展两个字节
			p[1] = 126
			ts.datast = 4

			var i = 0
			for ; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
			}
			//设置长度
			p[3] = byte(i % 128)
			p[2] = byte(i - (i % 128))

			//记录传输位置
			ts.idx += i
			//更新剩余数据大小
			ts.restDate -= i
			//再次判断 以防边界条件 能一次传完 不分片
			if ts.restDate <= 0 {
				p[0] = 1<<7 + byte(opcode)
				_, err := conn.conn.Write(p)
				if err != nil {
					return err
				}
				break
			}
			//写入
			ts.ismain = false
			_, err := conn.conn.Write(p[:i])
			if err != nil {
				return err
			}
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
				i := ts.restDate - ts.datast
				p[3] = byte(i % 128)
				p[2] = byte(i - (i % 128))
			}
			i := 0
			for i = 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}

			log.Printf("write p :%b", p[:i])
			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p[:i])
			if err != nil {
				return err
			}
			if opcode == 1 {
				log.Println("send:", string(p[ts.datast:len(msg)+ts.datast]))
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

	return conn.Write(msg, 1)
}

func (conn *MyConn) WriteString(s string, opts ...Option) error {
	//可选参数 设置读写时间
	op := ConnOptions{
		WriteTimeOut: time.Second,
	}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.WriteTimeOut = op.WriteTimeOut
	msg := []byte(s)
	return conn.Write(msg, 1)
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

	return conn.Write(msg, 2)
}
