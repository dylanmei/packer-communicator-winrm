package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/masterzen/winrm/winrm"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
	"io"
	"strings"
)

// A Communicator is the interface used to communicate with the machine
// that exists that will eventually be packaged into an image. Communicators
// allow you to execute remote commands, upload files, etc.
//
// Communicators must be safe for concurrency, meaning multiple calls to
// Start or any other method may be called at the same time.
type Communicator struct {
	host string
	user string
	pass string
}

// Start takes a RemoteCmd and starts it. The RemoteCmd must not be
// modified after being used with Start, and it must not be used with
// Start again. The Start method returns immediately once the command
// is started. It does not wait for the command to complete. The
// RemoteCmd.Exited field should be used for this.
func (c *Communicator) Start(rc *packer.RemoteCmd) (err error) {

	client := winrm.NewClient(c.host, c.user, c.pass)
	shell, err := client.CreateShell()
	if err != nil {
		return
	}
	defer shell.Close()

	cmd, err := shell.Execute(rc.Command)
	if err != nil {
		return
	}

	go io.Copy(rc.Stdout, cmd.Stdout)
	go io.Copy(rc.Stderr, cmd.Stderr)

	cmd.Wait()
	rc.SetExited(cmd.ExitCode())
	return
}

// Upload uploads a file to the machine to the given path with the
// contents coming from the given reader. This method will block until
// it completes.
func (c *Communicator) Upload(path string, r io.Reader) (err error) {

	client := winrm.NewClient(c.host, c.user, c.pass)
	shell, err := client.CreateShell()
	if err != nil {
		return
	}
	defer shell.Close()

	temp, err := runCommand(client, fmt.Sprintf(
		`echo %%TEMP%%\packer-%s.tmp`, uuid.TimeOrderedUUID()))

	if err != nil {
		return
	}

	temp = strings.TrimSpace(temp)
	bytes := make([]byte, 8000-len(temp))
	for {
		read, _ := r.Read(bytes)
		if read == 0 {
			break
		}

		_, err = runCommand(client, fmt.Sprintf(
			`echo %s >> %s`, base64.StdEncoding.EncodeToString(bytes[:read]), temp))

		if err != nil {
			return
		}
	}

	_, err = runPowershell(client, fmt.Sprintf(`
		$path = "%s"
		$temp = "%s"

		$dir = [System.IO.Path]::GetDirectoryName($path)
		if (-Not (Test-Path $dir)) {
			mkdir $dir
		} elseif (Test-Path $path) {
			rm $path
		}

		$lines = Get-Content $temp
		$value = [string]::join("",$lines)
		$bytes = [System.Convert]::FromBase64String($value)

		$file = [System.IO.Path]::GetFullPath($path)
		[System.IO.File]::WriteAllBytes($file, $bytes)
 
		del $temp
	`, path, temp))
	return
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
	panic("not implemented yet")
}

// Download downloads a file from the machine from the given remote path
// with the contents writing to the given writer. This method will
// block until it completes.
func (c *Communicator) Download(path string, w io.Writer) error {
	panic("not implemented yet")
}

func runCommand(client *winrm.Client, text string) (string, error) {
	stdout, stderr, err := client.RunWithString(text, "")

	if err != nil {
		return "", err
	}

	if stderr != "" {
		return "", errors.New("Error running command on guest: " + stderr)
	}

	return stdout, nil
}

func runPowershell(client *winrm.Client, text string) (string, error) {
	var bytes []byte
	for _, c := range []byte(text) {
		bytes = append(bytes, c, 0)
	}

	encoded := "powershell -NoProfile -EncodedCommand " +
		base64.StdEncoding.EncodeToString(bytes)

	return runCommand(client, encoded)
}
