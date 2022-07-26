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
	Name string `json:"monster_name"`
	Age  string `json:"monster_age"`
}

var monster = Monster{
	Name: "web",
	Age:  "socket",
}

func main() {
	r := gin.Default()
	r.GET("/ws", ping)
	r.Run(":9924")
}

func ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := up.Upgrade(c.Writer, c.Request, my_websocket.WithIOLOG(true))
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
		if string(ms) == "close" {
			break
		}
		err = ws.WriteJSON(monster)
		if err != nil {
			log.Println(err)
			break
		}
		err = ws.WriteString("hello my websocket")
		if err != nil {
			log.Println(err)
			break
		}
		//多次上传图片可能会因为客户端无法正常解析而关闭连接
		//err = ws.WriteImage("./example.png")
		//if err != nil {
		//	log.Println(err)
		//	break
		//}
	}
}
