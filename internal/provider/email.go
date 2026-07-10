package provider

import (
	"context"
	"net/smtp"
)

type EmailProvider struct {
	smtpHost     string
	smtpPort     string
	smtpUsername string
	smtpPassword string
	smtpSender   string
}

func NewEmailProvider(smtpHost, smtpPort, smtpUsername, smtpPassword, smtpSender string) *EmailProvider {
	return &EmailProvider{
		smtpHost:     smtpHost,
		smtpPort:     smtpPort,
		smtpUsername: smtpUsername,
		smtpPassword: smtpPassword,
		smtpSender:   smtpSender,
	}
}

func (p *EmailProvider) SendEmail(ctx context.Context, to, subject, body string) error {
	from := p.smtpSender
	if from == "" {
		from = "noreply@example.com"
	}
	message := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-Version: 1.0\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\n" +
		"\n" +
		body

	address := p.smtpHost + ":" + p.smtpPort
	return smtp.SendMail(address, smtp.PlainAuth("", p.smtpUsername, p.smtpPassword, p.smtpHost), from, []string{to}, []byte(message))
}
