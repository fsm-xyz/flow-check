package info

import (
	"fmt"
	"strconv"
	"strings"
)

func checkFlow(fb, ft string) bool {
	balance, _ := strconv.ParseFloat(fb, 64)
	total, _ := strconv.ParseFloat(ft, 64)
	if balance < C.Quota.MinFlowBalance || balance < total*C.Quota.MinFlowBalancePercent {
		return false
	}

	return true
}

func checkVoice(vb, vt string) bool {

	balance, _ := strconv.ParseFloat(strings.Trim(vb, "分钟"), 64)
	total, _ := strconv.ParseFloat(strings.Trim(vt, "分钟"), 64)
	if balance < C.Quota.MinVoiceBalance || balance < total*C.Quota.MinVoiceBalancePercent {
		return false
	}
	return true
}

func checkAndMail(vb, vt, fb, ft string) {
	var data1, data2 string
	if !checkFlow(fb, ft) {
		data1 = fmt.Sprintf("流量不足，总共: %sG, 剩余: %sG", ft, fb)
	}
	if !checkVoice(vb, vt) {
		data2 = fmt.Sprintf("语音不足，总共: %s, 剩余: %s", vt, vb)
	}

	if len(data1) > 0 || len(data2) > 0 {
		data := fmt.Sprintf("%s\n%s", data1, data2)
		fmt.Println(data)
		M.Send(data)
	}
}
