package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRestRouter() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", rootList).Methods("GET")
	r.HandleFunc("/{directory}", directoryList).Methods("GET")
	r.HandleFunc("/{directory}/{folder}", folderList).Methods("GET")
	return r
}

func directoryList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory, ok := Root.Directories[vars["directory"]]
	if !ok {
		http.NotFound(w, r)
		return
	}
	list := make([]string, 0)
	for key := range directory.Folders {
		list = append(list, key)
	}
	jsonResponse(w, list)
}

func rootList(w http.ResponseWriter, r *http.Request) {
	list := make([]string, 0, len(Root.Folders))
	for key := range Root.Folders {
		list = append(list, key)
	}
	jsonResponse(w, list)
}

func folderList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	directory, ok := Root.Directories[vars["directory"]]
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

func query(s string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		var err error
		if err != nil {
			log.Fatal(err)
		}
	}
}
