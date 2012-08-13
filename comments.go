package comments

import (
	"io"
)

var defaultDelim = &Delim{"#", "\n", true}

type Delim struct {
	Start, Stop string
	WriteStop   bool
}

type reader struct {
	r io.Reader
	state
	d *Delim
}

// NewReader returns an io.Reader which copies directly from r, ignoring '#' and
// any characters following it on the same line.
func NewReader(r io.Reader) io.Reader { return NewCustomReader(r, defaultDelim) }

// NewCustomReader is identical to NewReader, except that it accepts a custom
// delimeter
func NewCustomReader(r io.Reader, d *Delim) io.Reader { return &reader{r, fstate(text), d} }

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	buf = buf[:n]
	var wcount, rcount int
	for rcount != len(buf) {
		var written, read int
		written, read, r.state = r.state.run(buf[wcount:], buf[rcount:], r.d)
		wcount += written
		rcount += read
	}
	return wcount, err
}

type state interface {
	run(dst, src []byte, d *Delim) (written, read int, next state)
}

type fstate func(dst, src []byte, d *Delim) (written, read int, next state)

func (f fstate) run(dst, src []byte, d *Delim) (written, read int, next state) { return f(dst, src, d) }

func text(dst, src []byte, d *Delim) (written, read int, next state) {
	for i, b := range src {
		if b == d.Start[0] {
			return i, i + 1, fstate(comment)
		}
		dst[i] = b
	}
	return len(src), len(src), fstate(text)
}

func comment(dst, src []byte, d *Delim) (written, read int, next state) {
	for i, b := range src {
		if b == d.Stop[0] {
			if d.WriteStop {
				dst[0] = d.Stop[0]
				return 1, i + 1, fstate(text)
				return 0, i, &strWriter{d.Stop, fstate(text)}
			}
			return 0, i + 1, fstate(text)
		}
	}
	return 0, len(src), fstate(comment)
}

type strWriter struct {
	string
	state
}

func (s *strWriter) run(dst, _ []byte, _ *Delim) (written, read int, next state) {
	n := copy(dst, s.string)
	if n != len(s.string) {
		panic("Delimiter could not be written back to the buffer. Impossible.")
	}
	return n, n, s.state
}
