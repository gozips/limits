package limits

import (
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"testing"
)

func TestTotalSizeReadsTillTotalSizeIsMet(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := TotalSize((39 + 12 + 1), sources.HTTP)
	for _, v := range []struct {
		u, b string
		e    bool
	}{
		{u("39.txt"), `{"data": ["one"], "meta": {"code":200}}`, false},
		{u("12.txt"), "Hello World!", false},
		{u("3.txt"), "a", true},
		{u("one.txt"), "", true},
	} {
		_, r, _ := fn(v.u)
		b, err := readb(r)
		assert.Equal(t, v.b, b.String())
		assert.True(t, (err != nil) == v.e)
		if v.e {
			assert.Equal(t, "error: total size: exceeded limit", err.Error())
		}
	}
}
