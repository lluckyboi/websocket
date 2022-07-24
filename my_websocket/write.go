package my_websocket

import (
	"encoding/json"
	"log"
	"time"
)

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
	//初始化
	ts := &Writer{
		idx:      0,
		datast:   0,
		maskKey:  make([]byte, 4),
		restDate: len(msg),
		ismain:   true,
	}

	//开始传输
	for {
		p := make([]byte, conn.WriteBufferSize)
		//剩余数据大于缓冲 分片传输 协议头最多占用14字节
		if ts.restDate >= conn.WriteBufferSize-14 {
			// 0 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text 0表示当前为分片
			if ts.ismain {
				p[0] = 1
			} else {
				p[0] = 0
			}

			//大小显然大于125 8388608=2^(7+16)
			if ts.restDate <= 8388608 {
				p[1] = 254
				ts.datast = 4
			} else {
				p[1] = 255
				ts.datast = 10
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
			}

			//记录传输位置
			ts.idx += conn.WriteBufferSize
			//更新剩余数据大小
			ts.restDate -= conn.WriteBufferSize - ts.datast
			//再次判断 以防边界条件
			if ts.restDate <= 0 {
				p[0] = 1
				_, err := conn.conn.Write(p)
				if err != nil {
					return err
				}
				break
			}
			//写入
			ts.ismain = false
			_, err := conn.conn.Write(p)
			if err != nil {
				return err
			}
		} else { //剩余数据小于缓冲 一次发完
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			// 1 0 0 0 0 0 0 0 表示当前为最后一片 类型扩展数据
			//是否主片
			if ts.ismain {
				p[0] = 129
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
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}

			log.Printf("write p :%b", p)
			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p)
			if err != nil {
				return err
			}
			log.Println("send:", string(p[ts.datast:len(msg)+ts.datast]))
			break
		}
	}
	return nil
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

	msg := make([]byte, conn.WriteBufferSize)
	msg = []byte(s)

	//初始化
	ts := &Writer{
		idx:      0,
		datast:   0,
		maskKey:  make([]byte, 4),
		restDate: len(msg),
		ismain:   true,
	}

	for {
		p := make([]byte, conn.WriteBufferSize)
		//剩余数据大于缓冲 分片传输 协议头最多占用14字节
		if ts.restDate >= conn.WriteBufferSize-14 {
			// 0 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text 0表示当前为分片
			if ts.ismain {
				p[0] = 1
			} else {
				p[0] = 0
			}

			//大小显然大于125
			if ts.restDate <= (1<<7 + 1<<16) {
				p[1] = 126
				ts.datast = 4
			} else {
				p[1] = 127
				ts.datast = 10
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
			}

			//记录传输位置
			ts.idx += conn.WriteBufferSize
			//更新剩余数据大小
			ts.restDate -= conn.WriteBufferSize - ts.datast
			//再次判断 以防边界条件
			if ts.restDate <= 0 {
				p[0] = 1
				_, err := conn.conn.Write(p)
				if err != nil {
					return err
				}
				break
			}
			//写入
			ts.ismain = false
			_, err := conn.conn.Write(p)
			if err != nil {
				return err
			}
		} else { //剩余数据小于缓冲 一次发完
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			// 1 0 0 0 0 0 0 0 表示当前为最后一片 类型扩展数据
			//是否主片
			if ts.ismain {
				p[0] = 129
			} else {
				p[0] = 128
			}

			//todo 服务器发给客户端不应该掩码
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
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}

			log.Printf("write p :%b", p)
			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p)
			if err != nil {
				return err
			}
			log.Println("send:", string(p[ts.datast:len(msg)+ts.datast]))
			break
		}
	}
	return nil
}

func (conn *MyConn) WriteBinary(s []byte, opts ...Option) error {
	//可选参数 设置读写时间
	op := ConnOptions{
		WriteTimeOut: time.Second,
	}
	for _, option := range opts {
		option(&op)
	}
	conn.Opts.WriteTimeOut = op.WriteTimeOut

	msg := make([]byte, conn.WriteBufferSize)
	msg = []byte(s)

	//初始化
	ts := &Writer{
		idx:      0,
		datast:   0,
		maskKey:  make([]byte, 4),
		restDate: len(msg),
		ismain:   true,
	}

	for {
		p := make([]byte, conn.WriteBufferSize)
		//剩余数据大于缓冲 分片传输 协议头最多占用14字节
		if ts.restDate >= conn.WriteBufferSize-14 {
			// 0 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text 0表示当前为分片
			if ts.ismain {
				p[0] = 1
			} else {
				p[0] = 0
			}

			//大小显然大于125
			if ts.restDate <= (1<<7 + 1<<16) {
				p[1] = 126
				ts.datast = 4
			} else {
				p[1] = 127
				ts.datast = 10
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
			}

			//记录传输位置
			ts.idx += conn.WriteBufferSize
			//更新剩余数据大小
			ts.restDate -= conn.WriteBufferSize - ts.datast
			//再次判断 以防边界条件
			if ts.restDate <= 0 {
				p[0] = 1
				_, err := conn.conn.Write(p)
				if err != nil {
					return err
				}
				break
			}
			//写入
			ts.ismain = false
			_, err := conn.conn.Write(p)
			if err != nil {
				return err
			}
		} else { //剩余数据小于缓冲 一次发完
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			// 1 0 0 0 0 0 0 0 表示当前为最后一片 类型扩展数据
			//是否主片
			if ts.ismain {
				p[0] = 129
			} else {
				p[0] = 128
			}

			//todo 服务器发给客户端不应该掩码
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
			}

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}

			log.Printf("write p :%b", p)
			err := conn.conn.SetWriteDeadline(time.Now().Add(conn.Opts.WriteTimeOut))
			if err != nil {
				return err
			}
			_, err = conn.conn.Write(p)
			if err != nil {
				return err
			}
			log.Println("send:", string(p[ts.datast:len(msg)+ts.datast]))
			break
		}
	}
	return nil
}
