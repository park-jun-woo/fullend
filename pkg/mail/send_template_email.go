package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
)

// @func sendTemplateEmail
// @description Go 템플릿으로 HTML 이메일을 발송한다

type SendTemplateEmailInput struct {
	Host         string
	Port         int
	Username     string
	Password     string
	From         string
	To           string
	Subject      string
	TemplateName string            // 템플릿 파일 경로 또는 인라인 템플릿
	Data         map[string]string // 템플릿 변수
}

type SendTemplateEmailOutput struct{}

func SendTemplateEmail(in SendTemplateEmailInput) (SendTemplateEmailOutput, error) {
	tmpl, err := template.New("email").Parse(in.TemplateName)
	if err != nil {
		return SendTemplateEmailOutput{}, err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, in.Data); err != nil {
		return SendTemplateEmailOutput{}, err
	}
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		in.From, in.To, in.Subject, body.String())
	auth := smtp.PlainAuth("", in.Username, in.Password, in.Host)
	addr := fmt.Sprintf("%s:%d", in.Host, in.Port)
	err = smtp.SendMail(addr, auth, in.From, []string{in.To}, []byte(msg))
	return SendTemplateEmailOutput{}, err
}
