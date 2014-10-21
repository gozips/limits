package limits

import (
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"testing"
)

func TestCountCountsErroredReadAttempts(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := Count(3, sources.HTTP)
	i := 0
	for _, v := range []string{
		u("one.txt"),
		"thiswillerror.txt",
		u("three.txt"),
		u("four.txt"),
		u("five.txt"),
	} {
		_, _, err := fn(v)
		if err != nil && err.Error() == "error: count: exceeded limit" {
			i++
		}
	}
	assert.Equal(t, 2, i)
}

func TestCountLaxDoesNotCountErroredReadAttempts(t *testing.T) {
	var ts = tServer()
	defer ts.Close()
	u := urlfn(ts.URL)

	fn := CountLax(3, sources.HTTP)
	i := 0
	for _, v := range []string{
		u("one.txt"),
		"thiswillerror.txt",
		u("three.txt"),
		u("four.txt"),
		u("five.txt"),
	} {
		_, _, err := fn(v)
		if err != nil && err.Error() == "error: count: exceeded limit" {
			i++
		}
	}
	assert.Equal(t, 1, i)
}
