package my_websocket

import (
	"errors"
	"log"
	"strconv"
)

func (conn *MyConn) ReadMsg() (messagetype int, p []byte, err error) {
	var opcode int
	var data []byte

	for {
		//刷新心跳
		err = conn.conn.SetDeadline(conn.PingTimeOut())
		if err != nil {
			return 0, nil, err
		}

		msg := make([]byte, conn.ReadBufferSize)
		//缓冲
		msbuff := make([]byte, DefaultWriteBuffer)
		n, err := conn.conn.Read(msbuff)
		if conn.Opts.IOLog {
			log.Printf("read  p %d Bytes:%b", n, msbuff[:n])
		}
		if err != nil {
			return -1, nil, errors.New("read data err:" + err.Error())
		}
		//如果消息大小大于ReadBufferSize 自动扩容
		if uint64(n) > conn.ReadBufferSize {
			m := make([]byte, uint64(n)-conn.ReadBufferSize)
			msg = append(msg, m...)
		}
		copy(msg, msbuff[:n])

		//数据帧读取
		//每一位 右移加逻辑与 获得
		bit := make([]byte, 9)
		for i := 1; i <= 8; i++ {
			bit[i] = (msg[0] >> (8 - i)) & 1
		}

		//计算出opcode
		opcode = int(bit[5]<<3 + bit[6]<<2 + bit[7]<<1 + bit[8])

		//掩码标识
		mask := msg[1] >> 7
		if mask != 1 {
			return -1, nil, errors.New("no mask")
		}
		//数据长度
		playloadLength := msg[1] - 128
		var thislength uint64

		//找出masking-key起始字节索引
		maskst := 0

		//检查本次数据长度
		//小于等于125 真实长度
		//126后面两个字节补充长度
		//127后面八个字节补充长度
		if playloadLength <= 125 {
			maskst = 2
			thislength = uint64(playloadLength)
		} else if playloadLength == 126 {
			maskst = 4
			//后面补充16位
			thislength = uint64(playloadLength) + uint64(msg[2])<<8 + uint64(msg[3])
		} else if playloadLength == 127 {
			maskst = 10
			//后面补充64位
			thislength = uint64(playloadLength) + uint64(msg[2])<<56 + uint64(msg[3])<<48
			thislength += uint64(msg[4])<<40 + uint64(msg[5])<<32 + uint64(msg[6])<<24<<uint64(msg[7])<<16
			thislength += uint64(msg[8])<<8 + uint64(msg[9])
		}

		//获取掩码key
		maskKey := make([]byte, 4)
		for i := 0; i <= 3; i++ {
			maskKey[i] = msg[maskst+i]
		}

		//数据开始位置
		datast := uint64(maskst) + 4

		//数据
		datas := make([]byte, thislength)

		//掩码处理 从掩码后数据第一字节开始
		for i := uint64(0); i < thislength; i++ {
			j := i % 4
			datas[i] = msg[datast+i] ^ maskKey[j]
		}
		//追加解码后数据
		data = append(data, datas...)

		//检查分片
		if bit[1] == 0 {
			//如果分片
			if opcode == 0 {
				//如果是后续扩展数据帧
				continue
			} else {
				//如果是消息第一帧
				continue
			}
		} else {
			//是消息最后一帧
			if opcode == 0 {
				//如果分片
				break
			} else {
				//如果不分片
				break
			}
		}
	}

	trueLength := len(data)
	p = make([]byte, trueLength)
	if opcode == PingMessage {
		conn.pingHandler()
	} else if opcode == PongMessage {
		conn.pongHandler()
	} else if opcode == TextMessage {
		//收到文本消息返回
		messagetype = opcode
		copy(p, data[:trueLength])
	} else if opcode == BinaryMessage {
		//收到二进制消息返回
		messagetype = opcode
		copy(p, data[:trueLength])
	} else if opcode == CloseMessage {
		//收到关闭消息
		messagetype = opcode
		copy(p, data[:trueLength])
		return opcode, data, errors.New("received close message:" + string(data))
	} else if opcode == FileImageMessage {
		messagetype = opcode
		copy(p, data[:trueLength])
	} else {
		//读取到其他消息 返回错误
		return -1, data, errors.New("unknown data type:" + strconv.Itoa(opcode))
	}
	return
}
