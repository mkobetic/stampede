package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"sync"
)

type MailFolders []*MailFolder

func (a MailFolders) Len() int           { return len(a) }
func (a MailFolders) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MailFolders) Less(i, j int) bool { return a[i].Name < a[j].Name }

type MailFolder struct {
	Directory    *MailDirectory
	Path         string
	Name         string
	Messages     MessagesByDate
	MessagesById map[string]*MailMessage
}

func OpenFolder(directory *MailDirectory, path string, info os.FileInfo, wg *sync.WaitGroup) *MailFolder {
	folder := &MailFolder{Directory: directory, Path: path, Name: info.Name()}
	wg.Add(1)
	go openFolder(folder, info, wg)
	return folder
}

func openFolder(folder *MailFolder, info os.FileInfo, wg *sync.WaitGroup) {
	file, err := os.Open(folder.Path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	messages := make(MessagesByDate, 0, 20)
	messagesById := make(map[string]*MailMessage)
	position := int64(0)
	headerFinished := false
	message := &MailMessage{Folder: folder}
	for {
		line, err := reader.ReadSlice('\n')
		if len(line) == 0 {
			break
		}
		if startsNewMessage(line) {
			message.length = position - message.start
			messagesById[message.Id()] = message
			message = &MailMessage{Folder: folder, start: position}
			headerFinished = false
			messages = append(messages, message)
		} else if !headerFinished {
			headerFinished = !message.scanHeaderLine(line)
		}
		position += int64(len(line))
		for err == bufio.ErrBufferFull { // read the rest of the line
			line, err = reader.ReadSlice('\n')
			position += int64(len(line))
		}
	}
	message.length = position - message.start
	messagesById[message.Id()] = message

	if int64(position) != info.Size() {
		log.Fatalf("Folder %s length mismatch (%d != %d)!", folder.Name, position, info.Size())
	}

	sort.Sort(sort.Reverse(messages))
	log.Println(folder.UrlPath(), len(messages))
	folder.Messages = messages
	folder.MessagesById = messagesById
	wg.Done()
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

func (f *MailFolder) DumpOffsets() error {
	out, err := os.Create(f.Path + ".offsets2")
	if err != nil {
		log.Printf("dumping %s: %s", f.Path, err)
		return err
	}
	defer out.Close()
	byOffset := make(MessagesByOffset, len(f.Messages))
	copy(byOffset, f.Messages)
	sort.Sort(byOffset)
	for _, m := range byOffset {
		fmt.Fprintf(out, "%d,%d\n", m.start, m.length)
	}
	log.Printf("dumping %s: %d messages", f.Path, len(byOffset))
	return nil
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
	FROM_PREFIX          = []byte("From ")
	FROM_TIMESTAMP       = "Mon Jan 2 15:04:05 2006"
	FROM_TIMESTAMP_REGEX = regexp.MustCompile(`^` +
		`(Mon|Tue|Wed|Thu|Fri|Sat|Sun) ` +
		`(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec) [ \d]?\d ` +
		`[ \d]?\d:\d\d(:\d\d)? \d{4}$`)
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
	line = bytes.Trim(line[tsStart:len(line)-1], " \t\r\n")
	return FROM_TIMESTAMP_REGEX.Match(line)
}
