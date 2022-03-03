# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
    config.vm.box = 'digital_ocean'
    config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
    config.ssh.private_key_path = '~/.ssh/digitalocean'
    config.nfs.functional = false
    config.vm.allowed_synced_folder_types = :rsync
    config.vm.synced_folder "./Vagrant", "/Vagrant", disabled: true

    config.vm.define "go-minitwit", primary: true do |server|

      server.vm.provider :digital_ocean do |provider|
        provider.ssh_key_name = ENV["SSH_KEY_NAME"] # Get it here https://cloud.digitalocean.com/account/security
        provider.token = ENV["DIGITAL_OCEAN_TOKEN"] # Get it here https://cloud.digitalocean.com/account/api
        provider.image = 'ubuntu-20-04-x64'         # Ubuntu 20.04 | vagrant digitalocean-list images $DIGITAL_OCEAN_TOKEN
        provider.region = 'ams3'                    # Amsterdam 3 | vagrant digitalocean-list regions $DIGITAL_OCEAN_TOKEN
        provider.size = 's-1vcpu-1gb'               # 1 vCPU | vagrant digitalocean-list sizes $DIGITAL_OCEAN_TOKEN
        provider.privatenetworking = true
      end

      server.vm.hostname = "go-minitwit"
      server.vm.provision "shell", inline: "echo 'export DUCKDNS_TOKEN=" + ENV["DUCKDNS_TOKEN"] + "' >> ~/.profile" 
      server.vm.provision "shell", inline: "echo 'export DOCKER_PASSWORD=" + ENV["DOCKER_PASSWORD"] + "' >> ~/.profile" 
      server.vm.provision "shell", inline: "echo 'export DOCKER_USERNAME=" + ENV["DOCKER_USERNAME"] + "' >> ~/.profile" 
      server.vm.provision "shell", inline: "echo 'export BUGSNAG_API_KEY=" + ENV["BUGSNAG_API_KEY"] + "' >> ~/.profile"

      server.vm.provision "shell", privileged: true, inline: <<-SHELL
        sudo apt-get update
        source ~/.profile

        ### Install Go-Minitwit as a linux service (Not recommended)
        ### Install Go
        # sudo apt install -y build-essential
        # cd ~
        # curl -OL https://golang.org/dl/go1.16.7.linux-amd64.tar.gz
        # sudo tar -C /usr/local -xvf go1.16.7.linux-amd64.tar.gz
        # sudo rm go1.16.7.linux-amd64.tar.gz

        # echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
        

        # git clone https://github.com/ITU-DevOps-N/go-minitwit.git
        # cd go-minitwit

        # go mod download
        # go mod verify
        # go build -o go-minitwit main.go
        # go build -o go-minitwit-api api/api.go

        # sudo mv go-minitwit go-minitwit-api /usr/local/bin/
        # sudo mv /vagrant/Vagrant/go-minitwit.service /lib/systemd/system/go-minitwit.service
        # sudo mv /vagrant/Vagrant/go-minitwit-api.service /lib/systemd/system/go-minitwit-api.service
        
        # systemctl daemon-reload
        # systemctl enable go-minitwit.service
        # systemctl enable go-minitwit-api.service
        # systemctl start go-minitwit.service
        # systemctl start go-minitwit-api.service

        ### Install Go-Minitwit as Docker Container (Recommended)
        ### Install Docker
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
        apt-cache policy docker-ce
        sudo apt install -y docker-ce
        sudo systemctl status docker
        sudo usermod -aG docker ${USER}
        sudo apt install -y docker-compose

        git clone https://github.com/ITU-DevOps-N/go-minitwit.git
        cd go-minitwit
        echo $DOCKER_PASSWORD | docker login -u $DOCKER_USERNAME --password-stdin
        docker-compose up -d

        ### Installing DuckDNS
        sudo systemctl enable cron
        cd ~
        mkdir duckdns
        mv /vagrant/Vagrant/duck* ~/duckdns/
        chmod 700 ~/duckdns/duck.sh
        bash ~/duckdns/duck.sh
        sed -i "s|DUCKDNS_TOKEN|$DUCKDNS_TOKEN|g" ~/duckdns/duck.sh
        echo '*/5 * * * * ~/duckdns/duck.sh >/dev/null 2>&1' | crontab -
      SHELL
    end
  end