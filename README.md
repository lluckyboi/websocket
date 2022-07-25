# 🎉MyWebsocket

![pass](https://img.shields.io/badge/building-pass-green) ![pass](https://img.shields.io/badge/checks-pass-green)

## 🎁特点

支持**自动分片传输，自动扩容**

支持多种格式，文件传输，无需太过关心大小限制

用户可自定义**读写缓冲与读写超时**

**一键式**心跳管理，用户无需关心如何实现心跳

## ✨**已经实现：**

- [x] 升级协议

```go
 //通过填写upgrader 升级HTTP连接为websocket
 func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request,opts ...Option) (conn *MyConn, err error)

 //其中opts支持以下方法
 func WithPingWait(timeout time.Duration) Option
 func WithPongHandler(handler PongHandler)Option
```

- [x] 读取消息

```go
 //从连接中读取消息 返回数据类型、大小和错误
 func (conn *MyConn)ReadMsg()(messagetype int, p []byte, err error)
  
 //可通过以下方法设置读取缓冲大小
 func (conn *MyConn)SetWriteBuffersize(size int
```

- [x] 写入JSON、String

```go
 //将数据写入连接
 func (conn *MyConn) WriteJSON(v interface{}, opts ...Option) error
 func (conn *MyConn) WriteString(s string, opts ...Option) error

 //可通以下方法设置写入缓冲大小
 func (conn *MyConn)SetReadBuffersize(size int)
```

- [x] 关闭连接

```go
 //关闭连接
 func (conn *MyConn) Close()
```

- [x] 心跳


    Upgrade方法通过**可选参数**自定义心跳超时时间 (默认30秒)

    用户还可用通过WithPongHandler方法自定义服务端PongHandler

- [x] 文件传输(需要客户端设置自定义解析)

```go
 // 通过binary格式传输，目前只约定了jpg格式 
 func (conn *MyConn) WriteImageJPG(filePath string, opts ...Option) error
```
  
    
## 🛠正在实现：

- [ ] 分布式websocket

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



### 📑Reference

[后端2021红岩课件-websocket]https://www.yuque.com/gyxffu/uv3zph/gpib7h#Websocket

网络图片：

![img](https://img-blog.csdn.net/20140306233501843?watermark/2/text/aHR0cDovL2Jsb2cuY3Nkbi5uZXQvdTAxMDQ4NzU2OA==/font/5a6L5L2T/fontsize/400/fill/I0JBQkFCMA==/dissolve/70/gravity/SouthEast)