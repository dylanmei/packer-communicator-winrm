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

	h.CommandFunc("echo tacos", func(out, err io.Writer) int {
		out.Write([]byte("test"))
		return 0
	})

	stdout := bytes.NewBuffer(make([]byte, 0))
	stderr := bytes.NewBuffer(make([]byte, 0))
	rc := &packer.RemoteCmd{
		Command: "echo tacos",
		Stdout:  stdout,
		Stderr:  stderr,
	}

	comm := &Communicator{
		host: h.Hostname,
		port: h.Port,
	}

	err := comm.Start(rc)

	if err != nil {
		t.Errorf("Unexpected error %v", err)
	}

	rc.Wait()

	if rc.ExitStatus != 0 {
		t.Errorf(`expected rc.ExitStatus=0 but was %d`, rc.ExitStatus)
	}

	if stderr.String() != "" {
		t.Errorf(`expected sterr="" but was "%s"`, stderr)
	}

	if stdout.String() != "tacos" {
		t.Errorf(`expected stdout=tacos but was "%s"`, stdout)
	}
}
