# Go Minitwit!
[![DigitalOcean Referral Badge](https://web-platforms.sfo2.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=7cb197c4e0cb&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

## Go-Minitwit on Digital Ocean
To deploy our application on Digital Ocean, we used Vagrant. Vagrant is a tool that allows us to create and provision virtual machines.
In our case, Vagrant will create a droplet on Digital Ocean and provision it with the necessary software and configuration.
More specifically it will:
- Update the software on the VM
- Install Go, GCC, and other dependencies necessary for Go-Minitwit
- Clone Go-Minitwit repository
- Build Go-Minitwit application to binary
- Create a `systemd` service
- Start Go-Minitwit service
- 

### Deploy
Install Vagrant plugin for Digital Ocean:
`vagrant plugin install vagrant-digitalocean`

Create your access token [here](https://cloud.digitalocean.com/account/api/tokens)

```shell
export DIGITAL_OCEAN_TOKEN=''
export SSH_KEY_NAME=''
```

Make sure that your private key for `$SSH_KEY_NAME` is in your `~/.ssh/` directory.

Run vagrant up:
`vagrant up`

[Access to the application!](go-minitwit.duckdns.org)
