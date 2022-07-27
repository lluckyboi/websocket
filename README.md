# 🎉MyWebsocket

![pass](https://img.shields.io/badge/building-pass-green) ![pass](https://img.shields.io/badge/checks-pass-green) ![pass](https://img.shields.io/badge/tests-pass-green)

## 🎁特性

- [x] 支持**自动分片传输，自动扩容**


- [x] **简洁强大而灵活**，用户可自定义每个分片负载数据的大小


- [x] 支持多种格式，**文件传输**，无需太过关心大小限制


- [x] 用户可自定义**读写缓冲与读写超时**


- [x] 一键式**心跳管理**，用户无需关心如何实现心跳


- [x] 支持读写数据帧**追踪** 轻松debug


## 🎿快速开始

```go
var up = my_websocket.Upgrader{
HandshakeTimeout: time.Second * 5,
ReadBufferSize:   2048,
WriteBufferSize:  2048,
}

func main() {
r := gin.Default()
r.GET("/ws", ping)
r.Run(":9924")
}

func ping(c *gin.Context) {
//升级get请求为webSocket协议
ws, _ := up.Upgrade(c.Writer, c.Request)
defer ws.Close()
for {
//读取ws中的数据
_, ms, _ := ws.ReadMsg()
log.Println("received:", string(ms))
//写入string到ws连接
err=ws.WriteString("hello my websocket")
}
}

```


## ✨**已经实现：**

- [x] 升级协议

```go
 //通过填写upgrader 升级HTTP连接为websocket
 func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request,opts ...Option) (conn *MyConn, err error)

 //其中opts支持以下方法
 func WithPingWait(timeout time.Duration) Option    //心跳时间
 func WithPongHandler(handler PongHandler)Option    //自定义pongHandler
 func WithIOLOG(need bool) Option                   //读写数据帧追踪
```

- [x] 读取消息

```go
 //从连接中读取消息 返回数据类型、大小和错误
 func (conn *MyConn)ReadMsg()(messagetype int, p []byte, err error)
  
 //可通过以下方法设置读取缓冲大小
 func (conn *MyConn)SetWriteBuffersize(size int64)
```

- [x] 写入JSON、String

```go
 //将数据写入连接
 func (conn *MyConn) WriteJSON(v interface{}, opts ...Option) error
 func (conn *MyConn) WriteString(s string, opts ...Option) error

 //可通以下方法设置写入缓冲大小
 func (conn *MyConn)SetReadBuffersize(size int64)
```

- [x] 关闭连接

```go
 //关闭连接
 func (conn *MyConn) Close()
```

- [x] 心跳


​    Upgrade方法通过**可选参数**自定义心跳超时时间 (默认30秒)

​    用户还可用通过WithPongHandler方法自定义服务端PongHandler

- [x] 文件传输(需要客户端解析)

```go
 //占用binary位
func (conn *MyConn) WriteFile(filePath string, fileName string, opts ...Option) error
```

![uTools_1658799974993](http://typora.fengxiangrui.top/1658799978.png)

- [x] 分片传输


![uTools_1658734731483](http://typora.fengxiangrui.top/1658734761.png)



- [x] 读写数据帧追踪

```go
//除了在Upgrade时切换读写数据帧追踪，也可以调用以下方法随时切换
func (conn *MyConn)SetIOLog(need bool)
```

![uTools_1658807634153](http://typora.fengxiangrui.top/1658807652.png)

## 🛠正在实现：

- [ ] 客户端封装


## 🧪实现原理

根据websocket协议，读取数据帧并通过http/TCP进行通信

#### 升级协议的实现

首先根据协议 检查请求头等配置，然后从http.ResponseWriter重新拿到conn ，接下拿到浏览器生成的密钥

并与Websocket的 **Magic String(258EAFA5-E914-47DA-95CA-C5AB0DC85B11)** 拼接后进行**sha1加密+base64编码**

最后回复报文

#### 读取消息的实现

依旧是基于TCP连接，从连接中读取数据帧，按照协议进行处理，如果**数据大于缓冲，则自动扩容**

拿到数据后按照协议解码处理，再根据opcode找到相应handler

#### 写入消息的实现

先将入参转换为```[]byte```,再根据数据长度与缓冲大小决定是否**分片传输**

#### 心跳的实现

通过操作net包的```SetDeadline```

#### 文件传输

有三种实现方法

**第一种**是使用多个websocket协议的非控制保留位，与客户端约定文件类型与opcde的映射关系

**第二种**是使用一个非控制保留位，第一个字节前四位设置为 `0 0 0 0 ` 	后四位为`m` (m可为0x3-7)

表示本次传输无扩展协议，使用分片传输

类型是文件信息（文件类型、大小等等），然后后续附加数据帧传输文件，传输完成后由客户端根据第一帧的文件信息拼接、解析文件

**第三种**是占用binary，约定当数据帧为二进制数据时，Payload Data中的前20个字节为文件名(xxx.png等等)

目前采用第三种(不会被普通客户端识别为异常数据帧)

## 📑Reference

[后端2021红岩课件-websocket]:

https://www.yuque.com/gyxffu/uv3zph/gpib7h#Websocket

[网络图片]：

![img](https://img-blog.csdn.net/20140306233501843?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMDQ4NzU2OA==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)