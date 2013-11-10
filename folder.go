package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"path/filepath"
)

var FROM_PREFIX = []byte("From ")

type MailMessage struct{}

type MailFolder struct {
	Path     string
	Messages []*MailMessage
}

type MailDirectory struct {
	Path    string
	Folders []*MailFolder
}

func OpenRoot(path string) []*MailDirectory {
	root, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer root.Close()
	infos, err := root.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	directories := make([]*MailDirectory, 0, 10)
	for _, info := range infos {
		if info.IsDir() {
			directories = append(directories, OpenDirectory(filepath.Join(MAILDIR, info.Name()), info))
		}
	}
	return directories
}

func OpenDirectory(path string, info os.FileInfo) *MailDirectory {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	infos, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	folders := make([]*MailFolder, 0, 20)
	for _, info := range infos {
		if !(len(filepath.Ext(info.Name())) > 0) {
			folders = append(folders, OpenFolder(filepath.Join(path, info.Name()), info))
		}
	}
	return &MailDirectory{path, folders}
}

func OpenFolder(path string, info os.FileInfo) *MailFolder {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	messages := make([]*MailMessage, 20)
	var message *MailMessage
	for {
		line, err := reader.ReadSlice('\n')
		if len(line) == 0 {
			break
		}
		if bytes.HasPrefix(line, FROM_PREFIX) {
			message = &MailMessage{}
			messages = append(messages, message)
		}
		for err == bufio.ErrBufferFull { // read the rest of the line
			line, err = reader.ReadSlice('\n')
		}
	}
	return &MailFolder{path, messages}
}
