package winrm

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"github.com/dylanmei/packer-communicator-winrm/winrm/addressing"
	"github.com/dylanmei/packer-communicator-winrm/winrm/envelope"
	"github.com/dylanmei/packer-communicator-winrm/winrm/rsp"
	"github.com/dylanmei/packer-communicator-winrm/winrm/wsman"
	"github.com/mitchellh/packer/common/uuid"
	"io"
	"io/ioutil"
	"launchpad.net/xmlpath"
	"log"
	"net/http"
	"os"
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

func sendEnvelope(user, pass string, env *envelope.Envelope) (io.Reader, error) {
	xmlEnvelope, err := xml.MarshalIndent(env, " ", "	")
	if err != nil {
		return nil, err
	}

	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("sending", string(xmlEnvelope))
	}

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewReader(xmlEnvelope))
	request.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")
	request.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(user+":"+pass)))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, NewHttpError(response)
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if os.Getenv("WINRM_DEBUG") != "" {
		log.Println("receiving", string(body))
	}

	return bytes.NewReader(body), nil
}
