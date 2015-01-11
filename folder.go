package main

import (
	"bufio"
	"bytes"
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
	position := 0
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
			messagesById[message.Summary.Id] = message
		} else if !headerFinished {
			headerFinished = !message.scanHeaderLine(line)
		}
		position += len(line)
		for err == bufio.ErrBufferFull { // read the rest of the line
			line, err = reader.ReadSlice('\n')
			position += len(line)
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

func (f *MailFolder) UrlPath() string {
	return f.Directory.UrlPath() + "/" + f.Name
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
