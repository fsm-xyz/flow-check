package info

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/chromedp"
)

func Loadpage(taskCtx context.Context) {
	var (
		buf    []byte
		vb, vt string
		fb, ft string
	)

	if err := chromedp.Run(
		taskCtx,
		chromedp.Reload(),
	); err != nil {
		log.Println("定时刷新失败", err)
		return
	}

	time.Sleep(10 * time.Second)

	if err := chromedp.Run(
		taskCtx,
		chromedp.Evaluate(getVb(), &vb),
		chromedp.Evaluate(getVt(), &vt),
		chromedp.Evaluate(getFb(), &fb),
		chromedp.Evaluate(getFt(), &ft),
		chromedp.FullScreenshot(&buf, 100),
	); err != nil {
		log.Println("定时刷新失败", err)
		return
	}

	if err := ioutil.WriteFile("end.png", buf, 0o644); err != nil {
		log.Println("输出图片失败", err)
		return
	}

	fmt.Println("数据情况", vb, vt, fb, ft)
	checkAndMail(vb, vt, fb, ft)
}

func getFb() string {
	return `document.getElementById('gprs').innerText`
}

func getFt() string {
	return `document.getElementById('sy_gprs').innerText`
}

func getVb() string {
	return `document.getElementById('balance').innerText`
}

func getVt() string {
	return `document.getElementById('highFeeT').innerText`
}
