package limits

import (
	"github.com/gozips/source"
	"io"
)

// lc is a LimitedReader + ReadCloser
type lc struct {
	*io.LimitedReader
	io.ReadCloser
}

func newlc(n int64, r io.ReadCloser) *lc {
	return &lc{
		LimitedReader: &io.LimitedReader{
			R: r,
			N: n,
		},
		ReadCloser: r,
	}
}

// Read attempts to read one more byte beyond the limit, if successful it
// returns a ReadError
func (l lc) Read(b []byte) (int, error) {
	n, err := l.LimitedReader.Read(b)
	if err == io.EOF {
		r := l.LimitedReader.R
		m, _ := r.Read(make([]byte, 1))
		if m > 0 {
			return n, source.ReadError{"error: size: exceeded limit"}
		}
	}

	return n, err
}

// Size reads up to EOF or n bytes which ever comes first
func Size(n int64, s source.Func) source.Func {
	return func(p string) (string, io.ReadCloser, error) {
		name, r, err := s.Readfrom(p)
		if err != nil {
			return name, r, err
		}

		l := newlc(n, r)
		return name, l, err
	}
}
