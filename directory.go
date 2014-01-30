package main

import (
	"log"
	"os"
	"path/filepath"
)

type MailDirectory struct {
	Path    string
	Name    string
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
			directories = append(directories, OpenDirectory(filepath.Join(path, info.Name()), info))
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
	log.Println(info.Name())
	directory := &MailDirectory{path, info.Name(), nil}
	folders := make([]*MailFolder, 0, 20)
	for _, info := range infos {
		if !(len(filepath.Ext(info.Name())) > 0) {
			folders = append(folders, OpenFolder(directory, filepath.Join(path, info.Name()), info))
		}
	}
	directory.Folders = folders
	return directory
}
