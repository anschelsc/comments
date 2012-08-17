package comments

import (
	"bufio"
	"github.com/anschelsc/finder"
	"io"
)

type reader struct {
	r          io.ByteReader
	begin, end *finder.Finder
	add        []byte
	state
}

// NewCustomReader returns an io.Reader that copies directly from r, replacing
// any string between begin and end with add.
func NewCustomReader(r io.Reader, begin, end, add []byte) io.Reader {
	br, ok := r.(io.ByteReader)
	if !ok {
		br = bufio.NewReader(r)
	}
	ret := &reader{
		r:     br,
		begin: finder.Compile(begin),
		end:   finder.Compile(end),
		add:   add,
	}
	ret.state = text{finder.NewReader(ret.begin, ret.r)}
	return ret
}

// NewReader calls NewCustomReader with the appropriate arguments for parsing
// bash-style comments.
func NewReader(r io.Reader) io.Reader {
	return NewCustomReader(r, []byte{'#'}, []byte{'\n'}, []byte{'\n'})
}

func (r *reader) Read(buf []byte) (n int, err error) {
	for n != len(buf) && err == nil {
		var nn int
		nn, err, r.state = r.state.run(buf[n:], r)
		n += nn
	}
	return
}

type state interface {
	run(dst []byte, up *reader) (int, error, state)
}

type text struct {
	r io.Reader
}

func (t text) run(dst []byte, up *reader) (int, error, state) {
	n, err := io.ReadFull(t.r, dst)
	if err == finder.Found {
		return n, nil, comment{finder.NewReader(up.end, up.r)}
	}
	return n, err, t
}

type comment struct {
	r io.Reader
}

var cbuf [1024]byte

func (c comment) run(_ []byte, up *reader) (int, error, state) {
	_, err := io.ReadFull(c.r, cbuf[:])
	if err == finder.Found {
		return 0, nil, copier(up.add)
	}
	return 0, err, c
}

type copier []byte

func (c copier) run(dst []byte, up *reader) (int, error, state) {
	n := copy(dst, c)
	if n == len(c) {
		return n, nil, text{finder.NewReader(up.begin, up.r)}
	}
	return n, nil, c[n:]
}
