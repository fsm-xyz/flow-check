package info

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	getStatusUrl  = "https://sn.clientaccess.10086.cn/html5/indivbusi/familyNew/getStatus"
	getStatusBody = `{
		"cid": "0",
		"sn": "0",
		"cv": "0",
		"sv": "0",
		"os": "0",
		"token": "%s",
		"phoneNum": "%s",
		"xk": "0",
		"reqBody": {
		  "moduleType": 202201001
		}
	  }`
)

var token string

func getStatus() {
	getStatusReq := buildReq(getStatusUrl, fmt.Sprintf(getStatusBody, token, C.YD.Tel))
	data := httpDo(getStatusReq)
	st := &Status{}
	fmt.Println("信息结果: ", string(data))
	if err := json.Unmarshal(data, st); err != nil {
		log.Println("解析结果出错", "err: ", err)
		return
	}

	checkAndMail(st.ResBody.VoiceInfor.BalanceFeeTotal, st.ResBody.VoiceInfor.HighFeeTTotal, st.ResBody.FlowInfor.BalanceFeeTotal, st.ResBody.FlowInfor.HighFeeTotal)
}

func buildReq(url, body string) *http.Request {
	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		log.Println("请求错误", "err: ", err)
	}
	req.Close = false

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	// req.Header.Set("Cookie", "")
	req.Header.Set("Host", "sn.clientaccess.10086.cn")
	req.Header.Set("Origin", "https://sn.clientaccess.10086.cn")
	req.Header.Set("Referer", "https://sn.clientaccess.10086.cn/html5/sx/vfamilyN/index.html")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	return req
}

func httpDo(req *http.Request) []byte {
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Println("http请求出错", "err: ", err)
		return nil
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("读取结果出错", "err: ", err)
		return nil
	}

	return respBody
}

type Status struct {
	ResBody struct {
		VName     string  `json:"vName"`
		Blance    float64 `json:"blance"`
		Upflag    string  `json:"upflag"`
		Llcx      string  `json:"llcx"`
		FlowInfor struct {
			BalanceFeeTotal string `json:"balanceFeeTotal"`
			HighFeeTotal    string `json:"highFeeTotal"`
			List            []struct {
				Use      string `json:"use"`
				PerUse   string `json:"perUse"`
				PhoneNum string `json:"phoneNum"`
				ShortNum string `json:"shortNum"`
			} `json:"list"`
			DiscntFeeTotal string `json:"discntFeeTotal"`
		} `json:"flowInfor"`
		BeforFalg  bool `json:"beforFalg"`
		VoiceInfor struct {
			BalanceFeeTotal  string `json:"balanceFeeTotal"`
			DiscntValueTotal int    `json:"discntValueTotal"`
			HighFeeTTotal    string `json:"highFeeTTotal"`
			List             []struct {
				Use      string `json:"use"`
				PerUse   string `json:"perUse"`
				PhoneNum string `json:"phoneNum"`
				ShortNum string `json:"shortNum"`
			} `json:"list"`
		} `json:"voiceInfor"`
		UserList []struct {
			IsMain   bool   `json:"isMain"`
			PhoneNum string `json:"phoneNum"`
			Remark   string `json:"remark"`
			Pic      string `json:"pic"`
			ShortNum string `json:"shortNum"`
		} `json:"userList"`
		Sts      string `json:"sts"`
		IopSts   string `json:"iopSts"`
		VDesc    string `json:"vDesc"`
		EndFalg  bool   `json:"endFalg"`
		IsHandle bool   `json:"isHandle"`
	} `json:"resBody"`
}
