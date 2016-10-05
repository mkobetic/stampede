package main

import (
	"bytes"
	"encoding/base32"
	"errors"
	"html"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
)

type MessagesByOffset []*MailMessage

func (a MessagesByOffset) Len() int           { return len(a) }
func (a MessagesByOffset) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MessagesByOffset) Less(i, j int) bool { return a[i].start < a[j].start }

type MailMessages []*MailMessage

func (a MailMessages) Len() int           { return len(a) }
func (a MailMessages) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MailMessages) Less(i, j int) bool { return a[i].Summary.Date.Before(a[j].Summary.Date) }

type MailMessage struct {
	Folder        *MailFolder
	start, length int64
	header_length int64
	From          string // mbox separator from part
	Date          string // mbox separator date part
	Summary       MessageSummary
}

type MessageSummary struct {
	Id      string
	Subject string
	Date    time.Time
	From    string
	To      string
	Status  MozillaStatus
}

type MozillaStatus int

// Following definitions were copied from
// http://www.eyrich-net.org/mozilla/X-Mozilla-Status.html?en
const (
	// Message has been read.
	MSG_FLAG_READ MozillaStatus = 0x0001
	// A reply has been successfully sent.
	MSG_FLAG_REPLIED = 0x0002
	// The user has flagged this message.
	MSG_FLAG_MARKED = 0x0004
	// Already gone (when folder not compacted). Since actually removing a message
	// from a folder is a semi-expensive operation, we tend to delay it;
	// messages with this bit set will be removed the next time folder compaction is done.
	// Once this bit is set, it never gets un-set.
	MSG_FLAG_EXPUNGED = 0x0008
	// Whether subject has “Re:” on the front. The folder summary uniquifies all of the strings in it,
	// and to help this, any string which begins with “Re:” has that stripped first.
	// This bit is then set, so that when presenting the message, we know to put it back
	// (since the “Re:” is not itself stored in the file).
	MSG_FLAG_HAS_RE = 0x0010
	// Whether the children of this sub-thread are folded in the display.
	MSG_FLAG_ELIDED = 0x0020
	// DB has offline news or imap article.
	MSG_FLAG_OFFLINE = 0x0080
	// If set, this thread is watched.
	MSG_FLAG_WATCHED = 0x0100
	// If set, then this message's sender has been authenticated when sending this msg.
	// This means the POP3 server gave a positive answer to the XSENDER command.
	// Since this command is no standard and only known by few servers, this flag is unmeaning in most cases.
	MSG_FLAG_SENDER_AUTHED = 0x0200
	// If set, then this message's body contains not the whole message, and a link is available
	// in the message to download the rest of it from the POP server.
	// This can be only a few lines of the message (in case of size restriction for the download of messages)
	// or nothing at all (in case of “Fetch headers only”)
	MSG_FLAG_PARTIAL = 0x0400
	// If set, this message is queued for delivery. This only ever gets set on messages in the queue folder,
	// but is used to protect against the case of other messages having made their way in there somehow
	// – if some other program put a message in the queue, it won't be delivered later!
	MSG_FLAG_QUEUED = 0x0800
	// This message has been forwarded.
	MSG_FLAG_FORWARDED = 0x1000
	//These are used to remember the message priority in interal status flags.
	MSG_FLAG_PRIORITIES = 0xE000
)

func (s MozillaStatus) cClass() string {
	switch {
	case s&MSG_FLAG_EXPUNGED != 0:
		return "msg expunged"
	case s&MSG_FLAG_READ != 0:
		return "msg read"
	default:
		return "msg unread"
	}
}

func (m *MailMessage) UrlPath() string {
	return m.Folder.UrlPath() + "/" + m.Summary.Id
}

func (m *MailMessage) Id() string {
	if m.Summary.Id == "" {
		b := make([]byte, 15)
		rand.New(rand.NewSource(m.Summary.Date.Unix())).Read(b)
		m.Summary.Id = base32.StdEncoding.EncodeToString(b)
	}
	return m.Summary.Id
}

func (m *MailMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(m.Folder.Path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
	defer file.Close()
	if _, err := file.Seek(m.start, os.SEEK_SET); err != nil {
		http.Error(w, err.Error(), 500)
	}
	b := io.LimitReader(file, m.length)
	MessagePage(w, m, b)
}

// Message rendering.

func (m *MailMessage) hSent() string {
	return m.Summary.Date.Format("06/01/02 03:04:05")
}

func (m *MailMessage) hSubject() string {
	return html.EscapeString(m.Summary.Subject)
}

func (m *MailMessage) hSender() string {
	return html.EscapeString(m.Summary.From)
}

func (m *MailMessage) cClass() string {
	return m.Summary.Status.cClass()
}

// Message scanning.

func (m *MailMessage) scanHeaderLine(line []byte) bool {
	// FIXME: if any header line exceeds buffer size header_length will be off
	m.header_length += int64(len(line))
	if len(line) < 3 || len(bytes.TrimSpace(line)) == 0 {
		return false
	}
	switch line[0] {
	case 'S':
		m.scanHeaderLineS(line)
	case 'F':
		m.scanHeaderLineF(line)
	case 'D':
		m.scanHeaderLineD(line)
	case 'X':
		m.scanHeaderLineX(line)
	case 'M':
		m.scanHeaderLineM(line)
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

func (m *MailMessage) scanHeaderLineM(line []byte) {
	// Process Message-ID/Message-Id field.
	value := bytes.TrimPrefix(line, []byte("Message-I"))
	if len(line) == len(value) {
		return
	}
	value2 := bytes.TrimPrefix(value, []byte("d: "))
	if len(value) > len(value2) {
		m.Summary.Id = decodeString(value2)
	} else if value2 = bytes.TrimPrefix(value, []byte("D: ")); len(value) > len(value2) {
		m.Summary.Id = decodeString(value2)
	}
}

func (m *MailMessage) scanHeaderLineX(line []byte) {
	value := bytes.TrimPrefix(line, []byte("X-mozilla-status: "))
	if len(line) > len(value) {
		i, _ := strconv.ParseInt(string(value[:4]), 16, 0)
		m.Summary.Status = MozillaStatus(i)
	} else if value = bytes.TrimPrefix(line, []byte("X-Mozilla-Status: ")); len(line) > len(value) {
		i, _ := strconv.ParseInt(string(value[:4]), 16, 0)
		m.Summary.Status = MozillaStatus(i)
	}
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
	if fields == nil {
		return time.Time{}, errors.New("Invalid Date Format: " + string(line))
	}
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
