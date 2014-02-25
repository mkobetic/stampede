package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"time"
)

type MailFolder struct {
	Directory *MailDirectory
	Path      string
	Name      string
	Messages  []*MailMessage
}

func OpenFolder(directory *MailDirectory, path string, info os.FileInfo) *MailFolder {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	messages := make([]*MailMessage, 0, 20)
	position := 0
	headerFinished := false
	folder := &MailFolder{directory, path, info.Name(), nil}
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
		}
		position += len(line)
		for err == bufio.ErrBufferFull { // read the rest of the line
			line, err = reader.ReadSlice('\n')
			position += len(line)
		}
	}
	if int64(position) != info.Size() {
		log.Fatal("Folder length mismatch!")
	}
	log.Println("\t", info.Name(), len(messages))
	folder.Messages = messages
	return folder
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
