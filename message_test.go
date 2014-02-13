package main

import (
	"testing"
)

func TestMailMessageScanHeaderLine(t *testing.T) {
	var m *MailMessage = &MailMessage{}
	if m.scanHeaderLine([]byte("\n")) {
		t.Fatal()
	}
	if m.scanHeaderLine([]byte("\r\n")) {
		t.Fatal()
	}
	if m.scanHeaderLine([]byte("")) {
		t.Fatal()
	}
	if !m.scanHeaderLine([]byte("Subject: Hello World\n")) {
		t.Fatal()
	}
	if string(m.Summary.Subject) != "Hello World" {
		t.Fatal()
	}
}
