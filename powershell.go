package main

import (
	"bytes"
	"encoding/base64"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"log"
	"strings"
)

func powershell(shell *winrm.Shell, command string) (string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	shell.Stdout = stdout
	shell.Stderr = stderr

	var bytes []byte
	for _, c := range []byte(command) {
		bytes = append(bytes, c, 0)
	}

	log.Println("starting winrm command: powershell -Command", command)

	text := "powershell -NoProfile -EncodedCommand " +
		base64.StdEncoding.EncodeToString(bytes)
	c, err := shell.NewCommand(text)

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
