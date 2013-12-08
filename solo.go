package main

import (
	"encoding/xml"
	_ "flag"
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/dylanmei/packer-communicator-winrm/winrm/addressing"
	"github.com/dylanmei/packer-communicator-winrm/winrm/wsman"
	"github.com/mitchellh/packer/common/uuid"
	"log"
	"os"
)

func solo() bool {
	// args := os.Args[1:]
	// if args := os.Args[1:]; len(args) == 0 {
	// 	return false
	// }

	encodeSomething()

	// flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	// user := flags.String("user", "vagrant", "WinRM user to run as")
	// pass := flags.String("pass", "vagrant", "WinRM password for user")

	// if args[0] == "shell" {
	// 	flags.Parse(os.Args[2:])
	// 	shell(*user, *pass, flags.Args())
	// }

	return true
}

func shell(user, pass string, commands []string) {
	s, err := winrm.NewShell(user, pass)
	if err != nil {
		log.Fatal(err.Error())
	}

	defer s.Delete()

	fmt.Println("Shell Id: ", s.Id)
	//todo s.Execute(commands...)
}

type SoapHeader struct {
	XMLName          xml.Name                      `xml:"http://www.w3.org/2003/05/soap-envelope Header"`
	MessageID        string                        `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing MessageID"`
	Action           *addressing.AttributedURI     `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing Action"`
	ReplyTo          *addressing.EndpointReference `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing ReplyTo"`
	To               *addressing.AttributedURI     `xml:"http://schemas.xmlsoap.org/ws/2004/08/addressing To"`
	MaxEnvelopeSize  *wsman.MaxEnvelopeSize        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd MaxEnvelopeSize"`
	OperationTimeout string                        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd OperationTimeout"`
	ResourceURI      *wsman.AttributableURI        `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd ResourceURI"`
	OptionSet        *wsman.OptionSet              `xml:"http://schemas.dmtf.org/wbem/wsman/1/wsman.xsd OptionSet"`
}

func encodeSomething() {

	sh := &SoapHeader{
		MessageID: "uuid:" + uuid.TimeOrderedUUID(),
		Action: &addressing.AttributedURI{
			URI:            "http://schemas.xmlsoap.org/ws/2004/09/transfer/Create",
			MustUnderstand: true,
		},
		To: &addressing.AttributedURI{
			URI: "http://localhost:5985/wsman",
			//			MustUnderstand: true,
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

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent(" ", "	")
	if err := enc.Encode(sh); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Println()
}
