# -*- mode: ruby -*-
# vi: set ft=ruby :

# Vagrantfile API/syntax version. Don't touch unless you know what you're doing!
VAGRANTFILE_API_VERSION = "2"

# Original found at
# https://github.com/nathany/vagrant-gopher/blob/master/Vagrantfile


Vagrant.require_version ">= 1.5.0"

# See https://code.google.com/p/go/downloads/list
GO_ARCHIVES = {
  "linux" => "go1.3.linux-amd64.tar.gz"
}

INSTALL = {
  "linux" => "apt-get update -qq; apt-get install -qq -y git mercurial bzr curl"
}

# location of the Vagrantfile
def src_path
  File.dirname(__FILE__)
end

SRC = "/home/vagrant/src/"
PROJECT_SRC = SRC + "github.com/127biscuits/apihippo.com"

# shell script to bootstrap Go
def bootstrap(box)
  install = INSTALL[box]
  archive = GO_ARCHIVES[box]

  profile = <<-PROFILE
    export GOPATH=$HOME
    export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin
    export CDPATH=.:$GOPATH/src/github.com:$GOPATH/src/code.google.com/p:$GOPATH/src/bitbucket.org:$GOPATH/src/launchpad.net
  PROFILE

  <<-SCRIPT
  #{install}

  if ! [ -f /home/vagrant/#{archive} ]; then
    response=$(curl -O# https://storage.googleapis.com/golang/#{archive})
  fi
  tar -C /usr/local -xzf #{archive}

  echo '#{profile}' >> /home/vagrant/.profile

  chown -R vagrant.vagrant #{SRC}

  echo "\nRun: vagrant ssh #{box} -c 'cd project/path; go test ./...'"
  SCRIPT
end

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|


    config.vm.define "linux" do |linux|
        linux.vm.box = "ubuntu/trusty64"
        linux.vm.synced_folder src_path, PROJECT_SRC
        config.vm.hostname = "apihippo.dev"
        config.vm.network "private_network", ip: "192.168.53.10"
        config.hostsupdater.aliases = ["cdn.apihippo.dev", "random.apihippo.dev"]
        linux.vm.provision :shell, :inline => bootstrap("linux")
    end



end


