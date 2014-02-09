package main

import (
	"fmt"
	"github.com/dylanmei/packer-communicator-winrm/winrm"
	"github.com/mitchellh/packer/common/uuid"
	"io"
	"log"
)

func upload(shell *winrm.Shell, path string, r io.Reader) (err error) {
	_, err = powershell(shell, fmt.Sprintf(`
        if (Test-Path "%s") {
            rm "%s"
        }`, path, path))

	if err != nil {
		return
	}

	temp, err := cmd(shell,
		fmt.Sprintf("echo %%TEMP%%\\packer-%s.tmp", uuid.TimeOrderedUUID()))

	if err != nil {
		return
	}

	// Base64.encode64(IO.binread(from)).gsub("\n",'').chars.to_a.each_slice(8000-file_name.size) do |chunk|
	//   out = cmd("echo #{chunk.join} >> \"#{file_name}\"")
	// end

	log.Printf("writing to temporary file [%s]", temp)
	_, err = cmd(shell, fmt.Sprintf("echo aABlAGwAbABvACAAdwBvAHIAbABkAA== >> \"%s\"", temp))

	log.Printf("restoring file to [%s]", path)

	_, err = powershell(shell, fmt.Sprintf(`
        $path = "%s"
        $dir = [System.IO.Path]::GetDirectoryName($path)
        if (-Not (Test-Path $dir)) { mkdir $dir }
 
        $temp = "%s"
        $b64 = Get-Content $temp
        $bytes = [System.Convert]::FromBase64String($b64)

        $file = [System.IO.Path]::GetFullPath($path)
        [System.IO.File]::WriteAllBytes($file, $bytes)
 
        del $temp
    `, path, temp))
	return
}
