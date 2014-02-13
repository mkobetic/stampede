package main

import ()

type MailMessage struct {
	Folder        *MailFolder
	start, length int
	header_length int
	From          string // mbox separator from part
	Date          string // mbox separator date part
	Summary       MessageSummary
}

type MessageSummary struct {
	Subject string
	Date    string
	From    string
	To      string
	Status  string
}

func (m *MailMessage) scanHeaderLine(line []byte) bool {
	m.header_length += len(line) // FIXME: if any header line exceeds buffer size header_length will be off
	if len(line) < 3 {           // should be either \n or \r\n
		return false
	}
	return true
}
