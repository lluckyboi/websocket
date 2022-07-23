package my_websocket

import "encoding/json"

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
		maskst:   0,
		maskKey:  make([]byte, 4),
		restDate: len(msg),
	}
	//开始传输
	for {
		p := make([]byte, conn.WriteBufferSize)
		//剩余数据大于缓冲 分片传输 协议头最多占用14字节
		if ts.restDate >= conn.WriteBufferSize-14 {
			// 1 0 0 0 0 0 0 1 表示无扩展协议 传输分片 类型text
			p[0] = 129
			//大小显然大于125 8388608=2^(7+16)
			if ts.restDate <= 8388608 {
				p[1] = 254
				ts.maskst = 4
			} else {
				p[1] = 255
				ts.maskst = 10
			}
			ts.datast = ts.maskst + 4

			//获取掩码
			copy(ts.maskKey, CreatMuskKey())
			//写入MuskKey
			copy(p[ts.maskst:ts.maskst+3], ts.maskKey[:])
			//赋值并掩码处理
			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				j := i % 4
				p[ts.datast+i] = msg[ts.idx+i] ^ ts.maskKey[j]
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
			conn.conn.Write(p)
		} else {
			// 0 0 0 0 0 0 0 1 表示无扩展协议 传输不分片 类型text
			p[0] = 1
			//处理Payload len
			if ts.restDate <= 125 {
				p[1] = byte(1<<7 + ts.restDate)
				ts.maskst = 2
			} else {
				p[1] = 254
				ts.maskst = 4
			}
			ts.datast = ts.maskst + 4
			//获取掩码
			copy(ts.maskKey, CreatMuskKey())
			//写入MuskKey
			copy(p[ts.maskst:ts.maskst+3], ts.maskKey[:])
			//赋值并掩码处理
			for i := 0; i < conn.WriteBufferSize-ts.datast; i++ {
				j := i % 4
				p[ts.datast+i] = msg[ts.idx+i] ^ ts.maskKey[j]
				ts.restDate--
				//如果写完 退出循环
				if ts.restDate <= 0 {
					break
				}
			}
			conn.conn.Write(p)
			break
		}
	}
	return nil
}
