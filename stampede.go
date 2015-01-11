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
	http.HandleFunc("/", Navigate)
	// http.HandleFunc("/directory/", Directory)
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

func Navigate(w http.ResponseWriter, r *http.Request) {
	var path []string
	if r.URL.Path != "/" {
		path = strings.Split(r.URL.Path, "/")
		if len(path) > 0 && len(path[0]) == 0 {
			path = path[1:]
		}
	}
	log.Print(r.URL.Path, len(path), path)
	if d := Root.Find(path); d != nil {
		d.ServeHTTP(w, r)
	} else {
		http.Error(w, "Not Found", 404)
	}
}
