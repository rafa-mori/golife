package email

// IEmail defines the interface for an email.
type IEmail interface {
	GetHeader() IEmailHeader
	GetBody() string
	GetAttachments() []Attachment
	SetHeader(IEmailHeader)
	SetBody(string)
	SetAttachments([]Attachment)
}

// Email represents an email with a header, body, and attachments.
type Email struct {
	Header      IEmailHeader
	Body        string
	Attachments []Attachment
}

func (e *Email) GetHeader() IEmailHeader                 { return e.Header }
func (e *Email) GetBody() string                         { return e.Body }
func (e *Email) GetAttachments() []Attachment            { return e.Attachments }
func (e *Email) SetHeader(header IEmailHeader)           { e.Header = header }
func (e *Email) SetBody(body string)                     { e.Body = body }
func (e *Email) SetAttachments(attachments []Attachment) { e.Attachments = attachments }

// NewEmail creates a new Email instance with the provided header, body, and attachments.
func NewEmail(header IEmailHeader, body string, attachments []Attachment) *Email {
	return &Email{
		Header:      header,
		Body:        body,
		Attachments: attachments,
	}
}

// IAttachment defines the interface for an email attachment.
type IAttachment interface {
	GetFilename() string
	GetContentType() string
	GetData() []byte

	SetFilename(string)
	SetContentType(string)
	SetData([]byte)
}

// Attachment represents an email attachment.
type Attachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

func (a *Attachment) GetFilename() string               { return a.Filename }
func (a *Attachment) GetContentType() string            { return a.ContentType }
func (a *Attachment) GetData() []byte                   { return a.Data }
func (a *Attachment) SetFilename(filename string)       { a.Filename = filename }
func (a *Attachment) SetContentType(contentType string) { a.ContentType = contentType }
func (a *Attachment) SetData(data []byte)               { a.Data = data }

// NewAttachment creates a new Attachment instance with the provided filename, content type, and data.
func NewAttachment(filename, contentType string, data []byte) *Attachment {
	return &Attachment{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
}
