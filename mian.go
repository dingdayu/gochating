package main

import (
	"net/http"
	"github.com/dingdayu/gochatting/handlers"
)

func main() {
	// 注册一个路由
	http.HandleFunc("/hello", safeWebHandler(handlers.Hello))
	http.HandleFunc("/api/json", handlers.HelloJson)
	http.HandleFunc("/public/", handlers.PublicHandler)

	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// 服务器内部错误拦截
func safeWebHandler(fn http.HandlerFunc) http.HandlerFunc  {
	return func(w http.ResponseWriter, r *http.Request) {
		// 遇到错误时的扫尾工作
		defer func() {
			// 终止（拦截）错误的传递
			if e, ok := recover().(error); ok {
				http.Error(w, e.Error(), http.StatusInternalServerError)
				// 或者输出自定义的50x错误页面
				//w.WriteHeader(http.StatusInternalServerError)
				//t,err := template.ParseFiles("./templates/50x.html")
				//if err != nil {
				//	// 抛出错误，并向上层层抛出错误
				//	panic(err)
				//}
				//t.Execute(w, nil)
			}
		}()
		// 调用传入的方法名
		fn(w, r)
	}
}