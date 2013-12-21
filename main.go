package main

import (
	"github.com/mitchellh/packer/packer"
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

			communicator.Start(rc)
			rc.Wait()
		},
	})
}
