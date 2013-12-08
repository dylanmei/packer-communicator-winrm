package envelope

import (
	"encoding/xml"
	"github.com/dylanmei/packer-communicator-winrm/winrm/addressing"
	"github.com/dylanmei/packer-communicator-winrm/winrm/rsp"
	"github.com/dylanmei/packer-communicator-winrm/winrm/wsman"
)

type Envelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  *Header  `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	Body    *Body    `xml:"http://www.w3.org/2003/05/soap-envelope Body"`
}

type Header struct {
	MessageID        string                        `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing MessageID"`
	Action           *addressing.AttributedURI     `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing Action"`
	ReplyTo          *addressing.EndpointReference `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing ReplyTo"`
	To               *addressing.AttributedURI     `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing To"`
	MaxEnvelopeSize  *wsman.MaxEnvelopeSize        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd MaxEnvelopeSize"`
	OperationTimeout string                        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd OperationTimeout"`
	ResourceURI      *wsman.AttributableURI        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd ResourceURI"`
	OptionSet        *wsman.OptionSet              `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd OptionSet,omitempty"`
	SelectorSet      *wsman.SelectorSet            `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd SelectorSet,omitempty"`
}

type Body struct {
	Shell *rsp.Shell `xml:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell Shell,omitempty"`
}
