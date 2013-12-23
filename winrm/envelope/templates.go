package envelope

import (
	"bytes"
	"text/template"
)

type CreateShell struct {
	MessageId string
}

type DeleteShell struct {
	MessageId string
	ShellId   string
}

type CreateCommand struct {
	MessageId   string
	ShellId     string
	CommandText string
}

type Receive struct {
	MessageId string
	ShellId   string
	CommandId string
}

func (m *CreateShell) Xml() string {
	t := template.Must(template.New("CreateShell").Parse(CreateShellTemplate))
	return applyTemplate(t, m)
}

func (m *DeleteShell) Xml() string {
	t := template.Must(template.New("DeleteShell").Parse(DeleteShellTemplate))
	return applyTemplate(t, m)
}

func (m *CreateCommand) Xml() string {
	t := template.Must(template.New("CreateCommand").Parse(CreateCommandTemplate))
	return applyTemplate(t, m)
}

func (m *Receive) Xml() string {
	t := template.Must(template.New("Receive").Parse(ReceiveTemplate))
	return applyTemplate(t, m)
}

func applyTemplate(t *template.Template, data interface{}) string {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		panic(err)
	}
	return b.String()
}

const CreateShellTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="` + NS_SOAP_ENV + `" xmlns:a="` + NS_ADDRESSING + `" xmlns:rsp="` + NS_WIN_SHELL + `" xmlns:w="` + NS_WSMAN_DMTF + `">
  <env:Header>
    <a:To>http://localhost:5985/wsman</a:To>
    <a:ReplyTo>
      <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
    <a:MessageID>uuid:{{.MessageId}}</a:MessageID>
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

const DeleteShellTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="` + NS_SOAP_ENV + `" xmlns:a="` + NS_ADDRESSING + `" xmlns:w="` + NS_WSMAN_DMTF + `">
  <env:Header>
    <a:To>http://localhost:5985/wsman</a:To>
    <a:ReplyTo>
      <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
    <a:MessageID>uuid:{{.MessageId}}</a:MessageID>
    <w:OperationTimeout>PT60S</w:OperationTimeout>
    <w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
    <a:Action mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/09/transfer/Delete</a:Action>
    <w:SelectorSet>
      <w:Selector Name="ShellId">{{.ShellId}}</w:Selector>
    </w:SelectorSet>
  </env:Header>
  <env:Body/>
</env:Envelope>`

const CreateCommandTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="` + NS_SOAP_ENV + `" xmlns:a="` + NS_ADDRESSING + `" xmlns:rsp="` + NS_WIN_SHELL + `" xmlns:w="` + NS_WSMAN_DMTF + `">
  <env:Header>
    <a:To>http://localhost:5985/wsman</a:To>
    <a:ReplyTo>
      <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
    </a:ReplyTo>
    <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
    <a:MessageID>uuid:{{.MessageId}}</a:MessageID>
    <w:OperationTimeout>PT60S</w:OperationTimeout>
    <w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
    <a:Action mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Command</a:Action>
    <w:OptionSet>
      <w:Option Name="WINRS_CONSOLEMODE_STDIN">TRUE</w:Option>
      <w:Option Name="WINRS_SKIP_CMD_SHELL">FALSE</w:Option>
    </w:OptionSet>
    <w:SelectorSet>
      <w:Selector Name="ShellId">{{.ShellId}}</w:Selector>
    </w:SelectorSet>
  </env:Header>
  <env:Body>
    <rsp:CommandLine>
      <rsp:Command>{{.CommandText}}</rsp:Command>
    </rsp:CommandLine>
  </env:Body>
</env:Envelope>`

const ReceiveTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<env:Envelope xmlns:env="` + NS_SOAP_ENV + `" xmlns:a="` + NS_ADDRESSING + `" xmlns:rsp="` + NS_WIN_SHELL + `" xmlns:w="` + NS_WSMAN_DMTF + `">
    <env:Header>
      <a:To>http://localhost:5985/wsman</a:To>
      <a:ReplyTo>
        <a:Address mustUnderstand="true">http://schemas.xmlsoap.org/ws/2004/08/addressing/role/anonymous</a:Address>
      </a:ReplyTo>
      <w:MaxEnvelopeSize mustUnderstand="true">153600</w:MaxEnvelopeSize>
      <a:MessageID>uuid:{{.MessageId}}</a:MessageID>
      <w:OperationTimeout>PT60S</w:OperationTimeout>
      <w:ResourceURI mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/cmd</w:ResourceURI>
      <a:Action mustUnderstand="true">http://schemas.microsoft.com/wbem/wsman/1/windows/shell/Receive</a:Action>
      <w:SelectorSet>
        <w:Selector Name="ShellId">{{.ShellId}}</w:Selector>
      </w:SelectorSet>
    </env:Header>
    <env:Body>
      <rsp:Receive>
        <rsp:DesiredStream CommandId="{{.CommandId}}">stdout stderr</rsp:DesiredStream>
      </rsp:Receive>
    </env:Body>
  </env:Envelope>`
