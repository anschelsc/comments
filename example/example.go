package main

import (
	"fmt"
	"github.com/joshlf13/comments"
	"io"
	"os"
)

var (
	bash = []string{"#", "\n", ""}
	c1   = []string{"//", "\n", ""}
	c2   = []string{"/*", "*/", ""}
)

func main() {
	source := os.Stdin
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [<filename>]; falling back to stdin\n", os.Args[0])
	} else {
		var err error
		source, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Error openning file: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

	delim := c2

	// Try messing around with start and end delimeters
	rdr := comments.NewCustomReader(source, delim[0], delim[1], delim[2])
	io.Copy(os.Stdout, rdr)
}
