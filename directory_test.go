package main

import "testing"

func TestFind(t *testing.T) {
	r := &MailDirectory{Directories: make(map[string]*MailDirectory)}

	if r.Find([]string{}) != r {
		t.Fail()
	}

	if r.Find([]string{"X"}) != nil {
		t.Fail()
	}

	a := &MailDirectory{Name: "A", Directories: make(map[string]*MailDirectory)}
	r.Directories["A"] = a
	ab := &MailDirectory{Name: "B", Folders: make(map[string]*MailFolder)}
	a.Directories["B"] = ab

	abc := &MailFolder{Name: "C"}
	ab.Folders["C"] = abc

	if r.Find([]string{"A"}) != a {
		t.Fail()
	}

	if r.Find([]string{"B"}) != nil {
		t.Fail()
	}

	if r.Find([]string{"A", "B"}) != ab {
		t.Fail()
	}

	if r.Find([]string{"A", "B", "C"}) != abc {
		t.Fail()
	}

	if r.Find([]string{"A", "B", "C", "D"}) != abc {
		t.Fail()
	}
}
