package main

import (
	"encoding/base64"
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/mitchellh/packer/common/uuid"
	"io"
	"log"
)

func upload(shell *winrm.Shell, path string, r io.Reader) (err error) {

	temp, err := cmd(shell, fmt.Sprintf(
		`echo %%TEMP%%\packer-%s.tmp`, uuid.TimeOrderedUUID()))

	if err != nil {
		return
	}

	log.Printf("transfering file to", temp)

	bytes := make([]byte, 8000-len(temp))
	for {
		read, _ := r.Read(bytes)
		if read == 0 {
			break
		}
		_, err = cmd(shell, fmt.Sprintf(
			`echo %s >> %s`, base64.StdEncoding.EncodeToString(bytes[:read]), temp))

		if err != nil {
			return
		}
	}

	log.Printf("restoring file to", path)

	_, err = powershell(shell, fmt.Sprintf(`
        $path = "%s"
        $temp = "%s"

        if (Test-Path $path) {
            rm $path
        }

        $dir = [System.IO.Path]::GetDirectoryName($path)
        if (-Not (Test-Path $dir)) {
            mkdir $dir
        }

        $b64 = Get-Content $temp
        $bytes = [System.Convert]::FromBase64String($b64)

        $file = [System.IO.Path]::GetFullPath($path)
        [System.IO.File]::WriteAllBytes($file, $bytes)
 
        del $temp
    `, path, temp))
	return
}
