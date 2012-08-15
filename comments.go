// The comments package provides a wrapper for an io.Reader which ignores text inside comments.
package comments

import (
	"io"
)

var (
	defaultStart  = "#"
	defaultEnd    = "\n"
	defaultAppend = "\n"
)

type delim struct {
	// delimiter byte sequence
	del []byte
	// Buffer of back-write characters
	buf    []byte
	dellen int
	buflen int
	delpos int
	// Buffer start pos
	bufspos int
	// Buffer end pos
	bufepos int
}

type reader struct {
	r io.Reader
	state
	start  delim
	stop   delim
	append string
}

// NewReader returns an io.Reader which copies directly from r, ignoring '#' and
// any characters following it on the same line.
func NewReader(r io.Reader) io.Reader {
	return NewCustomReader(r, defaultStart, defaultEnd, defaultAppend)
}

// NewCustomReader is identical to NewReader, except that it accepts a custom
// start and stop delimiters, and an optional string to append after comments end.
// This may be useful for single-line comments (like bash's # or C's //) whose
// stop delimiter is a newline.
//
// NOTE: append MUST NOT BE longer than end; this will cause undefined behavior
// and possibly a nil-pointer dereference.
func NewCustomReader(r io.Reader, start, end, append string) io.Reader {
	return &reader{r, fstate(text),
		delim{[]byte(start), make([]byte, len(start)+len(append)), len(start), (len(start) + len(append)), 0, 0, 0},
		delim{[]byte(end), nil, len(end), 0, 0, 0, 0},
		append}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	buf = buf[:n]
	var wcount, rcount int
	for rcount != len(buf) {
		var written, read int
		written, read, r.state = r.state.run(buf[wcount:], buf[rcount:], r)
		wcount += written
		rcount += read
	}
	return wcount, err
}

type state interface {
	run(dst, src []byte, r *reader) (written, read int, next state)
}

type fstate func(dst, src []byte, r *reader) (written, read int, next state)

func (f fstate) run(dst, src []byte, r *reader) (written, read int, next state) { return f(dst, src, r) }

func text(dst, src []byte, r *reader) (written, read int, next state) {
	delim := &r.start
	slen := len(src)
	dlen := len(dst)
	var wc, rc int
	for rc < slen && wc < dlen {
		// If possible, write one
		// character from the buffer
		if delim.bufspos != delim.bufepos {
			dst[wc] = delim.buf[delim.bufspos]
			delim.bufspos = (delim.bufspos + 1) % delim.buflen
			wc++
		}

		if src[rc] == delim.del[delim.delpos] {
			delim.delpos++
			if delim.delpos == delim.dellen {
				delim.delpos = 0
				return wc, rc + 1, fstate(comment)
			}
		} else {
			if delim.delpos != 0 {
				for i := 0; i < delim.delpos; i++ {
					delim.buf[(delim.bufepos+i)%delim.buflen] = delim.del[i]
				}
				delim.bufepos = (delim.bufepos + delim.delpos) % delim.buflen
				delim.delpos = 0
			}
			// Either write another character
			// from the buffer and write one
			// from the input to the buffer or
			// write directly from input.
			if delim.bufspos != delim.bufepos {
				dst[wc] = delim.buf[delim.bufspos]
				delim.buf[delim.bufepos] = src[rc]
				delim.bufspos = (delim.bufspos + 1) % delim.buflen
				delim.bufepos = (delim.bufepos + 1) % delim.buflen
			} else {
				dst[wc] = src[rc]
			}
			wc++
		}
		rc++
	}

	// Flush the buffer
	for delim.bufspos != delim.bufepos && wc < dlen {
		dst[wc] = delim.buf[delim.bufspos]
		delim.bufspos = (delim.bufspos + 1) % delim.buflen
		wc++
	}
	return wc, rc, fstate(text)
}

func comment(dst, src []byte, r *reader) (written, read int, next state) {
	sdelim := &r.start
	edelim := &r.stop
	slen := len(src)
	dlen := len(dst)
	var wc, rc int
	for rc < slen && wc < dlen {
		// If possible, write one
		// character from the buffer
		// // fmt.Println(sdelim)
		if sdelim.bufspos != sdelim.bufepos {
			dst[wc] = sdelim.buf[sdelim.bufspos]
			sdelim.bufspos = (sdelim.bufspos + 1) % sdelim.buflen
			wc++
		}

		if src[rc] == edelim.del[edelim.delpos] {
			edelim.delpos++
			if edelim.delpos == edelim.dellen {
				edelim.delpos = 0
				// Add append to buffer
				i := 0
				for ; i < len(r.append) && wc < dlen; i++ {
					dst[i+wc] = r.append[i]
					wc++
				}
				for ; i < len(r.append); i++ {
					sdelim.buf[sdelim.bufspos] = r.append[i]
					sdelim.bufspos = (sdelim.bufspos + 1) % sdelim.buflen
				}
				return wc, rc + 1, fstate(text)
			}
		} else {
			if edelim.delpos != 0 {
				edelim.delpos = 0
			}

			// If there's anything left, keep
			// writing from the text buffer
			if sdelim.bufspos != sdelim.bufepos {
				dst[wc] = sdelim.buf[sdelim.bufspos]
				sdelim.bufspos = (sdelim.bufspos + 1) % sdelim.buflen
				wc++
			}
		}
		rc++
	}

	// Flush the buffer
	for sdelim.bufspos != sdelim.bufepos && wc < dlen {
		dst[wc] = sdelim.buf[sdelim.bufspos]
		sdelim.bufspos = (sdelim.bufspos + 1) % sdelim.buflen
		wc++
	}
	return 0, rc, fstate(comment)
}
