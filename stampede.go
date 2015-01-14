package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	root   = flag.String("root", "$HOME/mail", "directory containing the mail archive")
	assets = flag.String("assets", "assets", "directory containing asset files (css, ...)")
	Root   *MailDirectory
)

func main() {
	flag.Parse()
	Root = OpenDirectory(os.ExpandEnv(*root), nil)
	http.HandleFunc("/", Navigate)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(*assets))))
	log.Print("Listening at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func Navigate(w http.ResponseWriter, r *http.Request) {
	var path []string
	// URL.Path always starts with a /
	path = strings.Split(r.URL.Path[1:], "/")
	if len(path) > 0 && len(path[0]) == 0 {
		path = path[1:]
	}
	log.Printf("Navigate %#v", path)
	if d := Root.Find(path); d != nil {
		d.ServeHTTP(w, r)
	} else {
		http.Error(w, "Not Found", 404)
	}
}
