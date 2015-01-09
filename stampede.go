package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	root = flag.String("root", "$HOME/mail", "directory containing the mail archive")
	Root *MailDirectory
)

func main() {
	flag.Parse()
	Root = OpenDirectory(os.ExpandEnv(*root), nil)
	http.HandleFunc("/", Directory)
	http.HandleFunc("/directory/", Directory)
	//	http.HandleFunc("/folder", func(w http.ResponseWriter, r *http.Request) { Folder(w, Root) })
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Directory(w http.ResponseWriter, r *http.Request) {
	dn := strings.TrimPrefix(r.URL.Path, "/directory/")
	d, found := Root.Directories[dn]
	if !found {
		d = Root
	}
	DirectoryPage(w, d)
}
