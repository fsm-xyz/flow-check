package info

import (
	"encoding/json"
	"fmt"
	"headless-go/mail"
	"io/ioutil"
	"log"
)

var C = &Config{}

type Config struct {
	Mail struct {
		To       string `json:"to"`
		From     string `json:"from"`
		Subject  string `json:"subject"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"mail"`
	Quota struct {
		MinFlowBalancePercent  float64 `json:"min_flow_balance_percent"`
		MinVoiceBalancePercent float64 `json:"min_voice_balance_percent"`
		MinFlowBalance         float64 `json:"min_flow_balance"`
		MinVoiceBalance        float64 `json:"min_voice_balance"`
	} `json:"quota"`
	YD struct {
		Tel         string `json:"tel"`
		WaitSMSTime int    `json:"wait_sms_time"`
		ReloadTime  int    `json:"reload_time"`
	}
	B struct {
		Headless bool `json:"headless"`
	} `json:"b"`
	Port int  `json:"port"`
	Api  bool `json:"api"`
}

func BuildConfig(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("读取配置文件失败, err: %s\n", err)
	}

	if err = json.Unmarshal(data, C); err != nil {
		log.Panicf("解析配置文件失败, err: %s\n", err)
	}

	M = &mail.Mail{
		To:      C.Mail.To,
		From:    C.Mail.From,
		Subject: C.Mail.Subject,
		Host:    C.Mail.Host,
		Port:    C.Mail.Port,

		Password: C.Mail.Password,
	}
}

func RefreshToken(t string) {

	fmt.Println("token is ", token)
	token = t
}
