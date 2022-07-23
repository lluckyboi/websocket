package my_websocket

import (
	"encoding/json"
	"log"
)

func (conn *MyConn) WriteJSON(v interface{}) error {
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

			//获取掩码
			//copy(ts.maskKey, CreatMuskKey())
			//写入MuskKey
			//copy(p[ts.maskst:ts.maskst+3], ts.maskKey[:])
			//赋值并掩码处理

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
				conn.conn.Write(p)
				break
			}
			//写入
			ts.ismain = false
			conn.conn.Write(p)
		} else {
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			// 1 0 0 0 0 0 0 0 表示当前为最后一片 类型扩展数据
			if ts.ismain {
				p[0] = 129
			} else {
				p[0] = 128
			}

			//处理Payload len
			if ts.restDate <= 125 {
				p[1] = byte(ts.restDate)
				ts.datast = 2
			} else {
				p[1] = 126
				ts.datast = 4
			}

			//todo 服务器发给客户端不应该掩码！！！！
			//获取掩码
			//copy(ts.maskKey, CreatMuskKey())
			//写入MuskKey
			//copy(p[ts.maskst:ts.maskst+3], ts.maskKey[:])
			//赋值并掩码处理

			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				p[ts.datast+i] = msg[ts.idx+i]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}
			log.Printf("write p :%b", p)
			_, err = conn.conn.Write(p[:])
			if err != nil {
				log.Println("write into tcp err:", err)
			}
			break
		}
	}
	return nil
}
