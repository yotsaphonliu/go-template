package smtp_service

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
)

type EmailData struct {
	TxHeaderDisplay    string
	TxNameTitle        string
	TxNameValue        string
	TxOwnerTitle       string
	TxOwnerValue       string
	TxCreatedTimeTitle string
	TxCreatedTimeValue string
	TxSubHeader        string
	Url                string
}

func (sc *SmtpServiceClient) Send(to []string, subject string, templatePath string, data any) error {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", sc.Config.From)
	msg.SetHeader("To", to...)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body.String())

	//msg.Embed("src/core/smtp_service/templates/images/logo.png")
	//msg.Embed("src/core/smtp_service/templates/images/check.png")

	dialer := gomail.NewDialer(
		sc.Config.SMTPHost,
		sc.Config.SMTPPort,
		sc.Config.Username,
		sc.Config.Password,
	)

	return dialer.DialAndSend(msg)
}
