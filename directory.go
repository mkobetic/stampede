package main

import (
	"log"
	"os"
	"path/filepath"
)

type MailDirectory struct {
	Path        string
	Name        string
	Directories map[string]*MailDirectory
	Folders     map[string]*MailFolder
}

func OpenRoot(path string) map[string]*MailDirectory {
	root, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer root.Close()
	infos, err := root.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}
	directories := make(map[string]*MailDirectory)
	for _, info := range infos {
		if info.IsDir() {
			directory := OpenDirectory(filepath.Join(path, info.Name()), info)
			directories[directory.Name] = directory
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
