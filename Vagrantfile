# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
    config.vm.box = 'digital_ocean'
    config.vm.box_url = "https://github.com/devopsgroup-io/vagrant-digitalocean/raw/master/box/digital_ocean.box"
    config.ssh.private_key_path = '~/.ssh/digitalocean'
    config.nfs.functional = false
    config.vm.allowed_synced_folder_types = :rsync
    config.vm.synced_folder "./Vagrant", "/Vagrant", disabled: true

    config.vm.define "minitwit-manager", primary: true do |server|

      server.vm.provider :digital_ocean do |provider|
        provider.ssh_key_name = ENV["SSH_KEY_NAME"] # Get it here https://cloud.digitalocean.com/account/security
        provider.token = ENV["DIGITAL_OCEAN_TOKEN"] # Get it here https://cloud.digitalocean.com/account/api
        provider.image = 'ubuntu-20-04-x64'         # Ubuntu 20.04 | vagrant digitalocean-list images $DIGITAL_OCEAN_TOKEN
        provider.region = 'ams3'                    # Amsterdam 3 | vagrant digitalocean-list regions $DIGITAL_OCEAN_TOKEN
        provider.size = 's-2vcpu-2gb'               # 2 vCPU | vagrant digitalocean-list sizes $DIGITAL_OCEAN_TOKEN
        provider.tags = ['manager','minitwit_cluster']
        provider.privatenetworking = true
      end

      server.vm.hostname = "minitwit-manager"
      server.vm.provision "shell", inline: 'echo "export DUCKDNS_TOKEN=' + "'" + ENV["DUCKDNS_TOKEN"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_PASSWORD=' + "'" + ENV["DOCKER_PASSWORD"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_USERNAME=' + "'" + ENV["DOCKER_USERNAME"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export BUGSNAG_API_KEY=' + "'" + ENV["BUGSNAG_API_KEY"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PASS=' + "'" + ENV["DB_PASS"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_USER=' + "'" + ENV["DB_USER"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_HOST=' + "'" + ENV["DB_HOST"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PORT=' + "'" + ENV["DB_PORT"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_DATABASE=' + "'" + ENV["DB_DATABASE"] + "'" + '" >> ~/.profile'

      server.vm.provision "shell", privileged: true, inline: <<-SHELL
        sudo apt-get update
        source ~/.profile

        ### Install Go-Minitwit as Docker Container (Recommended)
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
        apt-cache policy docker-ce
        sudo apt install -y docker-ce
        sudo systemctl status docker
        sudo usermod -aG docker ${USER}
        sudo apt install -y docker-compose

        git config --global user.email "minitwit@example.com"
        git config --global user.name "Minitwit"
        git clone https://github.com/ITU-DevOps-N/go-minitwit.git

        ### Docker Swarm Cluster
        SWARM_MANAGER_IP=$(hostname -I | awk '{print $1}')
        docker swarm init --advertise-addr $SWARM_MANAGER_IP
        
        ufw allow 22/tcp && \
        ufw allow 2376/tcp && \
        ufw allow 2377/tcp && \
        ufw allow 7946/tcp && \
        ufw allow 7946/udp && \
        ufw allow 4789/udp && \
        ufw reload && \
        ufw --force  enable && \
        systemctl restart docker

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


    config.vm.define "minitwit-worker1", primary: true do |server|

      server.vm.provider :digital_ocean do |provider|
        provider.ssh_key_name = ENV["SSH_KEY_NAME"] # Get it here https://cloud.digitalocean.com/account/security
        provider.token = ENV["DIGITAL_OCEAN_TOKEN"] # Get it here https://cloud.digitalocean.com/account/api
        provider.image = 'ubuntu-20-04-x64'         # Ubuntu 20.04 | vagrant digitalocean-list images $DIGITAL_OCEAN_TOKEN
        provider.region = 'ams3'                    # Amsterdam 3 | vagrant digitalocean-list regions $DIGITAL_OCEAN_TOKEN
        provider.size = 's-2vcpu-2gb'               # 2 vCPU | vagrant digitalocean-list sizes $DIGITAL_OCEAN_TOKEN
        provider.tags = ['worker','minitwit_cluster']
        provider.privatenetworking = true
      end

      server.vm.hostname = "minitwit-worker1"
      server.vm.provision "shell", inline: 'echo "export DUCKDNS_TOKEN=' + "'" + ENV["DUCKDNS_TOKEN"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_PASSWORD=' + "'" + ENV["DOCKER_PASSWORD"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_USERNAME=' + "'" + ENV["DOCKER_USERNAME"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export BUGSNAG_API_KEY=' + "'" + ENV["BUGSNAG_API_KEY"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PASS=' + "'" + ENV["DB_PASS"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_USER=' + "'" + ENV["DB_USER"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_HOST=' + "'" + ENV["DB_HOST"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PORT=' + "'" + ENV["DB_PORT"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_DATABASE=' + "'" + ENV["DB_DATABASE"] + "'" + '" >> ~/.profile'

      server.vm.provision "shell", privileged: true, inline: <<-SHELL
        sudo apt-get update
        source ~/.profile

        ### Install Go-Minitwit as Docker Container (Recommended)
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
        apt-cache policy docker-ce
        sudo apt install -y docker-ce
        sudo systemctl status docker
        sudo usermod -aG docker ${USER}
        sudo apt install -y docker-compose
      SHELL
    end

    config.vm.define "minitwit-worker2", primary: true do |server|

      server.vm.provider :digital_ocean do |provider|
        provider.ssh_key_name = ENV["SSH_KEY_NAME"] # Get it here https://cloud.digitalocean.com/account/security
        provider.token = ENV["DIGITAL_OCEAN_TOKEN"] # Get it here https://cloud.digitalocean.com/account/api
        provider.image = 'ubuntu-20-04-x64'         # Ubuntu 20.04 | vagrant digitalocean-list images $DIGITAL_OCEAN_TOKEN
        provider.region = 'ams3'                    # Amsterdam 3 | vagrant digitalocean-list regions $DIGITAL_OCEAN_TOKEN
        provider.size = 's-2vcpu-2gb'               # 2 vCPU | vagrant digitalocean-list sizes $DIGITAL_OCEAN_TOKEN
        provider.tags = ['worker','minitwit_cluster']
        provider.privatenetworking = true
      end

      server.vm.hostname = "minitwit-worker2"
      server.vm.provision "shell", inline: 'echo "export DUCKDNS_TOKEN=' + "'" + ENV["DUCKDNS_TOKEN"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_PASSWORD=' + "'" + ENV["DOCKER_PASSWORD"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DOCKER_USERNAME=' + "'" + ENV["DOCKER_USERNAME"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export BUGSNAG_API_KEY=' + "'" + ENV["BUGSNAG_API_KEY"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PASS=' + "'" + ENV["DB_PASS"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_USER=' + "'" + ENV["DB_USER"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_HOST=' + "'" + ENV["DB_HOST"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_PORT=' + "'" + ENV["DB_PORT"] + "'" + '" >> ~/.profile'
      server.vm.provision "shell", inline: 'echo "export DB_DATABASE=' + "'" + ENV["DB_DATABASE"] + "'" + '" >> ~/.profile'

      server.vm.provision "shell", privileged: true, inline: <<-SHELL
        sudo apt-get update
        source ~/.profile

        ### Install Go-Minitwit as Docker Container (Recommended)
        sudo apt install -y apt-transport-https ca-certificates curl software-properties-common
        curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
        sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable"
        apt-cache policy docker-ce
        sudo apt install -y docker-ce
        sudo systemctl status docker
        sudo usermod -aG docker ${USER}
        sudo apt install -y docker-compose
      SHELL
    end
  end