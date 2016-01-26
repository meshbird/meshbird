#!/usr/bin/env bash

apt-get update && apt-get install -y git
curl -O https://storage.googleapis.com/golang/go1.5.3.linux-amd64.tar.gz
tar -xvf go1.5.3.linux-amd64.tar.gz
mv go /usr/local
rm -f go1.5.3.linux-amd64.tar.gz
echo "" >> /home/vagrant/.profile
echo "export PATH=\$PATH:/usr/local/go/bin:/home/vagrant/go/bin" >> /home/vagrant/.profile
echo "export GOPATH=/home/vagrant/go" >> /home/vagrant/.profile
mkdir -p /home/vagrant/go/src/github.com/meshbird
mkdir -p /home/vagrant/go/pkg
mkdir -p /home/vagrant/go/bin
ln -s /vagrant /home/vagrant/go/src/github.com/meshbird/meshbird
chown -R vagrant:vagrant /home/vagrant/go
