package limits

import (
	"fmt"
	"github.com/gozips/source"
	"io"
)

// Limited is a LimitedReader + ReadCloser. It truncates all files that exceed
// the set limit.
type Limited struct {
	Name string

	*io.LimitedReader
	io.ReadCloser
}

func NewLimited(n int64, r io.ReadCloser, opts ...func(*Limited)) *Limited {
	l := &Limited{
		LimitedReader: &io.LimitedReader{r, n},
		ReadCloser:    r,
	}

	for _, v := range opts {
		v(l)
	}

	return l
}

func (r Limited) Read(b []byte) (int, error) {
	n, err := r.LimitedReader.Read(b)
	if err == io.EOF {
		m, _ := r.LimitedReader.R.Read(make([]byte, 1))
		if m > 0 {
			return n, source.ReadError{fmt.Sprintf("error: size: %s exceeded limit",
				r.Name)}
		}
	}

	return n, err
}
