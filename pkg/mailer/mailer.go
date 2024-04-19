package mailer

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"shopscraper/pkg/config"
	"shopscraper/pkg/models"
)

type SmtpSender interface {
	SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

type RealSmtpSender struct{}

func (s *RealSmtpSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	return smtp.SendMail(addr, a, from, to, msg)
}

var (
	ErrConfigInvalid = errors.New("configuration is invalid")
)

func validateConfig(cfg config.EmailConfig) error {
	if cfg.Recipient == "" || cfg.Sender == "" || cfg.Subject == "" || cfg.Server == "" || cfg.Port == "" {
		return ErrConfigInvalid
	}
	return nil
}

func SendEmail(smtpSender SmtpSender, products []models.Product, programConfig config.ProgramConfig) error {
	if len(products) == 0 {
		log.Println("No products to send; skipping email.")
		return nil
	}

	if err := validateConfig(programConfig.Email); err != nil {
		log.Printf("Invalid email configuration: %v\n", err)
		return fmt.Errorf("%w: %v", ErrConfigInvalid, err)
	}

	password := os.Getenv("SHOPSCRAPER_SMTP_PASSWORD")
	if password == "" {
		log.Println("SMTP password is not set.")
		return ErrConfigInvalid
	}

	body := constructEmailBody(products)

	msg := fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: %s\n\n%s",
		programConfig.Email.Sender,
		programConfig.Email.Recipient,
		programConfig.Email.Subject,
		body,
	)

	log.Printf("Sending email to %s with %d items", programConfig.Email.Recipient, len(products))
	err := smtpSender.SendMail(
		fmt.Sprintf("%s:%s", programConfig.Email.Server, programConfig.Email.Port),
		smtp.PlainAuth("", programConfig.Email.Sender, password, programConfig.Email.Server),
		programConfig.Email.Sender, []string{programConfig.Email.Recipient}, []byte(msg),
	)

	if err != nil {
		log.Printf("Failed to send email: %v", err)
		return err
	}

	return nil
}

func constructEmailBody(products []models.Product) string {
	body := ""
	for _, p := range products {
		line := fmt.Sprintf("%s - %d", p.Name, p.Price)
		// If previous price is valid, include it
		if p.PreviousPrice.Valid {
			line += fmt.Sprintf(" (%d)", p.PreviousPrice.Int64)
		}
		line += fmt.Sprintf(" - %s\n", p.Shop)
		line += fmt.Sprintf("%s\n\n", p.Link)

		body += line
	}
	return body
}
