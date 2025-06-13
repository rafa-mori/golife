package email

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/rafa-mori/logz"
	"io"
)

// EmailService represents the email service with SMTP and IMAP configurations.
type EmailService struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	IMAPHost     string
	IMAPPort     string
	IMAPUsername string
	IMAPPassword string
}

// NewEmailService creates a new EmailService instance with the provided SMTP and IMAP configuration.
func NewEmailService(smtpHost, smtpPort, smtpUsername, smtpPassword, imapHost, imapPort, imapUsername, imapPassword string) *EmailService {
	return &EmailService{
		SMTPHost:     smtpHost,
		SMTPPort:     smtpPort,
		SMTPUsername: smtpUsername,
		SMTPPassword: smtpPassword,
		IMAPHost:     imapHost,
		IMAPPort:     imapPort,
		IMAPUsername: imapUsername,
		IMAPPassword: imapPassword,
	}
}

// ParseEmailBody parses the email body and returns the body as a string.
func (es *EmailService) ParseEmailBody(msgReader *mail.Part) (string, error) {
	// Check if the message reader is nil.
	if msgReader == nil {
		logz.Error("Message reader is nil", nil)
		return "", fmt.Errorf("message reader is nil")
	}

	// Read the email body.
	body := msgReader.Body

	bodyContent, err := io.ReadAll(body)
	if err != nil {
		logz.Error("Failed to read email body", map[string]interface{}{"error": err})
		return "", err
	}
	return string(bodyContent), nil
}

// getAddressListSlice parses the email header and returns a map of address lists.
func getAddressListSlice(header *imap.Envelope) (map[string][]*imap.Address, error) {
	addressList := make(map[string][]*imap.Address)
	addressList["From"] = header.From
	addressList["Sender"] = header.Sender
	addressList["Reply-To"] = header.ReplyTo
	addressList["To"] = header.To
	addressList["Cc"] = header.Cc
	addressList["Bcc"] = header.Bcc
	return addressList, nil
}
