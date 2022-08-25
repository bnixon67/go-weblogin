/*
Copyright 2022 Bill Nixon

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/
package weblogin

import (
	"bytes"
	"log"
	"net/smtp"
	"text/template"
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
		log.Printf("SendMail failed: %v", err)
		return err
	}

	log.Print("Sent email")
	return err
}
