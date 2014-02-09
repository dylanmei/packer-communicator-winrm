package main

import (
	"bytes"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"log"
	"strings"
)

func cmd(shell *winrm.Shell, command string) (string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	shell.Stdout = stdout
	shell.Stderr = stderr

	log.Println("starting winrm command:", command)

	c, err := shell.NewCommand(command)

	if err != nil {
		return "", err
	}

	if _, err = c.Receive(); err != nil {
		return "", err
	}

	if stderr.Len() > 0 {
		log.Println("winrm stderr: %s", stderr.String())
	}

	return strings.Trim(stdout.String(), " \r\n"), nil
}
