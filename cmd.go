package main

import (
	"bytes"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"strings"
)

func cmd(shell *winrm.Shell, command string) (string, error) {
	c, err := shell.NewCommand(command)

	if err != nil {
		return "", err
	}

	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	shell.Stdout = stdout
	shell.Stderr = stderr

	println("[CMD]")
	println(command)

	if _, err = c.Receive(); err != nil {
		return "", err
	}

	if stdout.Len() > 0 {
		println("[STDOUT]")
		println(stdout.String())
	}

	if stderr.Len() > 0 {
		println("[STDERR]")
		println(stderr.String())
	}

	return strings.Trim(stdout.String(), " \r\n"), nil
}
