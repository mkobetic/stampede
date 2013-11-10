package main

import (
	"fmt"
	"os"
)

var MAILDIR string = os.ExpandEnv("$HOME/mail")
var Root []*MailDirectory = OpenRoot(MAILDIR)

func main() {
	for _, directory := range Root {
		for _, folder := range directory.Folders {
			fmt.Println(folder.Path, len(folder.Messages))
		}
	}
}
