package info

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var C = &Config{}

type Config struct {
	Mail struct {
		To       string `json:"mail_to"`
		From     string `json:"mail_from"`
		Subject  string `json:"subject"`
		Password string `json:"password"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
	} `json:"mail"`
	Quota struct {
		MaxFlowBalancePercent  float64 `json:"max_flow_balance_percent"`
		MaxVoiceBalancePercent float64 `json:"max_voice_balance_percent"`
		MinFlowBalance         float64 `json:"min_flow_balance"`
		MinVoiceBalance        float64 `json:"min_voice_balance"`
	} `json:"quota"`
	YD struct {
		Token string `json:"token"`
		Tel   string `json:"tel"`
	}
}

func BuildConfig(filename string) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Panicf("读取配置文件失败, err: %s\n", err)
	}

	if err = json.Unmarshal(data, C); err != nil {
		log.Panicf("解析配置文件失败, err: %s\n", err)
	}
}
