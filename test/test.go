package main

import (
	"github.com/anschelsc/comments"
	"io"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		io.Copy(os.Stdout, comments.NewCustomReader(os.Stdin, os.Args[1][0]))
	} else {
		io.Copy(os.Stdout, comments.NewReader(os.Stdin))
	}
}
