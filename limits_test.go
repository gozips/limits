package limits

import (
	"fmt"
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"testing"
)

func e(p, s string) string {
	return fmt.Sprintf(p, s)
}

func TestCombinedLimits(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := Count(3, Size(7, sources.HTTP))
	i := 0
	for _, v := range []struct {
		u, b string
	}{
		{u("one.txt"), "one"},
		{u("12.txt"), "Hello W"},
		{u("two.txt"), "two"},
		{u("three.txt"), ""},
		{u("four.txt"), ""},
	} {
		_, b, err := fn(v.u)
		if err != nil && err.Error() == "error: count: exceeded limit" {
			i++
		} else {
			buf, _ := readb(b)
			assert.Equal(t, buf.String(), v.b)
		}
	}
	assert.Equal(t, 2, i)
}

func TestCombinedLimitsOnSize(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := Count(3, TotalSize(12, Size(7, sources.HTTP)))
	for _, v := range []struct {
		u, b string
	}{
		{"one.txt", "one"},
		{"12.txt", "Hello W"},
		{"two.txt", "tw"},
		{"three.txt", ""},
		{"four.txt", ""},
	} {
		_, b, err := fn(u(v.u))
		if v.b != "" {
			buf, err := readb(b)
			assert.Equal(t, buf.String(), v.b)

			if v.b == "Hello W" {
				assert.Equal(t, e("error: size: %s exceeded limit", v.u), err.Error())
			}

			if v.b == "tw" {
				assert.Equal(t,
					e("error: total size: %s exceeded limit; File truncated", v.u),
					err.Error())
			}
		} else {
			assert.Equal(t, "error: count: exceeded limit", err.Error())
		}
	}
}

func TestCombinedLimitsOnSizeReversedOrder(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := Count(3, Size(7, TotalSize(12, sources.HTTP)))
	for _, v := range []struct {
		u, b string
	}{
		{"one.txt", "one"},
		{"12.txt", "Hello W"},
		{"two.txt", "t"}, // NOTE we lose a byte if the previous exceeded size
		{"three.txt", ""},
		{"four.txt", ""},
	} {
		_, b, err := fn(u(v.u))
		if v.b != "" {
			buf, err := readb(b)
			assert.Equal(t, buf.String(), v.b)

			if v.b == "Hello W" {
				assert.Equal(t, e("error: size: %s exceeded limit", v.u), err.Error())
			}

			if v.b == "t" {
				assert.Equal(t,
					e("error: total size: %s exceeded limit; File truncated", v.u),
					err.Error())
			}
		} else {
			assert.Equal(t, "error: count: exceeded limit", err.Error())
		}
	}
}
