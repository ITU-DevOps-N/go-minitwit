#!/bin/bash

PATH_KEY="~/.ssh/digitalocean"

echo "Insert the last tag for the release (e.g. '32'):"
read TAG
echo "Insert DIGITAL_OCEAN_TOKEN:"
read DIGITAL_OCEAN_TOKEN
source ~/University/digital_ocean

echo "Vagrant up"
vagrant up


export SWARM_MANAGER_IP=`curl -s GET \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DIGITAL_OCEAN_TOKEN"\
    "https://api.digitalocean.com/v2/droplets" \
    | jq -r '.droplets|.[]|select(.name == "minitwit-manager-1")|.networks.v4|.[]|select(.type == "public")|.ip_address'`

echo "SWARM_MANAGER_IP=$SWARM_MANAGER_IP"


export WORKER1_IP=`curl -s GET \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $DIGITAL_OCEAN_TOKEN"\
    "https://api.digitalocean.com/v2/droplets" \
    | jq -r '.droplets|.[]|select(.name == "minitwit-worker1-1")|.networks.v4|.[]|select(.type == "public")|.ip_address'`

echo "WORKER1_IP=$WORKER1_IP"

echo "Getting join token from minitwit-manager"
export WORKER_TOKEN=`ssh -i $PATH_KEY -o "StrictHostKeyChecking no" root@$SWARM_MANAGER_IP "docker swarm join-token worker -q"`
echo "WORKER_TOKEN=$WORKER_TOKEN"
echo "Joining minitwit-worker1 to minitwit-manager swarm cluster"
ssh -i $PATH_KEY -o "StrictHostKeyChecking no" root@$WORKER1_IP "docker swarm join --token $WORKER_TOKEN $SWARM_MANAGER_IP:2377"

echo "Checking docker swarm cluster status"
ssh -i $PATH_KEY -o "StrictHostKeyChecking no" root@$SWARM_MANAGER_IP "docker node ls"


ssh -i $PATH_KEY -o "StrictHostKeyChecking no" root@$SWARM_MANAGER_IP "source ~/.profile && \
    cd go-minitwit && \
    sed -ir "s/TAG/$TAG/" docker-compose.prod.yml && \
    docker pull itudevops/go-minitwit:$TAG && \
    docker pull itudevops/go-minitwit-api:$TAG && \
    docker stack deploy --compose-file docker-compose.prod.yml minitwit"

echo "Waiting for docker swarm serivice to start"
sleep 12
ssh -i $PATH_KEY -o "StrictHostKeyChecking no" root@$SWARM_MANAGER_IP "docker service ls"

echo "Done!"