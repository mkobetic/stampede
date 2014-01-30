package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	root      = flag.String("root", "$HOME/mail", "directory containing the mail archive")
	Root      []*MailDirectory
	Templates = template.Must(template.New("all").ParseGlob("templates/*"))
)

func main() {
	flag.Parse()
	Root = OpenRoot(os.ExpandEnv(*root))
	log.Println("Templates: ")
	for _, t := range Templates.Templates() {
		log.Println(t.Name())
	}
	http.HandleFunc("/", render("stampede"))
	http.HandleFunc("/folder", render("folder"))
	http.HandleFunc("/message", render("message"))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func render(template string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := Templates.ExecuteTemplate(w, template, Root)
		if err != nil {
			log.Fatal(err)
		}
	}
}
