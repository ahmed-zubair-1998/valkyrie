#!/bin/bash -xe

export DEBIAN_FRONTEND=noninteractive
cd /home/ubuntu

sudo apt update
sudo apt -y install golang-go
git clone https://github.com/ahmed-zubair-1998/valkyrie.git
cd valkyrie/frontend

sudo go run . --dispatcher=http://${DISPATCHER_PUBLIC_DNS}:8090
