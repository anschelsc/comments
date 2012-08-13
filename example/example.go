package main

import (
	"github.com/joshlf13/comments"
	"io"
	"os"
)


func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %v <start> <stop> [<filename>]\n", os.Args[0])
		os.Exit(1)
	}
	
	source := os.Stdin
	var err error
	if len(os.Args) > 3 {
		source, err = os.Open(os.Args[3])
		if err != nil {
			fmt.Printf("Could not open resource: %v\nUsing stdin\n", os.Args[1])
			source = os.Stdin
		}
	}

	rdr := NewCustomReader(source, os.Args[1], os.Args[2])
	io.Copy(os.Stdout, rdr)
}
