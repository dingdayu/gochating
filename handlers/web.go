package handlers

import (
	"html/template"
	"net/http"
	"os"
	"io/ioutil"
	"path"
	"log"
)

const TEMPLATE_DIR  = "./templates"

var templates map[string]*template.Template

func init()  {
	// 获取目录下所有文件及文件夹
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
	}

	var templateName, templatePath string
	// 循环所有文件，这里目前进暂不支持子目录
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		// 非html后缀不加载
		if ext := path.Ext(templateName); ext != ".html" {
			continue
		}
		templatePath = TEMPLATE_DIR + "/" + templateName
		log.Println("Loading template:", templatePath)
		t := template.Must(template.ParseFiles(templatePath))
		templates = make(map[string] *template.Template)
		templates[templatePath] = t
	}
}

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
	err := templates[path].Execute(w, locals)
	if err != nil {
		// 抛出错误，并向上层层抛出错误
		panic(err)
	}
}

func Hello(w http.ResponseWriter, r *http.Request) {
	// 将字符串通过回写指针返回给浏览器
	locals := make(map[string]interface{})
	locals["name"] = "dingdayu"
	LoadHtml(w,  TEMPLATE_DIR + "/index.html", nil)
}
