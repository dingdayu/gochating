package main

import (
	"net/http"
	"github.com/dingdayu/gochatting/handlers"
)

func main() {
	// 注册一个路由
	http.HandleFunc("/hello", handlers.Hello)
	http.HandleFunc("/api/json", handlers.HelloJson)
	http.HandleFunc("/public/", handlers.PublicHandler)

	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
