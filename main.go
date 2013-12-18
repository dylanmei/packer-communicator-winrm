package main

import (
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"log"
	"os"
)

const endpoint = "http://localhost:5985/wsman"

// SOLO USAGE: ./packer-communicator-winrm shell -user vagrant -pass vagrant command-text

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		plugin()
		return
	}

	Run(&shell{
		Handle: func(user, pass string, commands ...string) {
			shell, err := winrm.NewShell(endpoint, user, pass)
			if err != nil {
				log.Fatal(err.Error())
			}

			defer shell.Delete()
			log.Println("Shell:", shell.Id)

			command, err := shell.NewCommand(commands[0])
			if err != nil {
				log.Println(err)
				return
			}

			output, err := command.Receive()
			if err != nil {
				log.Println(err)
				return
			}

			log.Printf("Command: %s, ExitCode: %d", command.CommandText, output.ExitCode)
			for _, value := range output.Stdout {
				log.Println("stdout", value)
			}
			for _, value := range output.Stderr {
				log.Println("stderr", value)
			}
		},
	})
}
