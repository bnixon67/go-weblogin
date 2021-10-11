package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
)

const emailTmpl = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

{{ .Body }}
`

type MailMessage struct {
	From    string
	To      string
	Subject string
	Body    string
}

func SendEmail(smtpUser, smtpPassword, smtpHost, smtpPort, to, subject, body string) error {

	mailMessage := MailMessage{
		From:    smtpUser,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	t := template.Must(template.New("email").Parse(emailTmpl))
	message := &bytes.Buffer{}
	err := t.Execute(message, mailMessage)
	if err != nil {
		log.Print(err)
		return err
	}

	fmt.Print(message)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, mailMessage.From, []string{mailMessage.To}, message.Bytes())
	if err != nil {
		log.Print("SendMail failed")
		log.Print(err)
		return err
	}

	log.Print("Sent email")
	return err
}
