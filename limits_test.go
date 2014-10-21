package limits

import (
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"testing"
)

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
