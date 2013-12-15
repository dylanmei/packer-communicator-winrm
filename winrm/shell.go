package winrm

import (
	"errors"
	"github.com/dylanmei/packer-communicator-winrm/envelope"
	"github.com/mitchellh/packer/common/uuid"
	"launchpad.net/xmlpath"
	"log"
)

type Shell struct {
	Id       string
	Owner    string
	password string
}

type Command struct {
	Id    string
	shell *Shell
}

func NewShell(user, pass string) (*Shell, error) {
	env := &envelope.CreateShell{uuid.TimeOrderedUUID()}
	response, err := deliver(user, pass, env)
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("//Body/Shell/ShellId")
	root, err := xmlpath.Parse(response)
	if err != nil {
		log.Fatal(err)
	}

	shellId, ok := path.String(root)
	if !ok {
		return nil, errors.New("Could not create shell.")
	}

	return &Shell{shellId, user, pass}, nil
}

func (s *Shell) NewCommand(cmd string) (*Command, error) {
	env := &envelope.CreateCommand{uuid.TimeOrderedUUID(), s.Id, cmd}
	response, err := deliver(s.Owner, s.password, env)
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("//Body/CommandResponse/CommandId")
	root, err := xmlpath.Parse(response)
	if err != nil {
		return nil, err
	}

	commandId, ok := path.String(root)
	if !ok {
		return nil, errors.New("Could not create command.")
	}

	return &Command{commandId, s}, nil
}

func (s *Shell) Delete() error {
	env := &envelope.DeleteShell{uuid.TimeOrderedUUID(), s.Id}
	_, err := deliver(s.Owner, s.password, env)

	if err != nil {
		log.Println(err.Error())
	}
	return err
}
