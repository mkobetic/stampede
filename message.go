package main

import ()

type MailMessage struct {
	Folder        *MailFolder
	start, length int
	header_length int
}

func (m *MailMessage) scanHeaderLine(line []byte) bool {
	m.header_length += len(line) // FIXME: if any header line exceeds buffer size header_length will be off
	return len(line) < 3         // should be either \n or \r\n, no valid header line should be this short
}
