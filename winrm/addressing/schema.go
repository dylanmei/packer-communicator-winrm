package addressing

type EndpointReference struct {
	Address *AttributedURI `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing Address"`
}

type AttributedURI struct {
	URI            string `xml:",chardata"`
	MustUnderstand bool   `xml:"mustUnderstand,attr,omitempty"`
}
