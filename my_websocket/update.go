package my_websocket

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
)


func (u *Upgrader)Upgrade(w http.ResponseWriter, r *http.Request)(conn *MyConn, err error){
	//设置默认值
	if u.ReadBufferSize==0{
		u.ReadBufferSize=DefaultReadBuffer
	}
	if u.WriteBufferSize==0{
		u.ReadBufferSize=DefaultReadBuffer
	}
	if u.HandshakeTimeout==time.Duration(0){
		u.HandshakeTimeout=DefaultTimeOut
	}

	conn.ReadBufferSize=u.ReadBufferSize
	conn.WriteBufferSize=u.WriteBufferSize

	//检查方法
	if r.Method != http.MethodGet {
		return &MyConn{},errors.New("the method is not get")
	}
	//检查请求头
	if r.Header.Get("Connection")!="Upgrade"{
		return &MyConn{},errors.New("the connection is not Upgrade")
	}
	if r.Header.Get("Upgrade")!="websocket"{
		return &MyConn{},errors.New("the connection is not websocket")
	}
	if r.Header.Get("Sec-Websocket-Version")!=WSVersion{
		return &MyConn{},errors.New("the version is not 13")
	}


	//从http.ResponseWriter重新拿到conn 出错就返回
	hijcakcer, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, http.StatusText(500),500)
		return &MyConn{},errors.New("upgrade conn err:get conn")
	}
	conn.conn, _, err = hijcakcer.Hijack()

	//拿到浏览器生成的密钥 并与Websocket的Magic String拼接
	wskey:=append([]byte(r.Header.Get("Sec-Websocket-Key")),[]byte(MagicString)...)
	respAccept:=SHA1AndBase64(string(wskey))

	//回复报文
	resp:=[]byte("HTTP/1.1 101 Switching Protocols\nUpgrade: websocket\nConnection: Upgrade\nSec-WebSocket-Accept: "+respAccept+"\n")
	conn.conn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
	conn.conn.Write(resp)
	return
}

func (conn *MyConn)ReadMsg(m Msg)(messagetype int,p []byte,err error){
	//按字节读
	msg:=make([]byte,conn.ReadBufferSize)
	n,err:=conn.conn.Read(msg)
	if err!=nil{
		return -1,nil,errors.New("read data err:"+err.Error())
	}
	//如果消息大小大于ReadBufferSize 一次读不完
	if n>conn.ReadBufferSize{
		return -1,nil,errors.New("ReadBufferSize is small than data size "+strconv.Itoa(conn.ReadBufferSize)+" < "+strconv.Itoa(n))
	}

	//数据帧读取
	//每一位 右移加逻辑与 获得
	bit:=make([]byte,9)
	for i:=1;i<=8;i++{
		bit[i]=(msg[0]>>(8-i))&1
	}

	//掩码标识
	mask:=msg[1]>>7
	if mask!=1{

	}
	//数据长度
	playloadLength:=msg[1]<<1
	var trueLength int64

	//找出masking-key起始字节索引
	maskst:=0
	//小于等于125 真实长度
	//126后面两个字节补充长度
	//127后面八个字节补充长度
	if playloadLength<=125{
		maskst=2
		trueLength=int64(playloadLength)
	}else if playloadLength==126{
		maskst=4
		//后面补充16位
		trueLength=int64(125)<<16+int64(msg[2])<<8+int64(msg[3])
	}else if playloadLength==127{
		maskst=10
		//后面补充64位
		trueLength=int64(125)<<64+int64(msg[2])<<56+int64(msg[3])<<48
		trueLength+=int64(msg[4])<<40+int64(msg[5])<<32+int64(msg[6])<<24<<int64(msg[7])<<16
		trueLength+=int64(msg[8])<<8+int64(msg[9])
	}
	//数据
 	data:=make([]byte,trueLength)

	//获取掩码key
	maskKey:=make([]byte,4)
	for i:=0;i<=3;i++{
		maskKey[i]=msg[maskst+i]
	}
	//数据开始位置
	datast:=int64(maskst)+4

	//检查是否消息分片
	//如果分片 接受完整数据 最多2048片
	msbuffer:=make([]byte,conn.ReadBufferSize)
	if bit[1]==0{
		i:=0
		for ;i<2048;i++{
			//读数据 错误返回 否则追加数据
			nx,err:=conn.conn.Read(msbuffer)
			if err!=nil{
				return -1,nil,errors.New("read data err:"+err.Error())
			}
			//长度检查 超长扩容
			if n+nx>=conn.ReadBufferSize{
				m:=make([]byte,conn.ReadBufferSize)
				msg=append(msg,m...)
			}

			//追加
			copy(msg[n:n+nx-1],msbuffer[:nx-1])
			//下一次追加从nx开始
			n+=nx
			//如果是消息最后一个分片 结束读取
			if msg[n]>>7==1{
				break
			}
		}
		//如果分片过多
		if i==2048{
			return -1,nil,errors.New("too many data slice: over "+strconv.Itoa(conn.ReadBufferSize))
		}
	}

	//掩码处理 从掩码后数据第一字节开始
	for i:=int64(0);i<trueLength;i++{
		j:=i%4
		data[i]=msg[datast+i]^maskKey[j]
	}

	//计算出opcode
	opcode:=int(bit[5]<<3+bit[6]<<2+bit[7]<<1+bit[8])

	if opcode==PingMessage{
		conn.pingHandler()
	}else if opcode==PongMessage{
		conn.pongHandler()
	}else if opcode==TextMessage{
		//收到文本消息返回
		m.Typ=opcode
		//m.content=读取到的内容
		m.content=msg[datast:datast+trueLength-1]
	}else if opcode==BinaryMessage{
		//收到二进制消息返回
		m.Typ=opcode
		//m.content=读取到的内容
		m.content=msg[datast:datast+trueLength-1]
	}else if opcode==CloseMessage {
		//收到关闭消息
		m.Typ=opcode
		return opcode,data,errors.New("received close message")
	}else {
		//读取到其他消息 返回错误
		return -1,data,errors.New("unknown data type:"+strconv.Itoa(opcode))
	}
	return
}

func (conn *MyConn) WriteJSON(v interface{}) error {
	//消息类型text
	opcode:=TextMessage

	//序列化
	msg, err :=json.Marshal(v)
	if err!=nil{
		return err
	}
	//一些用到的变量
	idx:=0
	datast:=0
	maskst:=0
	maskKey:=make([]byte,4)
	//缓冲扩容 保证大小是WriteBufferSize整数倍
	p:=make([]byte,conn.WriteBufferSize)
	for {
		//数据大于缓冲 分片传输
		if len(msg)>=conn.WriteBufferSize{
			// 1 0 0 0 表示无扩展协议 传输分片
			p[0]=8
			//大小显然大于125 8388608=2^(7+16)
			if len(msg)<=8388608{
				p[1]=254
				maskst=4
				datast=maskst+4
			}else {
				p[1]=255
				maskst=10
				datast=maskst+4
			}
			//掩码处理
			maskKey=CreatMuskKey()
			p[]
			conn.conn.Write()
			m:=make([]byte,conn.WriteBufferSize)
			p=append(p,m...)
		}else {


			conn.conn.Write()
			break
		}
	}

}

func (conn *MyConn)WriteMsg(m Msg)(err error){
	//按照数据帧写出数据
	p:=[]byte{}
	_,err =conn.conn.Write(p)
	if err != nil {
		//写不进去，咋办呢
	}
	return
}