package main

import (
	"github.com/mitchellh/packer/packer"
	rpc "github.com/mitchellh/packer/packer/plugin"
	"log"
	"os"
)

const endpoint = "http://localhost:5985/wsman"

// SOLO USAGE: ./packer-communicator-winrm cmd -user vagrant -pass vagrant command-text
// set WINRM_DEBUG=1 for more output

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		plugin()
		return
	}

	Run(&cmd{
		Handle: func(user, pass, command string) {
			communicator := &Communicator{endpoint, user, pass}
			rc := &packer.RemoteCmd{
				Command: command,
				Stdout:  os.Stdout,
				Stderr:  os.Stderr,
			}

			err := communicator.Start(rc)
			if err != nil {
				log.Printf("unable to run command: %s", err)
				return
			}

			rc.Wait()
		},
	})
}

func plugin() {
	server, err := rpc.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterCommunicator(new(Communicator))
	server.Serve()
}
