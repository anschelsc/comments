package comments

import (
	"io"
)

const (
	defaultStart = '#'
	defaultStop  = '\n'
)

type reader struct {
	r io.Reader
	state
	start []byte
	stop  []byte
}

// NewReader returns an io.Reader which copies directly from r, ignoring '#' and
// any characters following it on the same line.
func NewReader(r io.Reader) io.Reader {
	return &reader{r, text, []byte{defaultStart}, []byte{defaultStop}}
}

// NewCustomReader is identical to NewReader, except that it accepts custom
// start and stop delimeters
func NewCustomReader(r io.Reader, start, stop string) io.Reader {
	return &reader{r, text, []byte(start), []byte(stop)}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	buf = buf[:n]
	var wcount, rcount int
	for rcount != len(buf) {
		var written, read int
		_ = buf[rcount:]
		_ = buf[wcount:]
		written, read, r.state = r.state(buf[wcount:], buf[rcount:], r.start, r.stop)
		wcount += written
		rcount += read
	}
	return wcount, err
}

type state func(dst, src []byte, start, stop []byte) (written, read int, next state)

func text(dst, src []byte, start, stop []byte) (written, read int, next state) {
	dLen := len(start)
	sLen := len(src)
CHECK:
	for i, b := range src {
		if b == start[0] && dLen <= sLen-i {
			for j := 1; j < dLen; j++ {
				if src[i+j] != start[j] {
					continue CHECK
				}
			}
			return i, i + dLen, comment
		}
		dst[i] = b
	}
	return sLen, sLen, text
}

func comment(dst, src []byte, start, stop []byte) (written, read int, next state) {
	dLen := len(stop)
	sLen := len(src)
CHECK:
	for i, b := range src {
		if b == stop[0] && dLen <= sLen-i {
			for j := 1; j < dLen; j++ {
				if src[i+j] != stop[j] {
					continue CHECK
				}
			}
			return 0, i + dLen, text
		}
	}
	return 0, sLen, comment
}
