package comments

import (
	"io"
)

const (
	defaultDelim byte = '#'
)

type reader struct {
	r io.Reader
	state
	delim byte
}

// NewReader returns an io.Reader which copies directly from r, ignoring '#' and
// any characters following it on the same line.
func NewReader(r io.Reader) io.Reader { return &reader{r, text, defaultDelim} }

// NewCustomReader is identical to NewReader, except that it accepts a custom
// delimeter
func NewCustomReader(r io.Reader, delim byte) io.Reader { return &reader{r, text, delim} }

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	buf = buf[:n]
	var wcount, rcount int
	for rcount != len(buf) {
		var written, read int
		written, read, r.state = r.state(buf[wcount:], buf[rcount:], r.delim)
		wcount += written
		rcount += read
	}
	return wcount, err
}

type state func(dst, src []byte, delim byte) (written, read int, next state)

func text(dst, src []byte, delim byte) (written, read int, next state) {
	for i, b := range src {
		if b == delim {
			return i, i + 1, comment
		}
		dst[i] = b
	}
	return len(src), len(src), text
}

func comment(dst, src []byte, delim byte) (written, read int, next state) {
	for i, b := range src {
		if b == '\n' {
			dst[0] = '\n'
			return 1, i + 1, text
		}
	}
	return 0, len(src), comment
}
