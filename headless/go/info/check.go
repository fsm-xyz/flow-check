package info

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	var data string
	now := time.Now()

	if now.Hour() == 6 && (now.Minute() == 30 || now.Minute() == 0) {
		data = fmt.Sprintf("每日定时检测\n流量总共: %sG, 剩余: %sG\n语音总共: %s, 剩余: %s\n", ft, fb, vt, vb)
	}

	if !checkFlow(fb, ft) {
		data = fmt.Sprintf("%s流量不足, 总共: %sG, 剩余: %sG\n", data, ft, fb)
	}
	if !checkVoice(vb, vt) {
		data = fmt.Sprintf("%s语音不足, 总共: %s, 剩余: %s\n", data, vt, vb)
	}

	if len(data) > 0 {
		fmt.Println(data)
		M.Send(data)
	}
}
