package service

import (
	"fmt"

	"github.com/avalonprod/invoicepro/server/internal/config"
	"github.com/avalonprod/invoicepro/server/pkg/email"
)

type EmailService struct {
	sender email.Sender
	config config.EmailConfig
}

func NewEmailService(sender email.Sender, config config.EmailConfig) *EmailService {
	return &EmailService{
		sender: sender,
		config: config,
	}
}

type VerificationEmailInput struct {
	Name             string
	VerificationCode string
	Email            string
}

func (e *EmailService) SendUserVerificationEmail(input VerificationEmailInput) error {
	subject := fmt.Sprintf(e.config.Subjects.Verification, input.Name)
	sendInput := email.SendEmailInput{To: input.Email, Subject: subject}

	templateInput := VerificationEmailInput{Name: input.Name, VerificationCode: input.VerificationCode}

	err := sendInput.GenerateBodyFromHTML(e.config.Templates.Verification, templateInput)
	if err != nil {
		return err
	}
	err = e.sender.Send(sendInput)
	return err
}
