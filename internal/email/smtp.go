package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"rizon-test-task/internal/config"
)

// smtpSender implements EmailSender using SMTP.
type smtpSender struct {
	config *config.EmailConfig
}

// NewSMTPSender returns an SMTP implementation of EmailSender.
func NewSMTPSender(cfg *config.EmailConfig) EmailSender {
	return &smtpSender{
		config: cfg,
	}
}

// SendEmail sends an email using SMTP.
func (s *smtpSender) SendEmail(ctx context.Context, to, subject, body string) error {
	// Validate configuration
	if s.config.Username == "" {
		return fmt.Errorf("email username not configured")
	}
	if s.config.Password == "" {
		return fmt.Errorf("email password not configured")
	}

	// Create SMTP authentication
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.SMTPHost)

	// Format email message with headers
	from := s.config.FromEmail
	if from == "" {
		from = s.config.Username
	}

	// Build email message
	message := buildEmailMessage(from, to, subject, body)

	// Send email
	addr := s.config.Addr()
	
	// For Gmail (port 587), we need to use TLS
	if s.config.SMTPPort == 587 {
		return s.sendWithTLS(ctx, addr, auth, from, to, []byte(message))
	}
	
	// For port 465 (SSL), use SendMail directly (it handles SSL)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(message))
}

// sendWithTLS sends email with TLS (for port 587).
func (s *smtpSender) sendWithTLS(ctx context.Context, addr string, auth smtp.Auth, from, to string, message []byte) error {
	// Connect to SMTP server
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer client.Close()

	// Check if TLS is supported
	if ok, _ := client.Extension("STARTTLS"); ok {
		config := &tls.Config{
			ServerName: s.config.SMTPHost,
		}
		if err := client.StartTLS(config); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Set sender
	if err := client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipient
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to open data connection: %w", err)
	}

	if _, err := writer.Write(message); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write email body: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close data connection: %w", err)
	}

	// Quit
	if err := client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP connection: %w", err)
	}

	return nil
}

// buildEmailMessage formats the email message with proper headers.
func buildEmailMessage(from, to, subject, body string) string {
	// Escape special characters in headers
	from = escapeHeader(from)
	to = escapeHeader(to)
	subject = escapeHeader(subject)

	// Determine content type based on body
	contentType := "text/plain; charset=UTF-8"
	if strings.Contains(body, "<html") || strings.Contains(body, "<HTML") {
		contentType = "text/html; charset=UTF-8"
	}

	// Build message with headers
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += fmt.Sprintf("Content-Type: %s\r\n", contentType)
	message += "\r\n"
	message += body

	return message
}

// escapeHeader escapes special characters in email headers.
func escapeHeader(header string) string {
	// Simple escaping - wrap in quotes if contains special characters
	if strings.ContainsAny(header, "<>@,;:\\\"[]") {
		return fmt.Sprintf("\"%s\"", strings.ReplaceAll(header, "\"", "\\\""))
	}
	return header
}
