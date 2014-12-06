package limits

import (
	"fmt"
	"github.com/gozips/source"
	"io"
)

type stat struct {
	MaxReadSize int64
}

func (s *stat) Decrement(n int64) {
	s.MaxReadSize -= n
}

func (s stat) Exceeded() bool {
	return s.MaxReadSize <= 0
}

// Combined wraps Limited and keeps track of total combined read size for all
// reads. It truncates only the one file that exceeds the total limit set.
type Combined struct {
	*Limited
	*stat

	N int64 // track how much was currently read
}

func NewCombined(s *stat, r io.ReadCloser) *Combined {
	return &Combined{
		Limited: NewLimited(s.MaxReadSize, r),
		stat:    s,
	}
}

// Read decrements the stat n on each read
func (r *Combined) Read(b []byte) (int, error) {
	if r.stat.Exceeded() {
		x := "; File skipped"
		if r.N > 0 {
			x = "; File truncated"
		}

		return 0, source.ReadError{
			fmt.Sprintf("error: total size: %s exceeded limit%s", r.Name, x),
		}
	}

	n, err := r.Limited.Read(b)
	r.stat.Decrement(int64(n))
	r.N += int64(n) // increment current read attempt
	return n, err
}
