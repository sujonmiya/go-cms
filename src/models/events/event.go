package events

type Event uint8

const (
	ThemeChanged Event = iota
	DatabaseNameChanged
	MailServerHostChanged
	MailServerPortChanged
	MailServerLoginChanged
	MailServerPasswordChanged
)
