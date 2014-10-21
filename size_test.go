package limits

import (
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
		{u("39.txt"), `{"data"`, true},
		{u("12.txt"), `Hello W`, true},
		{u("3.txt"), `abc`, false},
	} {
		_, r, _ := fn(v.u)
		buf, err := readb(r)
		assert.Equal(t, buf.String(), v.b)
		if v.e {
			assert.Equal(t, "error: size: exceeded limit", err.Error())
		}
	}
}
