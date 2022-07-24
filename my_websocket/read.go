package my_websocket

import (
	"errors"
	"log"
	"strconv"
)

func (conn *MyConn) ReadMsg() (messagetype int, p []byte, err error) {
	err = conn.conn.SetDeadline(conn.PingTimeOut)
	if err != nil {
		return 0, nil, err
	}

	//按字节读
	msg := make([]byte, conn.ReadBufferSize)
	n, err := conn.conn.Read(msg)
	log.Printf("read  p :%b", msg)

	if err != nil {
		return -1, nil, errors.New("read data err:" + err.Error())
	}
	//如果消息大小大于ReadBufferSize 一次读不完
	if n > conn.ReadBufferSize {
		return -1, nil, errors.New("ReadBufferSize is small than data size " + strconv.Itoa(conn.ReadBufferSize) + " < " + strconv.Itoa(n))
	}

	//数据帧读取
	//每一位 右移加逻辑与 获得
	bit := make([]byte, 9)
	for i := 1; i <= 8; i++ {
		bit[i] = (msg[0] >> (8 - i)) & 1
	}

	//掩码标识
	mask := msg[1] >> 7
	if mask != 1 {
		return -1, nil, errors.New("no mask")
	}
	//数据长度
	playloadLength := msg[1] - 128
	var trueLength int64

	//找出masking-key起始字节索引
	maskst := 0
	//小于等于125 真实长度
	//126后面两个字节补充长度
	//127后面八个字节补充长度
	if playloadLength <= 125 {
		maskst = 2
		trueLength = int64(playloadLength)
	} else if playloadLength == 126 {
		maskst = 4
		//后面补充16位
		trueLength = int64(playloadLength) + int64(msg[2])<<8 + int64(msg[3])
	} else if playloadLength == 127 {
		maskst = 10
		//后面补充64位
		trueLength = int64(playloadLength) + int64(msg[2])<<56 + int64(msg[3])<<48
		trueLength += int64(msg[4])<<40 + int64(msg[5])<<32 + int64(msg[6])<<24<<int64(msg[7])<<16
		trueLength += int64(msg[8])<<8 + int64(msg[9])
	}

	//数据
	data := make([]byte, trueLength)

	//获取掩码key
	maskKey := make([]byte, 4)
	for i := 0; i <= 3; i++ {
		maskKey[i] = msg[maskst+i]
	}

	//数据开始位置
	datast := int64(maskst) + 4

	//检查是否消息分片
	//如果分片 接受完整数据 最多2048片
	msbuffer := make([]byte, conn.ReadBufferSize)
	if bit[1] == 0 {
		log.Println("将分片读取")
		i := 0
		for ; i < 2048; i++ {
			//读数据 错误返回 否则追加数据
			nx, err := conn.conn.Read(msbuffer)
			if err != nil {
				return -1, nil, errors.New("read data err:" + err.Error())
			}
			//长度检查 超长扩容
			if n+nx >= conn.ReadBufferSize {
				m := make([]byte, conn.ReadBufferSize)
				msg = append(msg, m...)
			}

			//追加
			copy(msg[n:n+nx-1], msbuffer[:nx-1])
			//下一次追加从nx开始
			n += nx
			//如果是消息最后一个分片 结束读取
			if msg[n]>>7 == 1 {
				break
			}
		}
		//如果分片过多
		if i == 2048 {
			return -1, nil, errors.New("too many data slice: over " + strconv.Itoa(conn.ReadBufferSize))
		}
	}

	//掩码处理 从掩码后数据第一字节开始
	for i := int64(0); i < trueLength; i++ {
		j := i % 4
		data[i] = msg[datast+i] ^ maskKey[j]
	}

	//计算出opcode
	opcode := int(bit[5]<<3 + bit[6]<<2 + bit[7]<<1 + bit[8])

	p = make([]byte, trueLength)
	if opcode == PingMessage {
		conn.pingHandler()
	} else if opcode == PongMessage {

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
	} else {
		//读取到其他消息 返回错误
		return -1, data, errors.New("unknown data type:" + strconv.Itoa(opcode))
	}
	return
}
