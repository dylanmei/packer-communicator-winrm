
## Packer WinRM Plugin

A [Packer](http://www.packer.io/) communicator plugin for interacting with machines using Windows Remote Management.

For more information on WinRM, please visit [Microsoft's WinRM site](http://msdn.microsoft.com/en-us/library/aa384426\(v=VS.85\).aspx).


### Status

This is a work in progress. *It is not a usable Packer plugin yet*. However, while the kinks are being worked out it is also a stand-alone command-line application.

[![wercker status](https://app.wercker.com/status/c702a1133a8359cc8830ad60487ee751/m "wercker status")](https://app.wercker.com/project/bykey/c702a1133a8359cc8830ad60487ee751)

### Usage

A Packer *communicator* plugin must support the following functionality: Execute a shell command, upload a file, download a file, and upload a directory.

#### Executing a shell command

    ./packer-communicator-winrm cmd -user vagrant -pass vagrant "echo hello packer"

#### Uploading a file

*not started*

#### Downloading a file

*not started*

#### Uploading a directory

*not started*

### Props

- joefitzgerald/packer-windows ([https://github.com/joefitzgerald/packer-windows](https://github.com/joefitzgerald/packer-windows))
- mitchellh/packer ([https://github.com/mitchellh/packer](https://github.com/mitchellh/packer))
- winrb/vagrant-windows ([https://github.com/WinRb/vagrant-windows](https://github.com/WinRb/vagrant-windows))
- winrb/winrm ([https://github.com/WinRb/WinRM](https://github.com/WinRb/WinRM))
 