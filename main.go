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

type js struct {
	jss string `json:"jss"`
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
		log.Println("up" + err.Error())
		return
	}
	js := js{jss: "123"}
	for {
		//写入ws数据
		err = ws.WriteJSON(js)
		if err != nil {
			log.Println(err)
			break
		}
		time.Sleep(time.Second)
	}
}
