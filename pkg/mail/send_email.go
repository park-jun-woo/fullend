package mail

import (
	"fmt"
	"net/smtp"
)

// @func sendEmail
// @description SMTP를 통해 이메일을 발송한다

type SendEmailInput struct {
	Host     string // SMTP 호스트 (예: smtp.gmail.com)
	Port     int    // SMTP 포트 (예: 587)
	Username string
	Password string
	From     string
	To       string
	Subject  string
	Body     string // plain text
}

type SendEmailOutput struct{}

func SendEmail(in SendEmailInput) (SendEmailOutput, error) {
	auth := smtp.PlainAuth("", in.Username, in.Password, in.Host)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		in.From, in.To, in.Subject, in.Body)
	addr := fmt.Sprintf("%s:%d", in.Host, in.Port)
	err := smtp.SendMail(addr, auth, in.From, []string{in.To}, []byte(msg))
	return SendEmailOutput{}, err
}
