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

type Monster struct {
	//Name string `json:"monster_name"`
	//Age  string `json:"monster_age"`
	Nm string `json:"nm"`
}

var monster = Monster{
	Nm: "sl",
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
	defer ws.Close()

	for {
		//读取ws中的数据
		_, ms, err := ws.ReadMsg()
		if err != nil {
			log.Println(err)
			break
		}
		log.Println("received:", string(ms))
		err = ws.WriteString("hello my websocket")
		if err != nil {
			log.Println(err)
			break
		}
		if string(ms) == "close" {
			break
		}
	}
}
