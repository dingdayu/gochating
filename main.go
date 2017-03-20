package main

import (
	"github.com/dingdayu/gochatting/handlers"
	"net/http"
	"log"
	"runtime/debug"
	"golang.org/x/net/websocket"
	"github.com/dingdayu/gochatting/handlers/api"
	"github.com/dingdayu/gochatting/handlers/web"
)

func main() {
	// 注册一个路由
	http.Handle("/websocket", websocket.Handler(handlers.Connection))

	http.HandleFunc("/", safeWebHandler(web.Hello))
	http.HandleFunc("/hello", safeWebHandler(web.Hello))
	http.HandleFunc("/login", safeWebHandler(web.Login))

	http.HandleFunc("/api/getOnlineUserList", api.GetOnlineUserList)
	http.HandleFunc("/api/json", api.HelloJson)
	http.HandleFunc("/api/login", api.Login)

	http.HandleFunc("/public/", web.PublicHandler)


	//models.GetUser()
	//utils.Browser("http://127.0.0.1:8080");
	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}

}

// 服务器内部错误拦截
func safeWebHandler(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 遇到错误时的扫尾工作
		defer func() {
			// 终止（拦截）错误的传递
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				//w.WriteHeader(http.StatusInternalServerError)
				//handlers.LoadHtml(w, "./templates/50x.html", nil)
				log.Println("WARN: panic in %v. - %v", fn, e)
				log.Println(string(debug.Stack()))
			}
		}()
		// 调用传入的方法名
		fn(w, r)
	}
}
