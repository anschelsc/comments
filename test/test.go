package main

import (
	"github.com/anschelsc/comments"
	"io"
	"os"
)

func main() {
	io.Copy(os.Stdout, comments.NewReader(os.Stdin))
}
