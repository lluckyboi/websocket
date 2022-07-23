package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	msg := make([]byte, 4)
	m := []byte(CreatMuskKey())
	fmt.Println(m)
	msg = append(msg, m...)
	fmt.Println(msg)
}

func CreatMuskKey() string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 4)
	for i := range b {
		b[i] = letterRunes[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(letterRunes))]
		//休眠一纳秒
		time.Sleep(time.Nanosecond)
	}
	return string(b)
}
