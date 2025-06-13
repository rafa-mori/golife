package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rafa-mori/golife/internal"
	"github.com/rafa-mori/logz"
	"net/http"
	"net/smtp"
)

type NotificationManager struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMSGateway   string
}

func (nm *NotificationManager) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", nm.SMTPUsername, nm.SMTPPassword, nm.SMTPHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail(nm.SMTPHost+":"+nm.SMTPPort, auth, nm.SMTPUsername, []string{to}, msg)
	if err != nil {
		logz.Error("Failed to send email", map[string]interface{}{"error": err})
		return err
	}
	logz.Info("Email sent successfully", nil)
	return nil
}

func (nm *NotificationManager) SendSMS(to, message string) error {
	data := map[string]string{
		"to":      to,
		"message": message,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		logz.Error("Failed to marshal SMS data", map[string]interface{}{"error": err})
		return err
	}
	resp, err := http.Post(nm.SMSGateway, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		logz.Error("Failed to send SMS", map[string]interface{}{"error": err})
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logz.Error("Failed to send SMS", map[string]interface{}{"status": resp.StatusCode})
		return fmt.Errorf("failed to send SMS, status code: %d", resp.StatusCode)
	}
	logz.Info("SMS sent successfully", nil)
	return nil
}

func (nm *NotificationManager) Notify(event internal.IManagedProcessEvents) {
	// Example notification logic
	if event.Event() == "process_started" {
		nm.SendEmail("user@example.com", "Process Started", "A process has started.")
		nm.SendSMS("+1234567890", "A process has started.")
	}
}

func NewNotificationManager(smtpHost, smtpPort, smtpUsername, smtpPassword, smsGateway string) *NotificationManager {
	return &NotificationManager{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		SMSGateway:   smsGateway,
	}
}
