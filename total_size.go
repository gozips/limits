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
	var st = &stat{n}

	return func(p string) (string, io.ReadCloser, error) {
		name, r, err := s.Readfrom(p)
		if err != nil {
			return name, r, err
		}

		c := NewCombined(st, r)
		c.Name = name

		return name, c, nil
	}
}
