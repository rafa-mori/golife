package email

import (
	"bytes"
	"github.com/emersion/go-imap"
	"io"
)

// EmailBodyPart represents a part of the email body.
type EmailBodyPart interface {
	GetContentType() string
	GetContent() (string, error)
}

// emailBodyPart is a concrete implementation of EmailBodyPart.
type emailBodyPart struct {
	contentType string
	content     string
}

func (ebp *emailBodyPart) GetContentType() string {
	return ebp.contentType
}

func (ebp *emailBodyPart) GetContent() (string, error) {
	return ebp.content, nil
}

// newEmailBodyPart creates a new emailBodyPart from an imap.BodyStructure and its content.
func newEmailBodyPart(part *imap.BodyStructure, body imap.Literal) (EmailBodyPart, error) {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, body)
	if err != nil {
		return nil, err
	}

	contentType := part.MIMEType + "/" + part.MIMESubType
	content := buffer.String()

	return &emailBodyPart{
		contentType: contentType,
		content:     content,
	}, nil
}
