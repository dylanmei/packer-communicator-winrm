package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/dylanmei/winrmtest"
	"github.com/mitchellh/packer/packer"
)

func Test_running_a_command(t *testing.T) {
	r := winrmtest.NewRemote()
	defer r.Close()

	r.CommandFunc("echo tacos", func(out, err io.Writer) int {
		out.Write([]byte("tacos"))
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
		host: r.Host,
		port: r.Port,
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
