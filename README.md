# 1. Go Minitwit!
[![DigitalOcean Referral Badge](https://web-platforms.sfo2.digitaloceanspaces.com/WWW/Badge%203.svg)](https://www.digitalocean.com/?refcode=7cb197c4e0cb&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

## 1.1. Go-Minitwit on Digital Ocean
We used Vagrant to deploy our application to [DigitalOcean](https://www.digitalocean.com/). Vagrant is a tool that allows us to create and provision virtual machines.
In our case, Vagrant will create a droplet on DigitalOcean and provision it with the necessary software and configuration.
Specifically, it will:
- Update the software on the VM
- Install Go, GCC, and other dependencies necessary for Go-Minitwit
- Clone Go-Minitwit repository
- Build Go-Minitwit application to binary
- Create a `systemd` service
- Start Go-Minitwit service

### 1.1.1. Deploy

In order to deploy the application to DigitalOcean you will need to have Vagrant and Virtualbox installed on your machine. Afterwards, you will need the Vagrant plugin for DigitalOcean.

Install the DigitalOCean plugin: `vagrant plugin install vagrant-digitalocean`

Create your access token [here](https://cloud.digitalocean.com/account/api/tokens)

```shell
export DIGITAL_OCEAN_TOKEN=''
export SSH_KEY_NAME=''
export DUCKDNS_TOKEN=''
```

Make sure that your private key for `$SSH_KEY_NAME` is in your `~/.ssh/` directory.

Run vagrant up:
`vagrant up`

[Access to the application!](http://go-minitwit.duckdns.org)

## 1.2. Development

You will need to setup environment variables to pull the image from [Docker Hub](https://hub.docker.com/). This is needed to run the application. The environment variables has been shared in the organization; if you do not have these environment variables, ask Emil or Gianmarco for them.

### 1.2.1. Environment variables

Once you have the environment variables, set it up like this:

1. Create a file on your system that contains the environment variables. This SHOULD not be at the project level or anywhere in version control. Example: `~/secrets/go-minitwit`
2. To expose the variables to the running terminal session, do `source ~/secrets/go-minitwit`
3. Confirm that the variables are exposed, example `echo $DOCKER_USERNAME`

### 1.2.3. Docker compose (running the application locally)

You will need to login to Docker Hub with the organization credentials. This can be once you have exposed the variables in your terminal session, 

Login to organization on Docker Hub: `docker login --username itudevops -p $DOCKER_PASSWORD`

Start the application locally with `docker-compose up`.

## 1.3 Bugsnag

Login credentials to Bugsnag is available through Emil or Gianmarco.


