package service

import (
	"dryve/internal/config"
	"dryve/internal/dto"
	"fmt"
	"strings"
)

// mockEmailService implementing EmailService
type mockEmailService struct {
	config config.EmailConfig
}

// NewMockEmailService creates a mock email implementation
// writing out to stdout the sent email.
func NewMockEmailService(c config.EmailConfig) EmailService {
	return &mockEmailService{
		config: c,
	}
}

// SendEmail just prints the email to stdout.
func (s mockEmailService) SendEmail(m dto.Email) error {
	var msg *strings.Builder

	msg.WriteString("From: ")
	msg.WriteString(m.From)
	msg.WriteString("\nTo: ")
	msg.WriteString(m.To)
	msg.WriteString("\nSubject: ")
	msg.WriteString(m.Subject)
	msg.WriteString("\nBody:\n")
	msg.WriteString(m.Body)

	if len(m.Attachments) > 0 {
		msg.WriteString("\nAttachments:")
		for _, a := range m.Attachments {
			msg.WriteString("\n\t- " + a)
		}
	}

	fmt.Println(msg.String())
	return nil
}
