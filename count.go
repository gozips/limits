package limits

import (
	"errors"
	fp "github.com/gozips/filepath"
	"github.com/gozips/source"
	"io"
)

// count increments at every attempt to Readfrom(p). If an error is returned and
// strict == false, then that attempt will not be counted
// Attempts after the n limit will return an error, no attemps to Readfrom(p)
// will be made
func count(n int, strict bool, s source.Func) source.Func {
	var i = 0

	return func(p string) (string, io.ReadCloser, error) {
		if i >= n {
			return fp.Base(p), nil, errors.New("exceeded maximum number of files")
		}

		name, r, err := s.Readfrom(p)
		defer func(er error) {
			i++
			if !strict && er != nil {
				i--
			}
		}(err)

		if err != nil {
			return name, r, err
		}
		return name, r, err
	}
}

// Count counts every attempt
func Count(n int, s source.Func) source.Func {
	return count(n, true, s)
}

// CountLax counts only attempts that result in no error
func CountLax(n int, s source.Func) source.Func {
	return count(n, false, s)
}
