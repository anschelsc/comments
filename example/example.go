package main

import (
	"fmt"
	"github.com/joshlf13/comments"
	"io"
	"os"
)

var (
	bash = []string{"#", "\n"}
	c1   = []string{"//", "\n"}
	c2   = []string{"/*", "*/"}
)

func main() {
	source := os.Stdin
	var err error
	if len(os.Args) > 1 {
		source, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("Could not open resource: %v\nUsing stdin\n", os.Args[1])
			source = os.Stdin
		}
	}

	delimeters := c2

	// Try messing around with start and end delimeters
	rdr := comments.NewCustomReader(source, delimeters[0], delimeters[1])
	io.Copy(os.Stdout, rdr)
}
