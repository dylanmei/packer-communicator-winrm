package main

import (
	"flag"
	"github.com/mitchellh/packer/packer"
	rpc "github.com/mitchellh/packer/packer/plugin"
	"github.com/rakyll/command"
	"log"
	"os"
)

const endpoint = "http://localhost:5985/wsman"

var user = flag.String("user", "vagrant", "user to run as")
var pass = flag.String("pass", "vagrant", "user's password")

// SOLO USAGE:
//   ./packer-communicator-winrm help
//   ./packer-communicator-winrm -user=vagrant -pass=vagrant run command-text
// Set WINRM_DEBUG=1 for more output

func main() {
	args := os.Args[1:]
	if len(args) != 0 {
		standalone()
	} else {
		server, err := rpc.Server()
		if err != nil {
			panic(err)
		}
		server.RegisterCommunicator(new(Communicator))
		server.Serve()
	}
}

func standalone() {
	command.On("run", "run a command", &RunCommand{})
	command.Parse()
	command.Run()
}

type RunCommand struct{}

func (cmd *RunCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (r *RunCommand) Run(args []string) {
	command := args[0]
	communicator := &Communicator{endpoint, *user, *pass}
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
}
