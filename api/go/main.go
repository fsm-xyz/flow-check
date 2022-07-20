package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"

	"xyz-zyz.io/fsm/flow-check/info"
	"xyz-zyz.io/fsm/flow-check/mail"
)

var run *bool

func main() {

	BeforeStart()

	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("t")
		RefreshToken(token)
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

	go func() {
		if *run {
			fmt.Println("开启API自动检测功能")
			info.Run()
		}
	}()

	http.ListenAndServe(":8000", nil)
}

func BeforeStart() {
	filename := flag.String("c", "conf/config.json", "配置文件")
	run = flag.Bool("r", false, "开启API检测")
	flag.Parse()

	info.BuildConfig(*filename)

	info.M = &mail.Mail{
		MailTo:   info.C.Mail.To,
		MailFrom: info.C.Mail.From,
		Subject:  info.C.Mail.Subject,
		Host:     info.C.Mail.Host,
		Port:     info.C.Mail.Port,

		Password: info.C.Mail.Password,
	}
}

func RefreshToken(token string) {
	info.C.YD.Token = token
}
