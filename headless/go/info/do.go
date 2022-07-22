package info

import (
	"context"
	"fmt"
	"headless-go/mail"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	once    sync.Once
	Running bool
)

func Run() {
	once.Do(headless)
}

var M *mail.Mail

const loginURL = "https://login.10086.cn/html/login/touch.html?channelID=20290&backUrl=https%3A%2F%2Fsn.clientaccess.10086.cn%2Fhtml5%2Fsx%2FvfamilyN%2Findex.html&timestamp=1658228994563"

func headless() {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.DisableGPU,
		chromedp.WindowSize(390, 844),
		chromedp.UserDataDir(dir),
	)

	fmt.Println(C.B.Headless)
	if !C.B.Headless {
		opts = append(opts, chromedp.Flag("headless", false))
	}

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	taskCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	//chromedp监听网页上弹出alert对话框
	chromedp.ListenTarget(taskCtx, func(ev interface{}) {
		if ev, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			fmt.Println("closing alert:", ev.Message)
			go func() {
				//自动关闭alert对话框
				if err := chromedp.Run(taskCtx,
					//注释掉下一行可以更清楚地看到效果
					page.HandleJavaScriptDialog(true),
				); err != nil {
					panic(err)
				}
			}()
		}
	})

	// fmt.Println(getSMS(tel))
	if err := chromedp.Run(
		taskCtx,
		chromedp.Navigate(loginURL),
		chromedp.Evaluate(getSMS(C.YD.Tel), nil),
	); err != nil {
		log.Fatal(err)
	}

	// fmt.Println(submit(getCode()))
	time.Sleep(time.Duration(C.YD.WaitSMSTime) * time.Second)
	chromedp.Run(
		taskCtx,
		chromedp.Evaluate(submit(getCode()), nil),
	)

	// chromedp.
	var (
		buf    []byte
		vb, vt string
		fb, ft string
	)
	for {
		log.Println("定时刷新")
		Running = false
		if err = chromedp.Run(
			taskCtx,
			chromedp.Reload(),
		); err != nil {
			log.Fatal("定时刷新失败", err)
		}

		time.Sleep(10 * time.Second)

		if err = chromedp.Run(
			taskCtx,
			chromedp.Evaluate(getVb(), &vb),
			chromedp.Evaluate(getVt(), &vt),
			chromedp.Evaluate(getFb(), &fb),
			chromedp.Evaluate(getFt(), &ft),
			chromedp.FullScreenshot(&buf, 100),
		); err != nil {
			log.Fatal("定时刷新失败", err)
		}

		if err := ioutil.WriteFile("end.png", buf, 0o644); err != nil {
			log.Fatal(err)
		}
		Running = true
		fmt.Println("数据情况", vb, vt, fb, ft)
		checkAndMail(vb, vt, fb, ft)
		time.Sleep(time.Duration(C.YD.ReloadTime) * time.Second)
	}
}

func getCode() string {
	data, err := ioutil.ReadFile("./code.txt")
	if err != nil {
		log.Fatal("读取code失败", err)
	}

	return string(data)
}

func getSMS(tel string) string {
	return fmt.Sprintf(`
		document.getElementById('p_phone').value = %s;
		console.log('获取验证码');
		document.getElementById('getSMSpwd').click();
		console.log('提交')`,
		tel)
}

func submit(code string) string {
	return fmt.Sprintf(`
		document.getElementById('p_sms').value = %s;
		console.log('输入验证码');
		document.getElementById('submit_bt').click();
		console.log('提交')`, code)
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
