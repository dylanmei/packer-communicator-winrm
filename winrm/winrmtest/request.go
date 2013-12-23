package winrmtest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type Request struct {
	*http.Request
	reader io.ReadSeeker
}

func newRequest(r *http.Request) *Request {
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	return &Request{
		Request: r, reader: bytes.NewReader(body),
	}
}

func (r *Request) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *Request) Seek(offset int64, whence int) (int64, error) {
	return r.reader.Seek(offset, whence)
}
