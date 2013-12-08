package main

import (
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/mitchellh/packer/packer/plugin"
)

// SOLO USAGE: ./packer-communicator-winrm shell -user vagrant -pass vagrant
// PLUGIN: ServeCommunicator doesn't exist. Fork https://github.com/dylanmei/packer/tree/communicator_plugin

func main() {
	if !solo() {
		plugin.ServeCommunicator(new(winrm.Communicator))
	}
}
