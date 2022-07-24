# 🎉MyWebsocket

![pass](https://img.shields.io/badge/checks-pass-green) ![pass](https://img.shields.io/badge/checks-pass-green)
#### ✨**已经实现：**

- [x] 升级协议
```go
  func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request,opts ...Option) (conn *MyConn, err error)
  //通过填写upgrader 升级HTTP连接为websocket
  ```

- [x] 读取消息
```go
  func (conn *MyConn)ReadMsg()(messagetype int, p []byte, err error)
  //从连接中读取消息 返回数据类型、大小和错误
  ```

- [x] 写入JSON、String、Binary
```go
  func (conn *MyConn) WriteJSON(v interface{}, opts ...Option) error
  func (conn *MyConn) WriteString(s string, opts ...Option) error
  func (conn *MyConn) WriteBinary(msg []byte, opts ...Option)error
  //将数据写入连接
  ```
- [x] 关闭连接
```go
  func (conn *MyConn) Close()
  //关闭连接
  ```
- [x] 心跳

Upgrade方法通过**可选参数**自定义心跳超时时间 (默认十秒)

用户还可用通过WithPongHandler方法自定义服务端PongHandler

#### 正在实现：


- [ ] 文件传输
- [ ] 扩展协议

### 🎁特点

支持**自动分片传输，自动扩容**

支持多种格式写入，无需太过关心大小限制

用户可自定义**读写缓冲与读写超时**

**一键式**心跳管理，用户无需关心如何实现心跳

