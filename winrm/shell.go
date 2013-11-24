package winrm

// http://msdn.microsoft.com/en-us/library/cc251731.aspx

import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

type Shell struct {
	Id       string
	Owner    string
	password string
}

func NewShell(owner, pass string) (*Shell, error) {
	command := `
<env:Envelope xmlns:xsd="http://www.w3.org/2001/XMLSchema"
	xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
	xmlns:env="http://www.w3.org/2003/05/soap-envelope"
	xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing"
	xmlns:b="http://schemas.dmtf.org/wbem/wsman/1/cimbinding.xsd"
	xmlns:n="http://schemas.xmlsoap.org/ws/2004/09/enumeration"
	xmlns:x="http://schemas.xmlsoap.org/ws/2004/09/transfer"
	xmlns:w="http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd"
	xmlns:p="http://schemas.microsoft.com/wbem/wsman/1/wsman.xsd"
	xmlns:rsp="http://schemas.microsoft.com/wbem/wsman/1/windows/shell"
	xmlns:cfg="http://schemas.microsoft.com/wbem/wsman/1/config"
>
  <env:Header>
	<a:To>http://localhost:5985/wsman</a:To>
	<a:ReplyTo>
	  <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
	</a:ReplyTo>
	<w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
	<a:MessageID>uuid:E266B619-7457-4B69-AEAB-633E5E36017A</a:MessageID>
	<w:Locale xml:lang="en-US" mustUnderstand="false"/>
	<p:DataLocale xml:lang="en-US" mustUnderstand="false"/>
	<w:OperationTimeout>PT60S</w:OperationTimeout>
	<w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
	<a:Action mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/09/transfer/Create</a:Action>
	<w:OptionSet>
	  <w:Option Name="WINRS_NOPROFILE">FALSE</w:Option>
	  <w:Option Name="WINRS_CODEPAGE">437</w:Option>
	</w:OptionSet>
  </env:Header>
  <env:Body>
	<rsp:Shell>
	  <rsp:InputStreams>stdin</rsp:InputStreams>
	  <rsp:OutputStreams>stdout stderr</rsp:OutputStreams>
	</rsp:Shell>
  </env:Body>
</env:Envelope>`

	request, _ := http.NewRequest("POST",
		"http://localhost:5985/wsman", bytes.NewBufferString(command))
	request.Header.Add("Content-Type", "application/soap+xml;charset=UTF-8")
	request.Header.Add("Authorization", "Basic "+
		base64.StdEncoding.EncodeToString([]byte(owner+":"+pass)))

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	fmt.Println("HTTP Status", response.Status)

	for key, value := range response.Header {
		fmt.Println(" ", key, ":", value)
	}

	if response.StatusCode != 200 {
		return nil, errors.New(response.Status)
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	rs := decodeResponse(bytes.NewBuffer(body))

	return &Shell{rs.ShellId, owner, pass}, nil
}

func (s *Shell) Delete() {
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
		log.Fatal(err)
	}

	fmt.Println("HTTP Status", response.Status)

	for key, value := range response.Header {
		fmt.Println(" ", key, ":", value)
	}
}

type RemoteShell struct {
	ShellId string `xml:"ShellId"`
}

func decodeResponse(reader io.Reader) *RemoteShell {
	decoder := xml.NewDecoder(reader)

	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}

		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Space == NS_WIN_SHELL && se.Name.Local == "Shell" {
				var rs RemoteShell
				decoder.DecodeElement(&rs, &se)
				return &rs
			}
		}
	}

	return nil
}
