package winrmtest

import (
	"bytes"
	"io"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"net/http"
	"strings"
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

func (r *Request) XmlString(query string) string {
	xpath, err := xmlpath.Compile(query)

	if err != nil {
		return ""
	}

	buffer, _ := ioutil.ReadAll(r.reader)
	r.reader.Seek(0, 0)

	node, err := xmlpath.Parse(bytes.NewReader(buffer))

	if err != nil {
		return ""
	}

	result, _ := xpath.String(node)
	return strings.Trim(result, " \r\n\t")
}
