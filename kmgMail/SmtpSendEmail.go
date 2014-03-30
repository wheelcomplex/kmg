package kmgMail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
)

type SmtpSendEmailRequest struct {
	SmtpHost     string
	SmtpPort     int
	From         string //also as stmp username
	SmtpPassword string
	To           []string
	Subject      string
	Message      string
}

func SmtpSendEmail(req SmtpSendEmailRequest) (err error) {
	parameters := &struct {
		From    string
		To      string
		Subject string
		Message string
	}{
		req.From,
		strings.Join([]string(req.To), ","),
		req.Subject,
		req.Message,
	}

	buffer := new(bytes.Buffer)

	t := template.Must(template.New("emailTemplate").Parse(_EmailScript()))
	t.Execute(buffer, parameters)

	auth := smtp.PlainAuth("", req.From, req.SmtpPassword, req.SmtpHost)

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", req.SmtpHost, req.SmtpPort),
		auth,
		req.From,
		req.To,
		buffer.Bytes())

	return err
}

// _EmailScript returns a template for the email message to be sent
func _EmailScript() (script string) {
	return `From: {{.From}}
To: {{.To}}
Subject: {{.Subject}}
MIME-version: 1.0
Content-Type: text/html; charset="UTF-8"

{{.Message}}`
}
