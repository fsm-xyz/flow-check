package info

import (
	"context"
	"fmt"
	"headless-go/mail"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

var (
	Running bool
)

func Run() {
	if !Running {
		headless()
	}
}

var M *mail.Mail

const loginURL = "https://login.10086.cn/html/login/touch.html?channelID=20290&backUrl=https%3A%2F%2Fsn.clientaccess.10086.cn%2Fhtml5%2Fsx%2FvfamilyN%2Findex.html&timestamp=1658228994563"

func headless() {
	dir, err := ioutil.TempDir("", "chromedp-example")
	if err != nil {
		log.Println("移除tmp目录失败", err)
		return
	}
	defer os.RemoveAll(dir)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(`Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1`),
		chromedp.DisableGPU,
		chromedp.WindowSize(390, 844),
		chromedp.UserDataDir(dir),
	)

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
		log.Println("获取短信失败", err)
		return
	}

	// fmt.Println(submit(getCode()))
	time.Sleep(time.Duration(C.YD.WaitSMSTime) * time.Second)
	chromedp.Run(
		taskCtx,
		chromedp.Evaluate(submit(getCode()), nil),
	)

	time.Sleep(10 * time.Second)

	// 登录成功，接下来可以通过刷新页面获取数据，或者直接调用接口
	// 刷新页面方式的话在linux服务器无法成功渲染
	loop(taskCtx)
}

func loop(taskCtx context.Context) {
	defer func() { Running = false }()
	var cookie string
	if C.Api {
		if err := chromedp.Run(
			taskCtx,
			chromedp.Evaluate(getCookie(), &cookie),
		); err != nil {
			log.Println("获取token失败", err)
			return
		}

		RefreshToken(getToken(cookie))
	}

	for {
		log.Println(C.YD.ReloadTime, "s定时刷新")
		Running = false

		if C.Api {
			getStatus()
		} else {
			Loadpage(taskCtx)
		}

		Running = true
		time.Sleep(time.Duration(C.YD.ReloadTime) * time.Second)
	}

}

func getCode() string {
	data, err := ioutil.ReadFile("./code.txt")
	if err != nil {
		log.Println("读取code失败", err)
		return ""
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

func getCookie() string {
	return `document.cookie`
}

func getToken(s string) string {
	for _, x := range strings.Split(s, ";") {
		if strings.Contains(x, "JSESSIONID") {
			return strings.Trim(x, " ")
		}
	}
	return ""
}
