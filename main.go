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

// 定义一个结构体
type Monster struct {
	Name string `json:"monster_name"` // 反射机制
	Age  string `json:"monster_age"`
}

// 将结构体进行序列化

var monster = Monster{
	Name: "sda",
	Age:  "50",
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

	for {
		//读取ws中的数据
		_, _, err := ws.ReadMsg()
		if err != nil {
			log.Println(err)
			break
		}
		//写入ws数据
		err = ws.WriteJSON(&monster)
		if err != nil {
			log.Println(err)
			break
		}
	}
}
