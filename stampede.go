package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	root = flag.String("root", "$HOME/mail", "directory containing the mail archive")
	Root []*MailDirectory
)

func main() {
	flag.Parse()
	Root = OpenRoot(os.ExpandEnv(*root))
	for _, directory := range Root {
		for _, folder := range directory.Folders {
			fmt.Println(folder.Path, len(folder.Messages))
		}
	}
}
