package mail

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

// @func sendTemplateEmail
// @description Sends HTML email using Go templates

type SendTemplateEmailRequest struct {
	To           string
	Subject      string
	TemplateName string
}

type SendTemplateEmailResponse struct{}

func SendTemplateEmail(req SendTemplateEmailRequest) (SendTemplateEmailResponse, error) {
	host := os.Getenv("SMTP_HOST")
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM")

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, req.To, req.Subject, req.TemplateName)
	auth := smtp.PlainAuth("", username, password, host)
	addr := fmt.Sprintf("%s:%d", host, port)
	err := smtp.SendMail(addr, auth, from, []string{req.To}, []byte(msg))
	return SendTemplateEmailResponse{}, err
}
