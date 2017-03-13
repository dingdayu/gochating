package handlers

import (
	"github.com/fsnotify/fsnotify"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"runtime"
)

const TEMPLATE_DIR = "./templates"

var templates map[string]*template.Template
var reloadTemplatesLock *sync.Mutex
var chs chan string

func init() {
	reloadTemplatesLock = &sync.Mutex{}
	chs = make(chan string, 10)
	templates = make(map[string]*template.Template)

	reloadTemplates()
	go watchTemplates(chs)

	go func() {
		for {
			paths := <- chs
			if ext := path.Ext(paths); ext == ".html" {
				log.Println(paths)
				retloadTemplate(paths)
			}
		}
	}()
}

// 实现文件改动的监听
func watchTemplates(chs chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//log.Println("modified file:", event.Name, event.Op)
				// 当文件夹下有写入，和创建动作时就重新加载目标文件 （复写map的元素）
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op == fsnotify.Create {
					//log.Println("modified file:", event.Name, event.Op)
					paths := "./" + strings.Replace(event.Name, "\\", "/", -1)
					//retloadTemplate(path, lock)
					//reloadTemplates()
					chs <- paths
					runtime.Gosched()

				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(TEMPLATE_DIR)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}

// 将目录文件夹下的目标文件加载到内存中
func reloadTemplates() {
	// 获取目录下所有文件及文件夹
	fileInfoArr, err := ioutil.ReadDir(TEMPLATE_DIR)
	if err != nil {
		panic(err)
	}

	var templateName, templatePath string
	// 循环所有文件，这里目前进暂不支持子目录
	for _, fileInfo := range fileInfoArr {
		templateName = fileInfo.Name()
		templatePath = TEMPLATE_DIR + "/" + templateName
		retloadTemplate(templatePath)
	}
}

func retloadTemplate(templatePath string) {
	// 非html后缀不加载
	reloadTemplatesLock.Lock()
	log.Println("Loading template:", templatePath)
	f, err := os.Open(templatePath)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	f.Close()
	t := template.Must(template.New("test").Parse(string(s)))
	//t := template.Must(template.ParseFiles(templatePath))
	templates[templatePath] = t
	reloadTemplatesLock.Unlock()
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
	chs <- "./templates/index.html"
	LoadHtml(w, TEMPLATE_DIR+"/index.html", nil)
}
