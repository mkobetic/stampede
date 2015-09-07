package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type MailDirectories []*MailDirectory

func (a MailDirectories) Len() int           { return len(a) }
func (a MailDirectories) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MailDirectories) Less(i, j int) bool { return a[i].Name < a[j].Name }

type MailDirectory struct {
	Path        string
	Name        string
	Directories map[string]*MailDirectory
	Folders     map[string]*MailFolder
}

func OpenDirectory(path string, info os.FileInfo, wg *sync.WaitGroup) *MailDirectory {
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
				directory := OpenDirectory(filepath.Join(path, info.Name()), info, wg)
				directories[directory.Name] = directory
			} else {
				folder := OpenFolder(directory, filepath.Join(path, info.Name()), info, wg)
				folders[folder.Name] = folder
			}
		}
	}
	directory.Folders = folders
	directory.Directories = directories
	return directory
}

func (d *MailDirectory) DirectoryList() (l MailDirectories) {
	for _, s := range d.Directories {
		l = append(l, s)
	}
	sort.Sort(l)
	return l
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

func (d *MailDirectory) Stats() *MessageStats {
	s := new(MessageStats)
	for _, sub := range d.Directories {
		s.Add(sub.Stats())
	}
	for _, f := range d.Folders {
		s.Add(f.Stats())
	}
	return s
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
