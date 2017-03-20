package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func ReturnJson(v interface{}, w http.ResponseWriter)  {
	if j, err := json.Marshal(v); err != nil {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, ``)
	} else {
		// 返回json的类型头信息
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, string(j))
	}
}