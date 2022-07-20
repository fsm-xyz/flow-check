package info

import (
	"fmt"
	"strconv"
)

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

func (st *Status) checkFlow() bool {
	balance, _ := strconv.ParseFloat(st.ResBody.FlowInfor.BalanceFeeTotal, 64)
	total, _ := strconv.ParseFloat(st.ResBody.FlowInfor.HighFeeTotal, 64)
	if balance < C.Quota.MinFlowBalance || balance < total*C.Quota.MaxFlowBalancePercent {
		return false
	}

	return true
}

func (st *Status) checkVoice() bool {
	balance, _ := strconv.ParseFloat(st.ResBody.VoiceInfor.BalanceFeeTotal, 64)
	total, _ := strconv.ParseFloat(st.ResBody.VoiceInfor.HighFeeTTotal, 64)
	if balance < C.Quota.MinVoiceBalance || balance < total*C.Quota.MaxVoiceBalancePercent {
		return false
	}

	return true
}

type mailType = int

const (
	FlowMail  mailType = 1
	VocieMail mailType = 2
)

func (st *Status) sendMail(t mailType) {
	var data string
	if t == FlowMail {
		data = fmt.Sprintf("流量不足，总共: %sG, 剩余: %sG", st.ResBody.FlowInfor.HighFeeTotal, st.ResBody.FlowInfor.BalanceFeeTotal)
	} else if t == VocieMail {
		data = fmt.Sprintf("语音不足，总共: %s分钟, 剩余: %s分钟", st.ResBody.VoiceInfor.HighFeeTTotal, st.ResBody.VoiceInfor.BalanceFeeTotal)
	}

	fmt.Println(data)
	m.Send(data)
}
