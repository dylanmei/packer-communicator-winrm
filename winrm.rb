require 'bundler'; Bundler.require
require 'winrm'

endpoint = 'http://localhost:5985/wsman'

winrm = WinRM::WinRMWebService.new(endpoint, :plaintext,
  :user => 'vagrant', :pass => 'D@rj33l1n9', :basic_auth_only => true)

puts 'running SET'
winrm.cmd('set') do |stdout, stderr|
  STDOUT.print stdout
  STDERR.print stderr
end