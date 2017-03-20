package utils

import (
	"runtime"
	"log"
	"time"
	"os"
	"os/exec"
)

func Browser(url string)  {
	// 自动打开web浏览器
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	case "darwin":
		cmd = exec.Command("open", url)
	}
	if cmd != nil {
		go func() {
			log.Println("Open the default browser after two seconds...")
			time.Sleep(time.Second * 2)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run()
		}()
	}
}