package handlers

import (
	"html/template"
	"net/http"
	"os"
)

const TEMPLATE_DIR  = "./templates"

func PublicHandler(w http.ResponseWriter, r *http.Request) {
	// 直接调用http包提供的文件服务方法，直接根据请求路径返回文件内容
	name := r.URL.Path[len("/"):]
	if exists := isExists(name); !exists {
		// 文件找不到返回404
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, name)
}

// 判断对应文件是否存在
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}

// 模板加载
func LoadHtml(w http.ResponseWriter, path string, locals map[string]interface{}) {
	t, err := template.ParseFiles(path)
	if err != nil {
		// 抛出错误，并向上层层抛出错误
		panic(err)
	}
	t.Execute(w, locals)
}

func Hello(w http.ResponseWriter, r *http.Request) {
	// 将字符串通过回写指针返回给浏览器
	locals := make(map[string]interface{})
	locals["name"] = "dingdayu"
	LoadHtml(w,  TEMPLATE_DIR + "/index.html", nil)
}
