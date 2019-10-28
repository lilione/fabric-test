#!/bin/bash
set -e

#dependencies
apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg2 \
    software-properties-common

#install Docker & Docker-compose
curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add -
add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   $(lsb_release -cs) \
   stable"
apt-get update && apt-get install -y \
    docker-ce \
    docker-ce-cli \
    containerd.io
curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

#fabric golang dependencies
go get github.com/syndtr/goleveldb/leveldb

#download github repos
mkdir -p $GOPATH/src/github.com/njeans/
cd $GOPATH/src/github.com/njeans/
git clone https://github.com/njeans/honeybadgerscc.git
git clone https://github.com/njeans/fabric-test.git
mkdir -p $GOPATH/src/github.com/hyperledger/
cd $GOPATH/src/github.com/hyperledger/
git clone https://github.com/njeans/fabric.git

#build system chaincode
cd $GOPATH/src/github.com/njeans/honeybadgerscc/
go build -buildmode=plugin -tags ""

#build fabric peer binary
export DOCKER_DYNAMIC_LINK=true
export GO_TAGS="$GO_TAGS pluginsenabled"
cd $GOPATH/src/github.com/hyperledger/fabric
make peer
#build fabric peer docker image
make peer-docker IN_DOCKER=true

#install and instantiate rockpaperscissors chaincode
cd $GOPATH/src/github.com/njeans/fabric-test/fabric-test
./byfn.sh generate
./byfn.sh up
