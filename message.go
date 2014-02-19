package main

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
	"time"
)

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
	Date    time.Time
	From    string
	To      string
	Status  string
}

func (m *MailMessage) scanHeaderLine(line []byte) bool {
	// FIXME: if any header line exceeds buffer size header_length will be off
	m.header_length += len(line)
	if len(line) < 3 || len(bytes.TrimSpace(line)) == 0 {
		return false
	}
	switch {
	case line[0] == 'S':
		m.scanHeaderLineS(line)
	case line[0] == 'F':
		m.scanHeaderLineF(line)
	case line[0] == 'D':
		m.scanHeaderLineD(line)
	case line[0] == 'X':
		m.scanHeaderLineX(line)
	}
	return true
}

func (m *MailMessage) scanHeaderLineS(line []byte) {
	value := bytes.TrimPrefix(line, []byte("Subject: "))
	if len(line) > len(value) {
		m.Summary.Subject = decodeString(value)
	}
}

func (m *MailMessage) scanHeaderLineF(line []byte) {
	value := bytes.TrimPrefix(line, []byte("From: "))
	if len(line) > len(value) {
		m.Summary.From = decodeString(value)
	}
}

func (m *MailMessage) scanHeaderLineD(line []byte) {
	value := bytes.TrimPrefix(line, []byte("Date: "))
	if len(line) > len(value) {
		m.Summary.Date, _ = decodeDate(value)
	}
}

func (m *MailMessage) scanHeaderLineX(line []byte) {
}

func decodeString(line []byte) string {
	return string(bytes.TrimRight(line, "\r\n"))
}

var (
	dateParser = regexp.MustCompile(`((Mon|Tue|Wed|Thu|Fri|Sat|Sun), )?(\d{1,2}) (Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) (\d{4}) (\d{1,2}):(\d\d):(\d\d) ([+-]\d{4}|Z)( \(?([A-Z]{3})\)?)?`)
	months     = map[string]time.Month{"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4, "May": 5, "Jun": 6, "Jul": 7, "Aug": 8, "Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12}
)

func decodeDate(line []byte) (time.Time, error) {
	fields := dateParser.FindStringSubmatch(string(line))
	day, err := strconv.Atoi(fields[3])
	if err != nil {
		return time.Time{}, errors.New("Invalid Day of Month")
	}
	month, present := months[fields[4]]
	if !present {
		return time.Time{}, errors.New("Invalid Month")
	}
	year, err := strconv.Atoi(fields[5])
	if err != nil || year < 1950 || year > 2050 {
		return time.Time{}, errors.New("Invalid Year")
	}
	hour, err := strconv.Atoi(fields[6])
	if err != nil || hour < 0 || hour > 23 {
		return time.Time{}, errors.New("Invalid Hour")
	}
	min, err := strconv.Atoi(fields[7])
	if err != nil || min < 0 || min > 59 {
		return time.Time{}, errors.New("Invalid Hour")
	}
	sec, err := strconv.Atoi(fields[8])
	if err != nil || sec < 0 || sec > 60 {
		return time.Time{}, errors.New("Invalid Second")
	}
	loc, err := decodeLocation(string(fields[9]))
	if err != nil {
		return time.Time{}, errors.New("Invalid Timezone Offset")
	}
	return time.Date(year, month, day, hour, min, sec, 0, loc), nil
}

func decodeLocation(input string) (*time.Location, error) {
	hours, err := strconv.Atoi(input[1:3])
	if err != nil || hours < 0 || hours > 12 {
		return nil, errors.New("Invalid Timezone Offset Hour")
	}
	minutes, err := strconv.Atoi(input[3:5])
	if err != nil || minutes < 0 || minutes > 59 {
		return nil, errors.New("Invalid Timezone Offset Minute")
	}
	offset := hours*3600 + minutes*60
	if input[0] == '-' {
		offset = -offset
	}
	return time.FixedZone(input, offset), nil
}
