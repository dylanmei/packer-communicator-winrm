package main

import (
	"encoding/base64"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
)

func powershell(shell *winrm.Shell, command string) (string, error) {
	var bytes []byte
	for _, c := range []byte(command) {
		bytes = append(bytes, c, 0)
	}

	text := "powershell -NoProfile -EncodedCommand " +
		base64.StdEncoding.EncodeToString(bytes)

	return cmd(shell, text)
}
