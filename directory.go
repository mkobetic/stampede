package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

type MailDirectory struct {
	Path        string
	Name        string
	Directories map[string]*MailDirectory
	Folders     map[string]*MailFolder
}

func OpenDirectory(path string, info os.FileInfo) *MailDirectory {
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()
	if info == nil {
		info, err = dir.Stat()
		if err != nil {
			log.Fatal(err)
		}
	}
	infos, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(info.Name())
	directory := &MailDirectory{Path: path, Name: info.Name()}
	directories := make(map[string]*MailDirectory)
	folders := make(map[string]*MailFolder)
	for _, info := range infos {
		if len(filepath.Ext(info.Name())) == 0 {
			if info.IsDir() {
				directory := OpenDirectory(filepath.Join(path, info.Name()), info)
				directories[directory.Name] = directory
			} else {
				folder := OpenFolder(directory, filepath.Join(path, info.Name()), info)
				folders[folder.Name] = folder
			}
		}
	}
	directory.Folders = folders
	directory.Directories = directories
	return directory
}

func (d *MailDirectory) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	DirectoryPage(w, d)
}

func (d *MailDirectory) UrlPath() string {
	return "/" + d.Name
}

func (d *MailDirectory) Label() string {
	return d.Name
}

func (d *MailDirectory) MessageStats() (count, size int64) {
	for _, s := range d.Directories {
		c, s := s.MessageStats()
		count += c
		size += s
	}
	for _, f := range d.Folders {
		count += int64(len(f.Messages))
	}
	return count, size
}

func (d *MailDirectory) Find(path []string) http.Handler {
	if len(path) == 0 {
		return d
	} else if f, found := d.Folders[path[0]]; found {
		return f.Find(path[1:])
	} else if s, found := d.Directories[path[0]]; found {
		return s.Find(path[1:])
	}
	return nil
}
