package service

import (
	"dryve/internal/config"
	"dryve/internal/dto"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(email dto.Email) error
}

// Default emailService implementing EmailService
type emailService struct {
	config config.EmailConfig
	dialer *gomail.Dialer
}

func NewEmailService(c config.EmailConfig) EmailService {
	dialer := gomail.NewDialer(c.Host, c.Port, c.User, c.Password)
	return &emailService{
		config: c,
		dialer: dialer,
	}
}

// SendEmail uses gomail dialer to send an email through SMTP
func (s emailService) SendEmail(m dto.Email) error {
	msg := gomail.NewMessage()

	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", m.Body)
	for _, a := range m.Attachments {
		msg.Attach(a)
	}

	return s.dialer.DialAndSend(msg)
}
