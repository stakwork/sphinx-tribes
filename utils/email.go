package utils

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stakwork/sphinx-tribes/logger"
)

type EmailService struct {
	apiKey        string
	fromEmail     string
	fromName      string
	supportEmail  string
}

func NewEmailService() *EmailService {
	return &EmailService{
		apiKey:       os.Getenv("SENDGRID_API_KEY"),
		fromEmail:    getEnvOrDefault("EMAIL_FROM_ADDRESS", "noreply@stakwork.com"),
		fromName:     getEnvOrDefault("EMAIL_FROM_NAME", "Sphinx Tribes"),
		supportEmail: getEnvOrDefault("SUPPORT_EMAIL", "support@stakwork.com"),
	}
}

func (e *EmailService) IsConfigured() bool {
	return e.apiKey != "" && e.supportEmail != ""
}

func (e *EmailService) SendNewUserNotification(userAlias, userPubKey, userUuid string) error {
	if !e.IsConfigured() {
		logger.Log.Info("Email service not configured, skipping new user notification")
		return nil
	}

	from := mail.NewEmail(e.fromName, e.fromEmail)
	to := mail.NewEmail("Support Team", e.supportEmail)
	
	subject := "New User Profile Created - Sphinx Tribes"
	
	plainTextContent := fmt.Sprintf(`A new user has created a profile on Sphinx Tribes.

User Details:
- Alias: %s
- Public Key: %s
- UUID: %s
- Profile URL: https://community.sphinx.chat/p/%s

This is an automated notification.`, userAlias, userPubKey, userUuid, userUuid)

	htmlContent := fmt.Sprintf(`
	<html>
	<body>
		<h2>New User Profile Created - Sphinx Tribes</h2>
		<p>A new user has created a profile on Sphinx Tribes.</p>
		
		<h3>User Details:</h3>
		<ul>
			<li><strong>Alias:</strong> %s</li>
			<li><strong>Public Key:</strong> %s</li>
			<li><strong>UUID:</strong> %s</li>
			<li><strong>Profile URL:</strong> <a href="https://community.sphinx.chat/p/%s">View Profile</a></li>
		</ul>
		
		<p><em>This is an automated notification.</em></p>
	</body>
	</html>`, userAlias, userPubKey, userUuid, userUuid)

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)

	client := sendgrid.NewSendClient(e.apiKey)
	response, err := client.Send(message)
	
	if err != nil {
		logger.Log.Error("Failed to send new user notification email: %v", err)
		return err
	}

	if response.StatusCode >= 400 {
		logger.Log.Error("SendGrid API returned error status %d: %s", response.StatusCode, response.Body)
		return fmt.Errorf("email service returned status %d", response.StatusCode)
	}

	logger.Log.Info("New user notification email sent successfully for user %s (UUID: %s)", userAlias, userUuid)
	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}