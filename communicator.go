package main

import (
	"errors"
	"github.com/mitchellh/packer/packer"
	"io"
)

// A Communicator is the interface used to communicate with the machine
// that exists that will eventually be packaged into an image. Communicators
// allow you to execute remote commands, upload files, etc.
//
// Communicators must be safe for concurrency, meaning multiple calls to
// Start or any other method may be called at the same time.
type Communicator struct {
}

// Start takes a RemoteCmd and starts it. The RemoteCmd must not be
// modified after being used with Start, and it must not be used with
// Start again. The Start method returns immediately once the command
// is started. It does not wait for the command to complete. The
// RemoteCmd.Exited field should be used for this.
func (c *Communicator) Start(rc *packer.RemoteCmd) error {
	return errors.New("Not implemented")
}

// Upload uploads a file to the machine to the given path with the
// contents coming from the given reader. This method will block until
// it completes.
func (c *Communicator) Upload(path string, r io.Reader) error {
	return errors.New("Not implemented")
}

// UploadDir uploads the contents of a directory recursively to
// the remote path. It also takes an optional slice of paths to
// ignore when uploading.
//
// The folder name of the source folder should be created unless there
// is a trailing slash on the source "/". For example: "/tmp/src" as
// the source will create a "src" directory in the destination unless
// a trailing slash is added. This is identical behavior to rsync(1).
func (c *Communicator) UploadDir(dst string, src string, exclude []string) error {
	return errors.New("Not implemented")
}

// Download downloads a file from the machine from the given remote path
// with the contents writing to the given writer. This method will
// block until it completes.
func (c *Communicator) Download(path string, w io.Writer) error {
	return errors.New("Not implemented")
}
