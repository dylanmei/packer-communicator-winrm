package winrm

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"github.com/dylanmei/packer-communicator-winrm/winrm/addressing"
	"github.com/dylanmei/packer-communicator-winrm/winrm/envelope"
	"github.com/dylanmei/packer-communicator-winrm/winrm/rsp"
	"github.com/dylanmei/packer-communicator-winrm/winrm/wsman"
	"github.com/mitchellh/packer/common/uuid"
	"io"
	"io/ioutil"
	"net/http"
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

	rs := decodeResponse(response)
	return &Shell{rs.ShellId, user, pass}, nil
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
	xml, err := xml.MarshalIndent(env, " ", "	")
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewReader(xml))
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
	return bytes.NewReader(body), nil
}

type remoteShell struct {
	ShellId string `xml:"ShellId"`
}

func decodeResponse(reader io.Reader) *remoteShell {
	decoder := xml.NewDecoder(reader)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Space == NS_SOAP_ENV && se.Name.Local == "Header" {
				var h envelope.Header
				decoder.DecodeElement(&h, &se)
			}
			if se.Name.Space == NS_WIN_SHELL && se.Name.Local == "Shell" {
				var rs remoteShell
				decoder.DecodeElement(&rs, &se)
				return &rs
			}
		}
	}

	return nil
}
