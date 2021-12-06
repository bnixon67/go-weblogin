package main

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
)

const emailTmpl = `From: {{ .From }}
To: {{ .To }}
Subject: {{ .Subject }}

{{ .Body }}
`

// MailMessage contains data to include in the email template.
type MailMessage struct {
	From    string
	To      string
	Subject string
	Body    string
}

// SendEmail will send an email using the values provided.
func SendEmail(smtpUser, smtpPassword, smtpHost, smtpPort, to, subject, body string) error {
	mailMessage := MailMessage{
		From:    smtpUser,
		To:      to,
		Subject: subject,
		Body:    body,
	}

	// TODO: cache template
	t, err := template.New("email").Parse(emailTmpl)
	if err != nil {
		log.Printf("unable to parse template, %v", err)
		return err
	}

	// fill message template
	message := &bytes.Buffer{}
	err = t.Execute(message, mailMessage)
	if err != nil {
		log.Print(err)
		return err
	}

	// authenticate to SMTP server
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	// send email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, mailMessage.From, []string{mailMessage.To}, message.Bytes())
	if err != nil {
		log.Print("SendMail failed")
		log.Print(err)
		return err
	}

	log.Print("Sent email")
	return err
}
