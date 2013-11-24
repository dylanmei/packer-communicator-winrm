package main

import (
	"flag"
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"log"
	"os"
)

func solo() bool {
	args := os.Args[1:]
	if args := os.Args[1:]; len(args) == 0 {
		return false
	}

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	user := flags.String("user", "vagrant", "WinRM user to run as")
	pass := flags.String("pass", "vagrant", "WinRM password for user")

	if args[0] == "shell" {
		flags.Parse(os.Args[2:])
		shell(*user, *pass, flags.Args())
	}
	return true
}

func shell(user, pass string, commands []string) {
	s, err := winrm.NewShell("vagrant", "vagrant")
	if err != nil {
		log.Fatal(err.Error())
	}

	defer s.Delete()

	fmt.Println("Shell Id: ", s.Id)
	//todo s.Execute(commands...)
}
