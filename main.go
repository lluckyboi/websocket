package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mywebsocket/my_websocket"
	"time"
)

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
	ws, err := up.Upgrade(c.Writer, c.Request)
	if err != nil {
		log.Println("up" + err)
		return
	}
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMsg()
		if err != nil {
			break
		}
		log.Println(mt, message)
		////写入ws数据
		//err = ws.WriteJSON(gin.H{"json":"json"})
		//if err != nil {
		//	break
		//}
	}
}
