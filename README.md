# 🎉MyWebsocket

![pass](https://img.shields.io/badge/building-pass-green)![pass](https://img.shields.io/badge/checks-pass-green)

#### ✨**已经实现：**

- [x]升级协议
```go
  func (u *Upgrader) Upgrade(w http.ResponseWriter, r *http.Request) (conn *MyConn, err error)
  //通过填写upgrader 升级HTTP连接为websocket
  ```

- [x]读取消息
```go
  func (conn *MyConn)ReadMsg(opts ...Option)(messagetype int, p []byte, err error)
  //从连接中读取消息 返回数据类型、大小和错误
  ```

- [x]写入JSON、String、Binary
```go
  func (conn *MyConn) WriteJSON(v interface{}, opts ...Option) error
  func (conn *MyConn) WriteString(s string, opts ...Option) error
  func (conn *MyConn) WriteBinary(msg []byte, opts ...Option)error
  //将数据写入连接
  ```

- [x]关闭连接
```go
  func (conn *MyConn) Close()
  //关闭连接
  ```

#### 正在实现：

- [ ] 心跳
- [ ] 文件传输
- [ ] 扩展协议

### 🎁特点

支持**分片传输，自动扩容**

用户可自定义**读写缓冲与读写超时**

