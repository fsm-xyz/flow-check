use anyhow::Result as aResult;
use headless_chrome::{protocol::cdp::Page::CaptureScreenshotFormatOption, Browser, LaunchOptions};
use lettre::transport::smtp::authentication::Credentials;
use lettre::{Message, SmtpTransport, Transport};
use serde::{Deserialize, Serialize};

use std::env;
use std::ffi;
use std::fs;
use std::thread;
use std::time;

fn main() {
    let args: Vec<String> = env::args().collect();
    // println!("{:?}", args);
    if args.len() < 2 {
        panic!("请指定配置文件，可执行文件后跟上配置文件")
    }
    let c = build(args[1].to_string());

    login(c).unwrap();
}

fn login(c: Config) -> aResult<()> {
    let login_url = "https://login.10086.cn/html/login/touch.html?channelID=20290&backUrl=https%3A%2F%2Fsn.clientaccess.10086.cn%2Fhtml5%2Fsx%2FvfamilyN%2Findex.html&timestamp=1658228994563";
    // let login_url = "https://google.com";

    let tel = c.yd.tel;
    let args = String::from("--disable-web-security");

    let options = LaunchOptions::default_builder()
        .window_size(Some((390, 844)))
        .headless(false)
        .args(vec![ffi::OsStr::new(&args)])
        .build()
        .expect("Couldn't find appropriate Chrome binary.");

    let browser = Browser::new(options)?;
    let tab = browser.wait_for_initial_tab()?;
    let jpeg_data = tab
        .navigate_to(login_url)?
        .wait_until_navigated()?
        .capture_screenshot(CaptureScreenshotFormatOption::Jpeg, Some(100), None, true)?;
    fs::write("login1.jpg", &jpeg_data)?;

    let get_sms = format!(
        "
        document.getElementById('p_phone').value = {};
        console.log('输入验证码');
        document.getElementById('getSMSpwd').click();
        console.log('提交')",
        tel,
    );
    tab.evaluate(&get_sms, true)?;
    thread::sleep(time::Duration::new(30, 0));
    let code = fs::read_to_string("./code.txt")?;
    println!("code: {}", code);

    let s = format!(
        "
        document.getElementById('p_sms').value = {};
        console.log('输入验证码');
        document.getElementById('submit_bt').click();
        console.log('提交')",
        code
    );
    tab.evaluate(&s, true)?;
    thread::sleep(time::Duration::new(10, 0));

    let s = format!(
        "
        document.getElementById('p_sms').value = {};
        console.log('输入验证码');
        document.getElementById('submit_bt').click();
        console.log('提交')",
        code
    );
    tab.evaluate(&s, true)?;

    println!("Screenshots successfully created.");

    loop {
        let fb = 18.0;
        let ft = 40.0;
        let vb = 50.0;
        let vt = 100.0;
        if !c.quota.check_flow(vb, vt) {
            let body = format!("语音余额不足, 剩余: {}分钟, 总量: {}分钟", vb, vt);
            c.mail.send(body)
        }
        if !c.quota.check_voice(fb, ft) {
            let body = format!("流量余额不足, 剩余: {}G, 总量: {}G", fb, ft);
            c.mail.send(body)
        }

        thread::sleep(time::Duration::new(10, 0));
        let jpeg_data = tab
            .reload(true, None)?
            .wait_until_navigated()?
            .capture_screenshot(CaptureScreenshotFormatOption::Jpeg, Some(100), None, true)?;
        fs::write("screenshot.jpg", &jpeg_data)?;
    }

    // Ok(())
}

fn build(file_name: String) -> Config {
    let data = fs::read_to_string(file_name).unwrap();
    let c: Config = serde_json::from_str(data.as_str()).unwrap();

    // Ok(())
    return c;
}

#[derive(Serialize, Deserialize)]
struct Config {
    mail: Mail,
    quota: Quota,
    yd: YD,
}

#[derive(Serialize, Deserialize)]
struct YD {
    tel: String,
    token: String,
}

#[derive(Serialize, Deserialize)]
struct Quota {
    min_flow_balance_percent: f64,
    min_voice_balance_percent: f64,
    min_flow_balance: f64,
    min_voice_balance: f64,
}
#[derive(Serialize, Deserialize)]
struct Mail {
    from: String,
    to: String,
    password: String,
    subject: String,

    host: String,
}

impl Mail {
    fn send(&self, body: String) {
        let email = Message::builder()
            .from(self.from.parse().unwrap())
            .to(self.to.parse().unwrap())
            .subject(&self.subject)
            .body(body)
            .unwrap();

        let creds = Credentials::new(self.from.to_string(), self.password.to_string());

        // Open a remote connection to gmail
        let mailer = SmtpTransport::relay(&self.host)
            .unwrap()
            .credentials(creds)
            .build();

        // Send the email
        match mailer.send(&email) {
            Ok(_) => println!("Email sent successfully!"),
            Err(e) => panic!("Could not send email: {:?}", e),
        }
    }
}

impl Quota {
    fn check_voice(&self, vb: f64, vt: f64) -> bool {
        println!(
            "语音, 剩余: {}分钟, 总量: {}分钟, 最少: {}分钟, 最少百分比: {}",
            vb, vt, self.min_voice_balance, self.min_voice_balance_percent
        );
        if vb < self.min_voice_balance || vb < vt * self.min_voice_balance_percent {
            println!(
                "语音余额不足, 剩余: {}分钟, 总量: {}分钟, 最少: {}分钟, 最少百分比: {}",
                vb, vt, self.min_voice_balance, self.min_voice_balance_percent
            );
            return false;
        }
        true
    }

    fn check_flow(&self, fb: f64, ft: f64) -> bool {
        println!(
            "流量, 剩余: {}G, 总量: {}G, 最少: {}G, 最少百分比: {}",
            fb, ft, self.min_flow_balance, self.min_flow_balance_percent
        );
        if fb < self.min_flow_balance || fb < ft * self.min_flow_balance_percent {
            println!(
                "流量余额不足, 剩余: {}G, 总量: {}G, 最少: {}G, 最少百分比: {}",
                fb, ft, self.min_flow_balance, self.min_flow_balance_percent
            );
            return false;
        }
        true
    }
}
