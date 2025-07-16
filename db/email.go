package db

import (
	"fmt"
	"net/smtp"
	"time"

	"github.com/stakwork/sphinx-tribes/config"
	"github.com/stakwork/sphinx-tribes/logger"
)

type EmailNotification struct {
	To      string
	Subject string
	Body    string
}

// SendNewUserNotification sends an email notification when a new user creates a profile
func (db database) SendNewUserNotification(person Person) {
	// Check if email configuration is available
	if config.SMTPHost == "" || config.SMTPUsername == "" || config.SMTPPassword == "" {
		logger.Log.Info("Email notification: SMTP configuration not found, skipping email")
		return
	}

	// Prepare email content
	subject := "New User Profile Created on Sphinx Tribes"
	body := fmt.Sprintf(`A new user has created a profile on Sphinx Tribes.

Profile Details:
- User Alias: %s
- Owner Public Key: %s
- Description: %s
- Created: %s
- Profile UUID: %s

You can view the profile at: %s/p?owner_id=%s

Best regards,
Sphinx Tribes System
`, 
		person.OwnerAlias,
		person.OwnerPubKey,
		person.Description,
		person.Created.Format(time.RFC3339),
		person.Uuid,
		config.Host,
		person.OwnerPubKey,
	)

	notification := EmailNotification{
		To:      config.SupportEmail,
		Subject: subject,
		Body:    body,
	}

	// Send email asynchronously
	go func() {
		err := sendEmail(notification)
		if err != nil {
			logger.Log.Error("Email notification: Failed to send new user notification: %v", err)
		} else {
			logger.Log.Info("Email notification: Successfully sent new user notification for %s", person.OwnerAlias)
		}
	}()
}

// sendEmail sends an email using SMTP
func sendEmail(notification EmailNotification) error {
	// SMTP server configuration
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort
	smtpAddr := smtpHost + ":" + smtpPort

	// Email authentication
	auth := smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, smtpHost)

	// Email message
	message := fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s\r\n",
		notification.To,
		config.EmailFrom,
		notification.Subject,
		notification.Body,
	)

	// Recipients
	recipients := []string{notification.To}

	// Send email
	err := smtp.SendMail(smtpAddr, auth, config.EmailFrom, recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendTestEmail sends a test email to verify configuration
func SendTestEmail(toEmail string) error {
	if config.SMTPHost == "" || config.SMTPUsername == "" || config.SMTPPassword == "" {
		return fmt.Errorf("SMTP configuration not available")
	}

	notification := EmailNotification{
		To:      toEmail,
		Subject: "Sphinx Tribes Email Test",
		Body: fmt.Sprintf(`This is a test email from Sphinx Tribes.

Email configuration is working correctly.

Sent at: %s

Best regards,
Sphinx Tribes System
`, time.Now().Format(time.RFC3339)),
	}

	return sendEmail(notification)
}