package wsman

type MaxEnvelopeSize struct {
	Value          uint `xml:",chardata"`
	MustUnderstand bool `xml:"mustUnderstand,attr,omitempty"`
}

type AttributableURI struct {
	URI            string `xml:",chardata"`
	MustUnderstand bool   `xml:"mustUnderstand,attr,omitempty"`
}

type OptionSet struct {
	Options []*Option `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd Option"`
}

type Option struct {
	Name  string `xml:",attr"`
	Value string `xml:",chardata"`
}

type SelectorSet struct {
	Selectors []*Selector `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd Selector"`
}

type Selector struct {
	Name  string `xml:",attr"`
	Value string `xml:",chardata"`
}
