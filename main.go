package main

import (
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"log"
	"os"
)

// SOLO USAGE: ./packer-communicator-winrm shell -user vagrant -pass vagrant command-text

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		plugin()
		return
	}

	Run(&shell{
		Handle: func(user, pass string, commands ...string) {
			shell, err := winrm.NewShell(user, pass)
			if err != nil {
				log.Fatal(err.Error())
			}

			defer shell.Delete()
			log.Println("shell:", shell.Id)
		},
	})
}
