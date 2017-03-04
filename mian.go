package main

import (
	"net/http"
	"text/template"
)

func main() {
	// 注册一个路由
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		// 将字符串通过回写指针返回给浏览器
		templates := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Title</title>
</head>
<body>
hello {{.}}!
</body>
</html>
		`
		name := "dingdayu"
		t := template.Must(template.New("templates").Parse(templates))
		t.Execute(w, name)
	})
	// 监听端口 8080
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
