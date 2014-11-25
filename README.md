## Packer WinRM Plugin

A [Packer](http://www.packer.io/) communicator plugin for interacting with machines using Windows Remote Management. For more information on WinRM, visit [Microsoft's WinRM site](http://msdn.microsoft.com/en-us/library/aa384426\(v=VS.85\).aspx).

### Status

This is a work in progress. *It is not a usable Packer plugin yet*. However, while the kinks are being worked out it is also a stand-alone command-line application.

[![Build Status](https://travis-ci.org/mefellows/packer-communicator-winrm.svg?branch=master)](https://travis-ci.org/mefellows/packer-communicator-winrm)
[![Coverage Status](https://coveralls.io/repos/mefellows/packer-communicator-winrm/badge.png?branch=master)](https://coveralls.io/r/mefellows/packer-communicator-winrm?branch=master)

### Usage

A Packer *communicator* plugin supports the following functionality: Execute a shell command, upload a file, download a file, and upload a directory.

#### Help

    alias pcw=`pwd`/packer-communicator-winrm
    pcw help

#### Executing a shell command

    pcw cmd "powershell Write-Host 'Hello' (Get-WmiObject -class Win32_OperatingSystem).Caption"

#### Uploading a file

    pcw file -from=./README.md -to=C:\\Windows\\Temp\\README.md
    pcw cmd "type C:\\Windows\\Temp\\README.md"

#### Uploading a directory

	pcw dir -from="./scripts/" -to="c:\\windows\\temp\\scripts"
	pcw cmd "type c:\\Windows\\Temp\\scripts\\test.sh"

#### Downloading a file

*not started*

### Props

- joefitzgerald/packer-windows ([https://github.com/joefitzgerald/packer-windows](https://github.com/joefitzgerald/packer-windows))
- masterzen/winrm ([https://github.com/masterzen/winrm](https://github.com/masterzen/winrm))
- mitchellh/packer ([https://github.com/mitchellh/packer](https://github.com/mitchellh/packer))
- WinRb/vagrant-windows ([https://github.com/WinRb/vagrant-windows](https://github.com/WinRb/vagrant-windows))
- WinRb/WinRM ([https://github.com/WinRb/WinRM](https://github.com/WinRb/WinRM))
 
