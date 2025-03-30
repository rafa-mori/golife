package email

// IEmailHeader defines the interface for email headers.
type IEmailHeader interface {
	GetMessageID() string
	GetFrom() string
	GetTo() string
	GetSubject() string
	GetDate() string
	GetCc() string
	GetBcc() string
	GetReplyTo() string
	GetInReplyTo() string
	GetReferences() string
	GetExtraHeaders() ExtraFields[string]

	SetMessageID(string)
	SetFrom(string)
	SetTo(string)
	SetSubject(string)
	SetDate(string)
	SetCc(string)
	SetBcc(string)
	SetReplyTo(string)
	SetInReplyTo(string)
	SetReferences(string)
	SetExtraHeaders(ExtraFields[string])
}

// EmailHeader represents the email header.
type EmailHeader struct {
	MessageID    string
	From         string
	To           string
	Subject      string
	Date         string
	Cc           string
	Bcc          string
	ReplyTo      string
	InReplyTo    string
	References   string
	ExtraHeaders ExtraFields[string]
}

func (eh *EmailHeader) GetMessageID() string  { return eh.MessageID }
func (eh *EmailHeader) GetFrom() string       { return eh.From }
func (eh *EmailHeader) GetTo() string         { return eh.To }
func (eh *EmailHeader) GetSubject() string    { return eh.Subject }
func (eh *EmailHeader) GetDate() string       { return eh.Date }
func (eh *EmailHeader) GetCc() string         { return eh.Cc }
func (eh *EmailHeader) GetBcc() string        { return eh.Bcc }
func (eh *EmailHeader) GetReplyTo() string    { return eh.ReplyTo }
func (eh *EmailHeader) GetInReplyTo() string  { return eh.InReplyTo }
func (eh *EmailHeader) GetReferences() string { return eh.References }
func (eh *EmailHeader) GetExtraHeaders() ExtraFields[string] {
	return eh.ExtraHeaders
}

func (eh *EmailHeader) SetMessageID(id string)          { eh.MessageID = id }
func (eh *EmailHeader) SetFrom(from string)             { eh.From = from }
func (eh *EmailHeader) SetTo(to string)                 { eh.To = to }
func (eh *EmailHeader) SetSubject(subject string)       { eh.Subject = subject }
func (eh *EmailHeader) SetDate(date string)             { eh.Date = date }
func (eh *EmailHeader) SetCc(cc string)                 { eh.Cc = cc }
func (eh *EmailHeader) SetBcc(bcc string)               { eh.Bcc = bcc }
func (eh *EmailHeader) SetReplyTo(replyTo string)       { eh.ReplyTo = replyTo }
func (eh *EmailHeader) SetInReplyTo(inReplyTo string)   { eh.InReplyTo = inReplyTo }
func (eh *EmailHeader) SetReferences(references string) { eh.References = references }
func (eh *EmailHeader) SetExtraHeaders(extraHeaders ExtraFields[string]) {
	eh.ExtraHeaders = extraHeaders
}

// NewEmailHeader creates a new EmailHeader instance with default values.
func NewEmailHeader() IEmailHeader {
	return &EmailHeader{
		MessageID:    "",
		From:         "",
		To:           "",
		Subject:      "",
		Date:         "",
		Cc:           "",
		Bcc:          "",
		ReplyTo:      "",
		InReplyTo:    "",
		References:   "",
		ExtraHeaders: &EmailExtraFields[string]{data: make(ExtraData[string])},
	}
}

// NewEmailHeaderWithData creates a new EmailHeader instance with the provided extra data.
func NewEmailHeaderWithData(data ExtraData[string]) IEmailHeader {
	return &EmailHeader{
		MessageID:    "",
		From:         "",
		To:           "",
		Subject:      "",
		Date:         "",
		Cc:           "",
		Bcc:          "",
		ReplyTo:      "",
		InReplyTo:    "",
		References:   "",
		ExtraHeaders: &EmailExtraFields[string]{data: data},
	}
}

// NewEmailHeaderWithAllData creates a new EmailHeader instance with all provided data.
func NewEmailHeaderWithAllData(messageID, from, to, subject, date, cc, bcc, replyTo, inReplyTo, references string) IEmailHeader {
	return &EmailHeader{
		MessageID:    messageID,
		From:         from,
		To:           to,
		Subject:      subject,
		Date:         date,
		Cc:           cc,
		Bcc:          bcc,
		ReplyTo:      replyTo,
		InReplyTo:    inReplyTo,
		References:   references,
		ExtraHeaders: &EmailExtraFields[string]{data: make(ExtraData[string])},
	}
}
