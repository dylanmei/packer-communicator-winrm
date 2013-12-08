package rsp

type Shell struct {
	InputStreams  string `xml:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell InputStreams,omitempty"`
	OutputStreams string `xml:"http://schemas.microsoft.com/wbem/wsman/1/windows/shell OutputStreams,omitempty"`
}
