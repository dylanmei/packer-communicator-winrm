package winrm

// http://msdn.microsoft.com/en-us/library/cc251731.aspx

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
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

	xml, err := xml.MarshalIndent(env, " ", "	")
	if err != nil {
		return nil, err
	}
	//	os.Stdout.Write(xml)
	//	fmt.Println()

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
	body, _ := ioutil.ReadAll(response.Body)
	rs := decodeResponse(bytes.NewBuffer(body))

	return &Shell{rs.ShellId, user, pass}, nil
}

func (s *Shell) Delete() error {
	command := fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="http://www.w3.org/2003/05/soap-envelope"
	xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing"
	xmlns:b="http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd"
	xmlns:n="http://schemas.xmlsoap.org/ws/2004/09/enumeration"
	xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer"
	xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd"
	xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd"
	xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell"
	xmlns:cfg="http://schemas.microsoft.com/wbem/wsman/1/config"
	xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
>
  <env:Header>
	<a:To>http://localhost:5985/wsman</a:To>
	<a:ReplyTo>
	  <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
	</a:ReplyTo>
	<w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
	<a:MessageID>uuid:4BFD28B2-CB5D-43CC-B086-4D2C3572314C</a:MessageID>
	<w:Locale xml:lang="en-US" mustUnderstand="false"/>
	<p:DataLocale xml:lang="en-US" mustUnderstand="false"/>
	<w:OperationTimeout>PT60S</w:OperationTimeout>
	<w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
	<a:Action mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete</a:Action>
	<w:SelectorSet>
	  <w:Selector Name="ShellId">%s</w:Selector>
	</w:SelectorSet>
  </env:Header>
  <env:Body/>
</env:Envelope>`, s.Id)

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewBufferString(command))
	request.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")
	request.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(s.Owner+":"+s.password)))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return NewHttpError(response)
	}

	return nil
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
				fmt.Println(h.Action.URI)
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
