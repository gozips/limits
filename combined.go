package limits

import (
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

// Combined wraps LimitedRc and keeps track of what it has read as a total
// combined
type Combined struct {
	*Limited
	*stat
}

func NewCombined(s *stat, r io.ReadCloser) *Combined {
	return &Combined{
		Limited: NewLimited(s.n, r),
		stat:    s,
	}
}

// Read decrements the stat n on each read
func (r Combined) Read(b []byte) (int, error) {
	if r.stat.Exceeded() {
		return 0, source.ReadError{"error: total size: exceeded limit"}
	}

	n, err := r.Limited.Read(b)
	r.stat.Decrement(int64(n))
	return n, err
}
