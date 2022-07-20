package main

import (
	"io/fs"
	"io/ioutil"
	"net/http"

	"xyz-zyz.io/fsm/flow-check/info"
)

func main() {
	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("t")
		info.RefreshToken(token)
		w.Write([]byte("token 更新成功"))
	})

	http.HandleFunc("/code", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("c")
		var data string = "code 保存成功"
		if err := ioutil.WriteFile("code.txt", []byte(token), fs.ModePerm); err != nil {
			data = "code 保存失败"
		}
		info.RefreshToken(token)
		w.Write([]byte(data))
	})

	go func() {
		info.Run()
	}()

	http.ListenAndServe(":8000", nil)
}
