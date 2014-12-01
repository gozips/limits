package limits

import (
	"fmt"
	"github.com/gozips/source"
	"io"
)

type stat struct {
	n int64
}

func (s *stat) Decrement(n int64) {
	s.n -= n
}

func (s stat) Exceeded() bool {
	return s.n <= 0
}

// lrc is a LimitedReader + ReadCloser
type lrc struct {
	name string
	*io.LimitedReader
	io.ReadCloser
}

func newLrc(name string, n int64, r io.ReadCloser) *lrc {
	return &lrc{
		name: name,
		LimitedReader: &io.LimitedReader{
			R: r,
			N: n,
		},
		ReadCloser: r,
	}
}

func (r lrc) Read(b []byte) (int, error) {
	n, err := r.LimitedReader.Read(b)
	if err == io.EOF {
		m, _ := r.LimitedReader.R.Read(make([]byte, 1))
		if m > 0 {
			return n, source.ReadError{fmt.Sprintf("error: size: %s exceeded limit", r.name)}
		}
	}

	return n, err
}

// slrc wraps lrc and keeps track of what it has read as a total combined
type slrc struct {
	*lrc
	*stat
}

func newSlrc(name string, s *stat, r io.ReadCloser) *slrc {
	return &slrc{
		lrc:  newLrc(name, s.n, r),
		stat: s,
	}
}

// Read decrements the stat n on each read
func (r slrc) Read(b []byte) (int, error) {
	if r.stat.Exceeded() {
		return 0, source.ReadError{"error: total size: exceeded limit"}
	}

	n, err := r.lrc.Read(b)
	r.stat.Decrement(int64(n))
	return n, err
}
