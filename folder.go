package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type MailFolder struct {
	Directory    *MailDirectory
	Path         string
	Name         string
	Messages     []*MailMessage
	MessagesById map[string]*MailMessage
}

func OpenFolder(directory *MailDirectory, path string, info os.FileInfo) *MailFolder {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	messages := make([]*MailMessage, 0, 20)
	messagesById := make(map[string]*MailMessage)
	position := int64(0)
	headerFinished := false
	folder := &MailFolder{Directory: directory, Path: path, Name: info.Name()}
	message := &MailMessage{Folder: folder}
	for {
		line, err := reader.ReadSlice('\n')
		if len(line) == 0 {
			break
		}
		if startsNewMessage(line) {
			message.length = position - message.start
			message = &MailMessage{Folder: folder, start: position}
			headerFinished = false
			messages = append(messages, message)
		} else if !headerFinished {
			headerFinished = !message.scanHeaderLine(line)
		} else if headerFinished {
			messagesById[message.Summary.Id] = message
		}
		position += int64(len(line))
		for err == bufio.ErrBufferFull { // read the rest of the line
			line, err = reader.ReadSlice('\n')
			position += int64(len(line))
		}
	}
	if int64(position) != info.Size() {
		log.Fatalf("Folder %s length mismatch (%d != %d)!", folder.Name, position, info.Size())
	}
	log.Println("\t", info.Name(), len(messages))
	folder.Messages = messages
	folder.MessagesById = messagesById
	return folder
}

func (f *MailFolder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	FolderPage(w, f)
}

func (f *MailFolder) Find(path []string) http.Handler {
	if len(path) == 0 {
		return f
	} else if m, found := f.MessagesById[path[0]]; found {
		return m
	}
	return nil
}

func (f *MailFolder) UrlPath() string {
	return f.Directory.UrlPath() + "/" + f.Name
}

func (f *MailFolder) cClass() string {
	if f.Stats().Unread.Count > 0 {
		return "folder unread"
	} else {
		return "folder"
	}
}

func (f *MailFolder) Label() string {
	u := f.Stats().Unread.Count
	if u > 0 {
		return fmt.Sprintf("%s (%d)", f.Name, u)
	} else {
		return f.Name
	}
}

func (f *MailFolder) Stats() *MessageStats {
	s := new(MessageStats)
	for _, m := range f.Messages {
		s.Count(m)
	}
	return s
}

type MessageStats struct {
	Read    Total
	Unread  Total
	Deleted Total
}

func (s *MessageStats) Add(s2 *MessageStats) {
	(&s.Read).Add(s2.Read)
	(&s.Unread).Add(s2.Unread)
	(&s.Deleted).Add(s2.Deleted)
}

func (s *MessageStats) Count(m *MailMessage) {
	switch m.cClass() {
	case "msg expunged":
		(&s.Deleted).AddSize(m.length)
	case "msg unread":
		(&s.Unread).AddSize(m.length)
	default:
		(&s.Read).AddSize(m.length)
	}
}

func (s *MessageStats) Total() Total {
	return Total{
		Count: s.Read.Count + s.Unread.Count + s.Deleted.Count,
		Size:  s.Read.Size + s.Unread.Size + s.Deleted.Size,
	}
}

type Total struct {
	Count int64
	Size  int64
}

func (t *Total) Add(u Total) {
	t.Count += u.Count
	t.Size += u.Size
}

func (t *Total) AddSize(s int64) {
	t.Count += 1
	t.Size += s
}

var (
	FROM_PREFIX    = []byte("From ")
	FROM_TIMESTAMP = "Mon Jan 2 15:04:05 2006"
)

func startsNewMessage(line []byte) bool {
	if !bytes.HasPrefix(line, FROM_PREFIX) {
		return false
	}
	line = line[len(FROM_PREFIX):]
	tsStart := bytes.IndexByte(line, ' ')
	if tsStart == -1 {
		return false
	}
	line = bytes.Trim(line[tsStart:len(line)-1], " \t\n")
	if len(line) > len(FROM_TIMESTAMP)+2 {
		return false
	}
	_, err := time.Parse(FROM_TIMESTAMP, string(line))
	return err == nil
}
