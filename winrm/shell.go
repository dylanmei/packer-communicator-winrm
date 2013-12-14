package winrm

import (
	"errors"
	"github.com/dylanmei/packer-communicator-winrm/winrm/addressing"
	"github.com/dylanmei/packer-communicator-winrm/winrm/envelope"
	"github.com/dylanmei/packer-communicator-winrm/winrm/rsp"
	"github.com/dylanmei/packer-communicator-winrm/winrm/wsman"
	"github.com/mitchellh/packer/common/uuid"
	"launchpad.net/xmlpath"
	"log"
)

type Shell struct {
	Id       string
	Owner    string
	password string
}

func NewShell(user, pass string) (*Shell, error) {
	header := &envelope.Header{
		MessageID: "uuid:" + uuid.TimeOrderedUUID(),
		Action: &addressing.AttributedURI{
			URI:            "http://schemas.xmlsoap.org/ws/2004/09/transfer/Create",
			MustUnderstand: true,
		},
		To: &addressing.AttributedURI{
			URI: "http://localhost:5985/wsman",
		},
		ReplyTo: &addressing.EndpointReference{
			&addressing.AttributedURI{
				URI:            "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
				MustUnderstand: true,
			},
		},
		MaxEnvelopeSize: &wsman.MaxEnvelopeSize{
			Value:          153600,
			MustUnderstand: true,
		},
		OperationTimeout: "PT60S",
		ResourceURI: &wsman.AttributableURI{
			URI:            "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
			MustUnderstand: true,
		},
		OptionSet: &wsman.OptionSet{
			Options: []*wsman.Option{
				&wsman.Option{"WINRS_NOPROFILE", "FALSE"},
				&wsman.Option{"WINRS_CODEPAGE", "437"},
			},
		},
	}

	env := &envelope.Envelope{
		Header: header,
		Body: &envelope.Body{
			&rsp.Shell{
				InputStreams:  "stdin",
				OutputStreams: "stdout stderr",
			},
		},
	}

	response, err := sendEnvelope(user, pass, env)
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("//Body/Shell/ShellId")
	root, err := xmlpath.Parse(response)
	if err != nil {
		log.Fatal(err)
	}

	shell, ok := path.String(root)
	if !ok {
		return nil, errors.New("Could not create shell.")
	}

	return &Shell{shell, user, pass}, nil
}

func (s *Shell) Delete() error {
	header := &envelope.Header{
		MessageID: "uuid:" + uuid.TimeOrderedUUID(),
		Action: &addressing.AttributedURI{
			URI:            "http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete",
			MustUnderstand: true,
		},
		To: &addressing.AttributedURI{
			URI: "http://localhost:5985/wsman",
		},
		ReplyTo: &addressing.EndpointReference{
			&addressing.AttributedURI{
				URI:            "http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous",
				MustUnderstand: true,
			},
		},
		MaxEnvelopeSize: &wsman.MaxEnvelopeSize{
			Value:          153600,
			MustUnderstand: true,
		},
		OperationTimeout: "PT60S",
		ResourceURI: &wsman.AttributableURI{
			URI:            "http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd",
			MustUnderstand: true,
		},
		SelectorSet: &wsman.SelectorSet{
			Selectors: []*wsman.Selector{
				&wsman.Selector{"ShellId", s.Id},
			},
		},
	}

	env := &envelope.Envelope{
		Header: header,
		Body:   &envelope.Body{},
	}

	_, err := sendEnvelope(s.Owner, s.password, env)
	return err
}
