package api

import (
	"fmt"
	"net/smtp"
	"net/mail"
	"net"
	"time"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/faelmori/logz"
)

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

func (es *EmailService) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", es.SMTPUsername, es.SMTPPassword, es.SMTPHost)
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	err := smtp.SendMail(es.SMTPHost+":"+es.SMTPPort, auth, es.SMTPUsername, []string{to}, msg)
	if err != nil {
		logz.Error("Failed to send email", map[string]interface{}{"error": err})
		return err
	}
	logz.Info("Email sent successfully", nil)
	return nil
}

func (es *EmailService) ReceiveEmails() ([]*mail.Message, error) {
	c, err := client.DialTLS(es.IMAPHost+":"+es.IMAPPort, nil)
	if err != nil {
		logz.Error("Failed to connect to IMAP server", map[string]interface{}{"error": err})
		return nil, err
	}
	defer c.Logout()

	if err := c.Login(es.IMAPUsername, es.IMAPPassword); err != nil {
		logz.Error("Failed to login to IMAP server", map[string]interface{}{"error": err})
		return nil, err
	}

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		logz.Error("Failed to select INBOX", map[string]interface{}{"error": err})
		return nil, err
	}

	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope, imap.FetchBody}, messages)
	}()

	var result []*mail.Message
	for msg := range messages {
		r := msg.GetBody(&imap.BodySectionName{})
		if r == nil {
			logz.Warn("Server didn't return message body", nil)
			continue
		}

		m, err := mail.CreateReader(r)
		if err != nil {
			logz.Error("Failed to create mail reader", map[string]interface{}{"error": err})
			continue
		}

		result = append(result, m)
	}

	if err := <-done; err != nil {
		logz.Error("Failed to fetch emails", map[string]interface{}{"error": err})
		return nil, err
	}

	logz.Info("Emails received successfully", nil)
	return result, nil
}

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
