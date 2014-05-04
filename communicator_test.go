package main

import (
	"bytes"
	"github.com/masterzen/winrm/winrmtest"
	"github.com/mitchellh/packer/packer"
	"io"
	"testing"
)

func Test_running_a_command(t *testing.T) {
	h := winrmtest.NewHost()
	defer h.Close()

	h.Command("dir C:\\Temp", func(out, err io.Writer) int {
		out.Write([]byte("test"))
		return 0
	})

	comm := &Communicator{*host, *user, *pass}
	rc := &packer.RemoteCmd{
		Command: "dir C:\\Temp",
		Stdout:  bytes.NewBuffer(make([]byte, 0)),
		Stderr:  bytes.NewBuffer(make([]byte, 0)),
	}

	err := comm.Start(rc)
	rc.Wait()

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}
	if rc.ExitStatus != 0 {
		t.Errorf("Expected ExitStatus 0, found %d", rc.ExitStatus)
	}
}
