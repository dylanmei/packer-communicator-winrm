package main

import (
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/mitchellh/packer/packer/plugin"
)

// SOLO USAGE: ./packer-communicator-winrm shell -user vagrant -pass vagrant
// PLUGIN: ServeCommunicator doesn't exist. Fork https://github.com/dylanmei/packer/tree/communicator_plugin

func main() {
	if !solo() {
		server, err := plugin.Server()
		if err != nil {
			panic(err)
		}
		server.RegisterCommunicator(new(winrm.Communicator))
		server.Serve()
	}
}
