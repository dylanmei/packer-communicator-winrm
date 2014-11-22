package main

import (
	"flag"
	"github.com/mefellows/winrm/winrm"
	"github.com/mitchellh/packer/packer"
	rpc "github.com/mitchellh/packer/packer/plugin"
	"github.com/rakyll/command"
	"log"
	"os"
)

var host = flag.String("host", "localhost", "host machine")
var port = flag.Int("port", 5985, "host port")
var user = flag.String("user", "vagrant", "user to run as")
var pass = flag.String("pass", "vagrant", "user's password")

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
	command.On("cmd", "run a command", &RunCommand{}, []string{})
	command.On("file", "copy a file", &FileCommand{}, []string{})
	command.Parse()
	command.Run()
}

type RunCommand struct{}

func (r *RunCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	return fs
}

func (r *RunCommand) Run(args []string) {
	command := args[0]
	// TODO: set the correct timeout or defaults or Config it
	communicator, _ := New(&winrm.Endpoint{*host, *port}, *user, *pass, 500)
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

type FileCommand struct {
	to   *string
	from *string
}

func (f *FileCommand) Flags(fs *flag.FlagSet) *flag.FlagSet {
	f.to = fs.String("to", "", "destination file path")
	f.from = fs.String("from", "", "source file path")
	return fs
}

func (f *FileCommand) Run(args []string) {
	communicator, _ := New(&winrm.Endpoint{*host, *port}, *user, *pass, 500)

	_, err := os.Stat(*f.from)
	if err != nil {
		log.Panicln("unable to stat file", err.Error())
	}

	file, err := os.Open(*f.from)
	if err != nil {
		log.Panicln("unable to open file", err.Error())
	}

	err = communicator.Upload(*f.to, file, nil)
	if err != nil {
		log.Printf("unable to copy file: %s", err)
	}
}
