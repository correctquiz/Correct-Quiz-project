package service

import (
	"context"
	"fmt"
	"log"
	"os"

	brevo "github.com/getbrevo/brevo-go/lib"
)

type EmailServiceInterface interface {
	SendEmail(toEmail, toName, subject, htmlContent string) error
}

type BrevoEmailService struct {
	apiKey string
	client *brevo.APIClient
}

func NewBrevoEmailService() EmailServiceInterface {
	apiKey := os.Getenv("BREVO_API_KEY")

	if apiKey == "" {
		log.Println("WARNING: BREVO_API_KEY is not set. Emails will not be sent.")
	}

	cfg := brevo.NewConfiguration()
	cfg.AddDefaultHeader("api-key", apiKey)

	client := brevo.NewAPIClient(cfg)

	return &BrevoEmailService{
		apiKey: apiKey,
		client: client,
	}
}

func (s *BrevoEmailService) SendEmail(toEmail, toName, subject, htmlContent string) error {

	if s.apiKey == "" {
		return fmt.Errorf("cannot send email: Brevo API key is not configured")
	}

	sender := brevo.SendSmtpEmailSender{
		Name:  "Thun Thimrattanakul",
		Email: "thunthim2546@gmail.com",
	}

	to := []brevo.SendSmtpEmailTo{
		{
			Email: toEmail,
			Name:  toName,
		},
	}

	smtpEmail := brevo.SendSmtpEmail{
		Sender:      &sender,
		To:          to,
		HtmlContent: htmlContent,
		Subject:     subject,
	}

	_, _, err := s.client.TransactionalEmailsApi.SendTransacEmail(context.Background(), smtpEmail)

	if err != nil {
		log.Printf("Brevo SendEmail error: %v\n", err)
		return fmt.Errorf("failed to send email via Brevo: %w", err)
	}

	log.Printf("Brevo email sent successfully to %s\n", toEmail)
	return nil
}
