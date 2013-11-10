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
	total := 0
	for _, directory := range Root {
		for _, folder := range directory.Folders {
			total = total + len(folder.Messages)
			fmt.Println(folder.Path, len(folder.Messages))
		}
	}
	fmt.Println("Total: ", total)
}
