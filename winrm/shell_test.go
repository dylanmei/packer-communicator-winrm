package winrm

import (
	"fmt"
	. "github.com/onsi/gomega"
	"io"
	"launchpad.net/xmlpath"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mux *http.ServeMux

func setup(t *testing.T) (string, func()) {
	RegisterTestingT(t)

	mux = http.NewServeMux()
	server := httptest.NewServer(mux)
	return server.URL, func() { server.Close() }
}

func Test_creating_a_shell(t *testing.T) {
	url, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Expect(r.Method).To(Equal("POST"))
		Expect(r.Header.Get("Content-Type")).To(ContainSubstring("application/soap+xml"))
		Expect(r.Header.Get("Authorization")).To(Equal("Basic dmFncmFudDp2YWdyYW50"))
		fmt.Fprintf(w, `
            <Envelope>
                <s:Body>
                    <rsp:Shell>
                        <rsp:ShellId>ABCXYZ</rsp:ShellId>
                    </rsp:Shell>
                </s:Body>
            </Envelope>`)
	})

	s, err := NewShell(url, "vagrant", "vagrant")

	Expect(err).To(BeNil())
	if err == nil {
		Expect(s.Endpoint).To(Equal(url))
		Expect(s.Id).To(Equal("ABCXYZ"))
		Expect(s.Owner).To(Equal("vagrant"))
	}
}

func Test_deleting_a_shell(t *testing.T) {
	url, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Expect(r.Method).To(Equal("POST"))
		Expect(r.Body).To(ContainXml("//Header/SelectorSet[Selector='ABCXYZ']"))
		fmt.Fprintf(w, `
            <Envelope>
                <s:Body></s:Body>
            </Envelope>`)
	})

	s := &Shell{
		Id:       "ABCXYZ",
		Endpoint: url,
	}

	err := s.Delete()
	Expect(err).To(BeNil())
}

func Test_authentication_failure(t *testing.T) {
	url, teardown := setup(t)
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
	})

	_, err := NewShell(url, "", "")
	Expect(err).ToNot(BeNil())
	Expect(err).To(BeAssignableToTypeOf((*HttpError)(nil)))
}

func ContainXml(expected string) OmegaMatcher {
	return &xpathMatcher{
		text: expected,
		path: xmlpath.MustCompile(expected),
	}
}

type xpathMatcher struct {
	text string
	path *xmlpath.Path
}

func (matcher *xpathMatcher) Match(actual interface{}) (success bool, message string, err error) {
	reader, ok := actual.(io.Reader)
	if !ok {
		return false, "", fmt.Errorf("ContainXml expects a []byte or an io.Reader")
	}

	node, err := xmlpath.Parse(reader)
	if err != nil {
		return false, "", err
	}

	_, ok = matcher.path.String(node)
	if ok {
		return true, fmt.Sprintf("Expected\n\t%#v\nnot to match xml-path\n\t%#v", "<todo/>", matcher.text), nil
	} else {
		return false, fmt.Sprintf("Expected\n\t%#v\nto to match xml-path\n\t%#v", "<todo/>", matcher.text), nil
	}
}
