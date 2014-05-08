package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/masterzen/winrm/winrm"
	"github.com/mitchellh/packer/common/uuid"
	"github.com/mitchellh/packer/packer"
)

// A Communicator is the interface used to communicate with the machine
// that exists that will eventually be packaged into an image. Communicators
// allow you to execute remote commands, upload files, etc.
//
// Communicators must be safe for concurrency, meaning multiple calls to
// Start or any other method may be called at the same time.
type Communicator struct {
	host string
	port int
	user string
	pass string
}

// Start takes a RemoteCmd and starts it. The RemoteCmd must not be
// modified after being used with Start, and it must not be used with
// Start again. The Start method returns immediately once the command
// is started. It does not wait for the command to complete. The
// RemoteCmd.Exited field should be used for this.
func (c *Communicator) Start(rc *packer.RemoteCmd) error {
	client := winrm.NewClient(&winrm.Endpoint{c.host, c.port}, c.user, c.pass)
	shell, err := client.CreateShell()
	if err != nil {
		return err
	}
	defer shell.Close()

	cmd, err := shell.Execute(rc.Command)
	if err != nil {
		return err
	}

	//	go func() {
	go io.Copy(rc.Stdout, cmd.Stdout)
	go io.Copy(rc.Stderr, cmd.Stderr)

	cmd.Wait()
	rc.SetExited(cmd.ExitCode())
	//	}()

	return nil
}

// Upload uploads a file to the machine to the given path with the
// contents coming from the given reader. This method will block until
// it completes.
func (c *Communicator) Upload(path string, r io.Reader) (err error) {

	client := winrm.NewClient(&winrm.Endpoint{c.host, c.port}, c.user, c.pass)
	shell, err := client.CreateShell()
	if err != nil {
		return
	}
	defer shell.Close()

	temp, err := runCommand(shell, fmt.Sprintf(
		`echo %%TEMP%%\packer-%s.tmp`, uuid.TimeOrderedUUID()))

	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	temp = strings.TrimSpace(temp)
	for _, chunk := range encodeChunks(bytes, 8000-len(temp)) {

		_, err = runCommand(shell,
			fmt.Sprintf(`echo %s >> %s`, chunk, temp))

		if err != nil {
			return
		}
	}

	_, err = runPowershell(shell, fmt.Sprintf(`
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

func runCommand(shell *winrm.Shell, text string) (string, error) {
	cmd, err := shell.Execute(text)

	if err != nil {
		return "", err
	}

	var stdout, stderr bytes.Buffer
	go io.Copy(&stdout, cmd.Stdout)
	go io.Copy(&stderr, cmd.Stderr)

	cmd.Wait()

	if stderr.Len() > 0 {
		return "", errors.New("Error running command on guest: " + stderr.String())
	}

	return stdout.String(), nil
}

func runPowershell(shell *winrm.Shell, text string) (string, error) {
	var bytes []byte
	for _, c := range []byte(text) {
		bytes = append(bytes, c, 0)
	}

	encoded := "powershell -NoProfile -EncodedCommand " +
		base64.StdEncoding.EncodeToString(bytes)

	return runCommand(shell, encoded)
}

func encodeChunks(bytes []byte, chunkSize int) []string {
	text := base64.StdEncoding.EncodeToString(bytes)
	reader := strings.NewReader(text)

	var chunks []string
	chunk := make([]byte, chunkSize)

	for {
		n, _ := reader.Read(chunk)
		if n == 0 {
			break
		}

		chunks = append(chunks, string(chunk[:n]))
	}

	return chunks
}
