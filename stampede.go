package main

import (
	"encoding/json"
	"flag"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

var (
	root      = flag.String("root", "$HOME/mail", "directory containing the mail archive")
	Root      map[string]*MailDirectory
	Templates = template.Must(template.New("all").ParseGlob("*.template"))
)

func main() {
	flag.Parse()
	Root = OpenRoot(os.ExpandEnv(*root))
	log.Println("Templates: ")
	for _, t := range Templates.Templates() {
		log.Println(t.Name())
	}
	//	http.HandleFunc("/", render("stampede.template"))
	r := mux.NewRouter()
	r.HandleFunc("/", rootList).Methods("GET")
	r.HandleFunc("/{directory}", directoryList).Methods("GET")
	r.HandleFunc("/{directory}/{folder}", folderList).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func directoryList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory, ok := Root[vars["directory"]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	list := make([]string, 0, len(directory.Folders))
	for key := range directory.Folders {
		list = append(list, key)
	}
	jsonResponse(w, list)
}

func rootList(w http.ResponseWriter, r *http.Request) {
	list := make([]string, 0, len(Root))
	for key := range Root {
		list = append(list, key)
	}
	jsonResponse(w, list)
}

func folderList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory, ok := Root[vars["directory"]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	folder, ok := directory.Folders[vars["folder"]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	list := make([]*MessageSummary, 0, len(folder.Messages))
	for _, message := range folder.Messages {
		list = append(list, &message.Summary)
	}
	jsonResponse(w, list)
}

func jsonResponse(w http.ResponseWriter, r interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(r)
	if err != nil {
		log.Fatal(err)
	}
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

func query(s string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var err error
		if err != nil {
			log.Fatal(err)
		}
	}
}
