package mail

import (
	"crypto/tls"
	"fmt"

	gomail "gopkg.in/gomail.v2"
)

type Mail struct {
	MailTo   string
	MailFrom string
	Password string

	Subject string

	Host string
	Port int
}

func (m *Mail) Send(data string) {
	d := gomail.NewDialer(m.Host, m.Port, m.MailFrom, m.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.MailFrom)
	msg.SetHeader("To", m.MailTo)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/plain", data)

	if err := d.DialAndSend(msg); err != nil {
		fmt.Println("发送邮件失败", err, data)
		return
	}
}
