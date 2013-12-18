package winrm

import (
	"encoding/base64"
	"errors"
	"github.com/dylanmei/packer-communicator-winrm/envelope"
	"github.com/mitchellh/packer/common/uuid"
	"launchpad.net/xmlpath"
	"strconv"
	"strings"
)

type Command struct {
	shell       *Shell
	Id          string
	CommandText string
}

type CommandOutput struct {
	ExitCode int
	Stdout   []string
	Stderr   []string
}

func (c *Command) Receive() (*CommandOutput, error) {
	env := &envelope.Receive{uuid.TimeOrderedUUID(), c.shell.Id, c.Id}
	response, err := deliver(c.shell.Endpoint, c.shell.Owner, c.shell.password, env)
	if err != nil {
		return nil, err
	}

	state := xmlpath.MustCompile("//Body/ReceiveResponse/CommandState/@State")
	exitcode := xmlpath.MustCompile("//Body/ReceiveResponse/CommandState/ExitCode")
	stdout := xmlpath.MustCompile("//Body/ReceiveResponse/Stream[@Name='stdout']")
	stderr := xmlpath.MustCompile("//Body/ReceiveResponse/Stream[@Name='stderr']")

	root, err := xmlpath.Parse(response)
	if err != nil {
		return nil, err
	}

	value, ok := state.String(root)
	if !ok {
		return nil, errors.New("Could not discover command state")
	}

	if !strings.HasSuffix(value, "Done") {
		panic("TODO: appending output")
	}

	value, ok = exitcode.String(root)
	if !ok {
		return nil, errors.New("Expected an exit code")
	}

	result, _ := strconv.Atoi(value)
	output := &CommandOutput{
		ExitCode: result,
		Stdout:   collectStream(root, stdout),
		Stderr:   collectStream(root, stderr),
	}

	return output, nil
}

func collectStream(node *xmlpath.Node, path *xmlpath.Path) []string {
	iter := path.Iter(node)
	values := make([]string, 0)

	for iter.Next() {
		node := iter.Node()
		data := node.String()
		if len(data) > 0 {
			b, _ := base64.StdEncoding.DecodeString(node.String())
			values = append(values, string(b))
		}
	}

	return values
}
