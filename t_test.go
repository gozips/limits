package limits

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

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

	mux.HandleFunc("/file/large.txt", bodyh(`{"data": ["one"], "meta": {"code":200}}`))
	mux.HandleFunc("/file/medium.txt", bodyh("Hello World!"))
	mux.HandleFunc("/file/small.txt", bodyh("abc"))

	ts := httptest.NewServer(mux)
	return ts
}

func urlfn(urlstr string) func(string) string {
	return func(f string) string {
		return fmt.Sprintf("%s/file/%s", urlstr, f)
	}
}

func readb(r io.ReadCloser) (*bytes.Buffer, error) {
	defer r.Close()
	var b []byte
	buf := bytes.NewBuffer(b)
	_, err := io.Copy(buf, r)
	return buf, err
}
