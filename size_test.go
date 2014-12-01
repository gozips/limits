package limits

import (
	"fmt"
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"testing"
)

func TestSizeReturnsUpToLimit(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := Size(7, sources.HTTP)
	for _, v := range []struct {
		u, b string
		e    bool
	}{
		{"39.txt", `{"data"`, true},
		{"12.txt", `Hello W`, true},
		{"3.txt", `abc`, false},
	} {
		_, r, _ := fn(u(v.u))
		buf, err := readb(r)
		assert.Equal(t, buf.String(), v.b)
		if v.e {
			assert.Equal(t, fmt.Sprintf("error: size: %s exceeded limit", v.u),
				err.Error())
		}
	}
}
