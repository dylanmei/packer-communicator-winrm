package main

import (
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	if !solo() {
		plugin.ServeCommunicator(new(winrm.Communicator))
	}
}
