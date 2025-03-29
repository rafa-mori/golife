package email

import (
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-message/mail"
	"github.com/faelmori/logz"
	"io"
	"strings"
)

// ParseEmail parses the entire email, including the header, body, and attachments.
func ParseEmail(msg *imap.Message) (*Email, error) {
	if msg == nil {
		logz.Error("Message is nil", nil)
		return nil, fmt.Errorf("message is nil")
	}

	headerReader, err := mail.CreateReader(msg.GetBody(&imap.BodySectionName{}))
	if err != nil {
		logz.Error("Failed to create mail reader", map[string]interface{}{"error": err})
		return nil, err
	}

	mlHeader, err := ParseEmailHeader(&headerReader.Header)
	if err != nil {
		logz.Error("Failed to parse email header", map[string]interface{}{"error": err})
		return nil, err
	}

	body, attachments, err := ParseEmailBody(headerReader)
	if err != nil {
		logz.Error("Failed to parse email body", map[string]interface{}{"error": err})
		return nil, err
	}

	return NewEmail(mlHeader, body, attachments), nil
}

// ParseEmailHeader parses the email header and returns an IEmailHeader instance.
func ParseEmailHeader(header *mail.Header) (IEmailHeader, error) {
	from, _ := header.AddressList("From")
	to, _ := header.AddressList("To")
	cc, _ := header.AddressList("Cc")
	bcc, _ := header.AddressList("Bcc")
	replyTo, _ := header.AddressList("Reply-To")
	inReplyTo := header.Get("In-Reply-To")
	references := header.Get("References")
	subject := header.Get("Subject")
	date := header.Get("Date")
	messageID := header.Get("Message-ID")

	extraHeaders := &EmailExtraFields[string]{data: make(ExtraData[string])}

	for key, values := range header.Map() {
		if key != "From" && key != "To" && key != "Cc" && key != "Bcc" && key != "Reply-To" && key != "In-Reply-To" && key != "References" && key != "Subject" && key != "Date" && key != "Message-ID" {
			extraHeaders.data[key] = values[0]
		}
	}

	ttt := strings.Join(extractAddresses(to), ", ")
	ccc := strings.Join(extractAddresses(cc), ", ")
	bbb := strings.Join(extractAddresses(bcc), ", ")

	return &EmailHeader{
		From:         from[0].Address,
		To:           ttt,
		Cc:           ccc,
		Bcc:          bbb,
		Subject:      subject,
		Date:         date,
		MessageID:    messageID,
		InReplyTo:    inReplyTo,
		References:   references,
		ReplyTo:      replyTo[0].Address,
		ExtraHeaders: extraHeaders,
	}, nil
}

// extractAddresses extracts email addresses from a list of mail.Address and returns them as a slice of strings.
func extractAddresses(addresses []*mail.Address) []string {
	var result []string
	for _, address := range addresses {
		result = append(result, address.Address)
	}
	return result
}

// ParseEmailBody parses the email body and returns the body content and attachments.
func ParseEmailBody(reader *mail.Reader) (string, []Attachment, error) {
	var body string
	var attachments []Attachment

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", nil, err
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			bodyBytes, err := io.ReadAll(part.Body)
			if err != nil {
				return "", nil, err
			}
			body = string(bodyBytes)
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			contentType, _, _ := h.ContentType()
			attachmentBytes, err := io.ReadAll(part.Body)
			if err != nil {
				return "", nil, err
			}
			attachments = append(attachments, Attachment{
				Filename:    filename,
				ContentType: contentType,
				Data:        attachmentBytes,
			})
		}
	}

	return body, attachments, nil
}
