package application

import (
	"bytes"
	"github.com/go-mail/mail"
	"github.com/mmmajder/zms-devops-auth-service/domain"
	"html/template"
	"log"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

func (service *EmailService) GetSendVerificationCodeEmailBody(receiverEmail string, verification domain.Verification) string {
	t, _ := template.ParseFiles(domain.VerificationEmailTemplate)
	var body bytes.Buffer
	t.Execute(&body, struct {
		Code int
		Url  string
	}{
		Code: verification.Code,
		Url:  "http://localhost/app/booking/auth/verify/" + verification.Id.Hex() + "?email=" + receiverEmail + "&userId=" + verification.UserId,
	})

	return body.String()
}

func (service *EmailService) SendEmail(subject, body string) {
	m := mail.NewMessage()
	m.SetHeader("From", domain.SenderEmailAddress)
	m.SetHeader("To", domain.SenderEmailAddress)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := mail.NewDialer(domain.EmailHost, domain.EmailPort, domain.SenderEmailAddress, domain.AppPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Fatal(err)
		log.Printf(err.Error())
	}
}
