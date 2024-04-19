package mailer

import (
	"database/sql"
	"net/smtp"
	"os"
	"shopscraper/pkg/config"
	"shopscraper/pkg/models"
	"strings"
	"testing"
)

type MockSmtpSender struct {
	SendMailFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
	Calls        []struct {
		Addr string
		Auth smtp.Auth
		From string
		To   []string
		Msg  []byte
	}
}

func (m *MockSmtpSender) SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	m.Calls = append(m.Calls, struct {
		Addr string
		Auth smtp.Auth
		From string
		To   []string
		Msg  []byte
	}{Addr: addr, Auth: a, From: from, To: to, Msg: msg})
	return m.SendMailFunc(addr, a, from, to, msg)
}

func TestSendEmail(t *testing.T) {
	// Save current value (if exists) and defer its restoration
	originalValue, isSet := os.LookupEnv("SHOPSCRAPER_SMTP_PASSWORD")
	defer func() {
		if isSet {
			os.Setenv("SHOPSCRAPER_SMTP_PASSWORD", originalValue)
		} else {
			os.Unsetenv("SHOPSCRAPER_SMTP_PASSWORD")
		}
	}()

	// Set the environment variable for the purpose of the test
	os.Setenv("SHOPSCRAPER_SMTP_PASSWORD", "test")

	products := []models.Product{
		{Name: "Product 1", Price: 10, Link: "https://example.com/product1"},
		{Name: "Product 2", PreviousPrice: sql.NullInt64{Int64: 10, Valid: true}, Price: 19, Link: "https://example.com/product2"},
	}
	programConfig := config.ProgramConfig{
		Email: config.EmailConfig{
			Recipient: "recipient@example.com",
			Sender:    "sender@example.com",
			Subject:   "New Products",
			Server:    "smtp.example.com",
			Port:      "587",
		},
	}
	mockSender := &MockSmtpSender{
		SendMailFunc: func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		},
	}

	err := SendEmail(mockSender, products, programConfig)
	if err != nil {
		t.Errorf("SendEmail() failed: %v", err)
	}

	if len(mockSender.Calls) != 1 {
		t.Errorf("Expected SendMail to be called once, got %v", len(mockSender.Calls))
	}

	call := mockSender.Calls[0]
	if call.From != programConfig.Email.Sender {
		t.Errorf("Expected sender email to be '%v', got '%v'", programConfig.Email.Sender, call.From)
	}

	if !strings.Contains(string(call.Msg), products[0].Name) || !strings.Contains(string(call.Msg), products[1].Name) {
		t.Errorf("Email body does not contain product names")
	}

	if call.To[0] != programConfig.Email.Recipient {
		t.Errorf("Expected recipient email to be '%v', got '%v'", programConfig.Email.Recipient, call.To[0])
	}
}

// pkg/mailer/mailer_test.go

func TestSendEmail_NoProducts(t *testing.T) {
	// Set up environment variable for SMTP password
	originalPassword, isSet := os.LookupEnv("SHOPSCRAPER_SMTP_PASSWORD")
	defer func() {
		if isSet {
			os.Setenv("SHOPSCRAPER_SMTP_PASSWORD", originalPassword)
		} else {
			os.Unsetenv("SHOPSCRAPER_SMTP_PASSWORD")
		}
	}()
	os.Setenv("SHOPSCRAPER_SMTP_PASSWORD", "test")

	// Define test data
	products := []models.Product{} // No products
	programConfig := config.ProgramConfig{
		Email: config.EmailConfig{
			Recipient: "recipient@example.com",
			Sender:    "sender@example.com",
			Subject:   "New Products",
			Server:    "smtp.example.com",
			Port:      "587",
		},
	}
	mockSender := &MockSmtpSender{
		SendMailFunc: func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			return nil
		},
	}

	// Call the function under test
	err := SendEmail(mockSender, products, programConfig)

	// Check that no error was returned
	if err != nil {
		t.Errorf("SendEmail() should not return an error when there are no products: %v", err)
	}

	// Check that SendMail was not called
	if len(mockSender.Calls) != 0 {
		t.Errorf("Expected SendMail to not be called, but it was called %v times", len(mockSender.Calls))
	}
}
