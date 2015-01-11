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
	defer Root.Close()
	http.HandleFunc("/", Navigate)
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
