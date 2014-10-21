package limits

import (
	"bytes"
	"fmt"
	"github.com/gozips/sources"
	"github.com/nowk/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func bodyh(str string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(str))
	}
}

func tServer() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/file/one.txt", bodyh("one"))
	mux.HandleFunc("/file/two.txt", bodyh("two"))
	mux.HandleFunc("/file/three.txt", bodyh("three"))
	mux.HandleFunc("/file/four.txt", bodyh("four"))
	mux.HandleFunc("/file/five.txt", bodyh("five"))

	ts := httptest.NewServer(mux)
	return ts
}

func bodyrc(str string) io.ReadCloser {
	b := bytes.NewReader([]byte(str))
	return ioutil.NopCloser(b)
}

func urlfn(urlstr string) func(string) string {
	return func(f string) string {
		return fmt.Sprintf("%s/file/%s", urlstr, f)
	}
}

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
