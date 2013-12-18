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
	Endpoint string
	Owner    string
	password string
}

func NewShell(endpoint, user, pass string) (*Shell, error) {
	env := &envelope.CreateShell{uuid.TimeOrderedUUID()}
	response, err := deliver(endpoint, user, pass, env)
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("//Body/Shell/ShellId")
	root, err := xmlpath.Parse(response)
	if err != nil {
		return nil, err
	}

	id, ok := path.String(root)
	if !ok {
		return nil, errors.New("Could not create shell.")
	}

	return &Shell{id, endpoint, user, pass}, nil
}

func (s *Shell) NewCommand(cmd string) (*Command, error) {
	env := &envelope.CreateCommand{uuid.TimeOrderedUUID(), s.Id, cmd}
	response, err := deliver(s.Endpoint, s.Owner, s.password, env)
	if err != nil {
		return nil, err
	}

	path := xmlpath.MustCompile("//Body/CommandResponse/CommandId")
	root, err := xmlpath.Parse(response)
	if err != nil {
		return nil, err
	}

	id, ok := path.String(root)
	if !ok {
		return nil, errors.New("Could not create command.")
	}

	return &Command{s, id, cmd}, nil
}

func (s *Shell) Delete() error {
	env := &envelope.DeleteShell{uuid.TimeOrderedUUID(), s.Id}
	_, err := deliver(s.Endpoint, s.Owner, s.password, env)

	if err != nil {
		log.Println(err.Error())
	}
	return err
}
