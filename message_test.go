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
	if m.scanHeaderLine([]byte(" \n")) {
		t.Fatal()
	}
	if m.scanHeaderLine([]byte(" 	 	 	\n")) {
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
	if !m.scanHeaderLine([]byte("From: John Doe\n")) {
		t.Fatal()
	}
	if string(m.Summary.From) != "John Doe" {
		t.Fatal()
	}
}

func TestDateParser(t *testing.T) {
	for _, date := range []string{
		"Mon, 9 Feb 2009 12:21:36 -0800 (PST)",
		"9 Feb 2009 12:21:36 -0800 (PST)",
		"Mon, 9 Feb 2009 12:21:36 -0800",
		"9 Feb 2009 12:21:36 -0800",
		"09 Feb 2009 12:21:36 -0800",
		"9 Feb 2009 12:21:36 +0000",
		"9 Feb 2009 2:21:36 -0800",
	} {
		d := dateParser.FindSubmatch([]byte(date))
		if len(d) == 0 {
			t.Fatal(date)
		}
	}
}

func TestDecodeDate(t *testing.T) {
	for _, date := range []string{
		"Mon, 9 Feb 2009 12:21:36 -0800 (PST)",
		"9 Feb 2009 12:21:36 -0800 (PST)",
		"Mon, 9 Feb 2009 12:21:36 -0800",
		"9 Feb 2009 12:21:36 -0800",
		"09 Feb 2009 12:21:36 -0800",
		"9 Feb 2009 12:21:36 +0000",
		"9 Feb 2009 2:21:36 -0800",
	} {
		_, err := decodeDate([]byte(date))
		if err != nil {
			t.Fatal(date)
		}
	}
}

func TestDecodeLocation(t *testing.T) {
	for _, input := range []string{
		"-0800",
		"+0830",
		"-0000",
		"+0000",
		"-1200",
		"+1205",
	} {
		_, err := decodeLocation(input)
		if err != nil {
			t.Fatal(input)
		}
	}
}
