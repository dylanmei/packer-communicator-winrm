require 'bundler'; Bundler.require
require 'winrm'

endpoint = 'http://localhost:5985/wsman'

winrm = WinRM::WinRMWebService.new(endpoint, :plaintext,
  :user => 'vagrant', :pass => 'vagrant', :basic_auth_only => true)

# winrm.cmd('ipconfig /all') do |stdout, stderr|
#   p "= pconfig /all\n"
#   STDOUT.print stdout
#   STDERR.print stderr
# end
puts
winrm.cmd('set') do |stdout, stderr|
#  p "= set\n"
  STDOUT.print stdout
  STDERR.print stderr
end