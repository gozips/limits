package limits

import (
	"github.com/gozips/source"
	"io"
)

// TotalSize reads up to a total of n of the combined sizes
// It's order combination counts. It should wrap the Size limiter, ex:
//
//		TotalSize(N, Size(M, source.Func))
//
// Running it in reverse will lose a byte each time the Size limiter exceeds
func TotalSize(n int64, s source.Func) source.Func {
	st := stat{
		n: n,
	}

	return func(p string) (string, io.ReadCloser, error) {
		name, r, err := s.Readfrom(p)
		if err != nil {
			return name, r, err
		}

		c := newSlrc(&st, r)
		return name, c, err
	}
}
