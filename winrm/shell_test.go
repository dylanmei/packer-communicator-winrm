package winrm

import (
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm/winrmtest"
	. "github.com/dylanmei/packer-communicator-winrm/winrm/winrmtest/matchers"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

func Test_creating_a_shell(t *testing.T) {
	RegisterTestingT(t)
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
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

	s, err := NewShell(fixture.Endpoint, "vagrant", "vagrant")

	Expect(err).To(BeNil())
	if err == nil {
		Expect(s.Endpoint).To(Equal(fixture.Endpoint))
		Expect(s.Id).To(Equal("ABCXYZ"))
		Expect(s.Owner).To(Equal("vagrant"))
	}
}

func Test_creating_a_shell_command(t *testing.T) {
	RegisterTestingT(t)
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		Expect(r).To(MatchXmlPath("//Header/SelectorSet[Selector='ABCXYZ']"))
		Expect(r).To(MatchXmlPath("//Body/CommandLine[Command='foo bar']"))
		fmt.Fprintf(w, `
			<Envelope>
				<s:Body>
					<rsp:CommandResponse>
						<rsp:CommandId>123789</rsp:CommandId>
					</rsp:CommandResponse>
				</s:Body>
			</Envelope>`)
	})

	s := &Shell{
		Id:       "ABCXYZ",
		Endpoint: fixture.Endpoint,
	}

	c, err := s.NewCommand("foo bar")

	Expect(err).To(BeNil())
	if err == nil {
		Expect(c.Id).To(Equal("123789"))
		Expect(c.CommandText).To(Equal("foo bar"))
	}
}

func Test_deleting_a_shell(t *testing.T) {
	RegisterTestingT(t)
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		Expect(r).To(MatchXmlPath("//Header/SelectorSet[Selector='ABCXYZ']"))
		fmt.Fprintf(w, `
			<Envelope>
				<s:Body></s:Body>
			</Envelope>`)
	})

	s := &Shell{
		Id:       "ABCXYZ",
		Endpoint: fixture.Endpoint,
	}

	err := s.Delete()
	Expect(err).To(BeNil())
}

func Test_authentication_failure(t *testing.T) {
	RegisterTestingT(t)
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		w.WriteHeader(401)
	})

	_, err := NewShell(fixture.Endpoint, "", "")
	Expect(err).ToNot(BeNil())
	Expect(err).To(BeAssignableToTypeOf((*HttpError)(nil)))
}
