package limits

import (
	"github.com/gozips/source"
	"io"
)

// Size reads up to EOF or n bytes, on each file, which ever comes first
func Size(n int64, s source.Func) source.Func {
	return func(p string) (string, io.ReadCloser, error) {
		name, r, err := s.Readfrom(p)
		if err != nil {
			return name, r, err
		}

		l := newLrc(n, r)
		return name, l, nil
	}
}
