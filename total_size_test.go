package limits

import (
	"fmt"
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
		t    int
	}{
		{"39.txt", `{"data": ["one"], "meta": {"code":200}}`, false, 0},
		{"12.txt", "Hello World!", false, 0},
		{"3.txt", "a", true, 1},
		{"one.txt", "", true, 2},
	} {
		_, r, _ := fn(u(v.u))
		b, err := readb(r)
		assert.Equal(t, v.b, b.String())
		assert.True(t, (err != nil) == v.e)
		if v.e {
			x := "; File truncated"
			if v.t == 2 {
				x = "; File skipped"
			}

			assert.Equal(t,
				fmt.Sprintf("error: total size: %s exceeded limit%s", v.u, x),
				err.Error())
		}
	}
}
