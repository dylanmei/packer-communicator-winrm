package winrm

import (
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm/winrmtest"
	"net/http"
	"testing"
)

func Test_creating_a_shell(t *testing.T) {
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		if r.Header.Get("Authorization") != "Basic dmFncmFudDp2YWdyYW50" {
			t.Fatal("bad authorization")
		}
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

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if s.Endpoint != fixture.Endpoint {
		t.Fatal("bad endpoint:", s.Endpoint)
	}

	if s.Id != "ABCXYZ" {
		t.Fatal("bad shell id:", s.Id)
	}

	if s.Owner != "vagrant" {
		t.Fatal("bad owner:", s.Owner)
	}
}

func Test_creating_a_shell_command(t *testing.T) {
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		if r.XmlString("//Header/SelectorSet[Selector='ABCXYZ']") == "" {
			t.Fatal("bad request: selector")
		}
		if r.XmlString("//Body/CommandLine[Command='foo bar']") == "" {
			t.Fatal("bad request: command")
		}

		fmt.Fprintf(w, `
			<Envelope>
				<s:Body>
					<rsp:CommandResponse>
						<rsp:CommandId>123456</rsp:CommandId>
					</rsp:CommandResponse>
				</s:Body>
			</Envelope>`)
	})

	s := &Shell{
		Id:       "ABCXYZ",
		Endpoint: fixture.Endpoint,
	}

	c, err := s.NewCommand("foo bar")

	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if c.Id != "123456" {
		t.Fatal("bad command id:", c.Id)
	}

	if c.CommandText != "foo bar" {
		t.Fatal("bad command text:", c.CommandText)
	}
}

func Test_deleting_a_shell(t *testing.T) {
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		if r.XmlString("//Header/SelectorSet[Selector='ABCXYZ']") == "" {
			t.Fatal("bad request: selector")
		}
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

	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func Test_authentication_failure(t *testing.T) {
	fixture := winrmtest.NewFixture()

	fixture.HandleFunc(func(w http.ResponseWriter, r *winrmtest.Request) {
		w.WriteHeader(401)
	})

	_, err := NewShell(fixture.Endpoint, "", "")

	if err == nil {
		t.Fatal("bad: no error")
	}

	herr, ok := err.(*HttpError)
	if !ok {
		t.Fatal("bad: not an http error")
	}

	if herr.StatusCode != 401 {
		t.Fatal("bad: http status code", herr.StatusCode)
	}
}
