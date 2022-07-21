package mail

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type Mail struct {
	To       string
	From     string
	Password string

	Subject string

	Host string
	Port int
}

func (m *Mail) Send(data string) {
	d := gomail.NewDialer(m.Host, m.Port, m.From, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/plain", data)

	if err := d.DialAndSend(msg); err != nil {
		fmt.Println("发送邮件失败", err, data)
		return
	}
}
