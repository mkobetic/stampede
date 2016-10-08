package main

import (
	"net/http"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	r := &MailDirectory{Directories: make(map[string]*MailDirectory)}
	a := &MailDirectory{Name: "A", Directories: make(map[string]*MailDirectory)}
	r.Directories["A"] = a
	ab := &MailDirectory{Name: "B", Folders: make(map[string]*MailFolder)}
	a.Directories["B"] = ab
	abc := &MailFolder{Name: "C", MessagesById: make(map[string]*MailMessage)}
	ab.Folders["C"] = abc
	m := &MailMessage{}
	abc.MessagesById["D"] = m

	tf := func(s string, dir http.Handler) {
		path := strings.Split(s, "")
		t.Run(strings.Join(path, ""), func(t *testing.T) {
			if got := r.Find(path); got != dir {
				t.Fatal(got)
			}
		})
	}

	tf("", r)
	tf("X", nil)
	tf("A", a)
	tf("B", nil)
	tf("AB", ab)
	tf("ABC", abc)
	tf("ABCD", m)
}
