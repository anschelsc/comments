package comments

import (
	"io"
)

type Reader struct {
	r io.Reader
	state
}

func NewReader(r io.Reader) *Reader { return &Reader{r, text} }

func (r *Reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	buf = buf[:n]
	var wcount, rcount int
	for rcount != len(buf) {
		var written, read int
		_ = buf[rcount:]
		_ = buf[wcount:]
		written, read, r.state = r.state(buf[wcount:], buf[rcount:])
		wcount += written
		rcount += read
	}
	return wcount, err
}

type state func(dst, src []byte) (written, read int, next state)

func text(dst, src []byte) (written, read int, next state) {
	for i, b := range src {
		if b == '#' {
			return i, i + 1, comment
		}
		dst[i] = b
	}
	return len(src), len(src), text
}

func comment(dst, src []byte) (written, read int, next state) {
	for i, b := range src {
		if b == '\n' {
			dst[0] = '\n'
			return 1, i + 1, text
		}
	}
	return 0, len(src), comment
}