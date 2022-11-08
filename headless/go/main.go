package main

import (
	"flag"
	"fmt"
	"headless-go/info"
	"io/fs"
	"io/ioutil"
	"net/http"
	"time"
)

var run *bool

func main() {

	BeforeStart()

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
		w.Write([]byte(data))
	})

	http.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		go func() {
			fmt.Println("手动开启API自动检测功能")
			info.Run()
		}()
	})

	http.HandleFunc("/running", func(w http.ResponseWriter, r *http.Request) {
		data := "运行失败"
		if info.Running {
			data = "运行中"
		}
		w.Write([]byte(data))
	})

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		go info.API()
	})

	go func() {
		if *run {
			fmt.Println("开启API自动检测功能")
			info.Run()
		}
	}()

	http.ListenAndServe(fmt.Sprintf(":%d", info.C.Port), nil)
}

func BeforeStart() {
	filename := flag.String("c", "conf/config.json", "配置文件")
	run = flag.Bool("r", false, "开启检测")
	flag.Parse()

	info.BuildConfig(*filename)

	go func() {
		t := time.NewTicker(time.Second * 30)
		defer t.Stop()

		for range t.C {
			fmt.Println("定时加载配置")
			info.BuildConfig(*filename)
		}
	}()
}
