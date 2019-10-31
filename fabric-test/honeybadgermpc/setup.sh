#!/bin/bash

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
if [ $? -ne 0 ]
then
  cd honeybadgerscc
  git pull 
  cd ..
fi
git clone https://github.com/njeans/fabric-test.git
if [ $? -ne 0 ]
then
  cd fabric-test
  git pull 
  cd ..
fi
mkdir -p $GOPATH/src/github.com/hyperledger/
cd $GOPATH/src/github.com/hyperledger/
git clone https://github.com/njeans/fabric.git
if [ $? -ne 0 ]
then
  cd fabric
  git pull 
  cd ..
fi

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

cd $GOPATH/src/github.com/njeans/fabric-test/fabric-test

#copy new docker-compose file that uses named volumes instead of bind mounts
cp /opt/docker-compose-from-docker.yaml $GOPATH/src/github.com/njeans/fabric-test/fabric-test/docker-compose-cli.yaml

./byfn.sh generate

#copy releavant crypto data and binaries to named volumes 
cp $GOPATH/src/github.com/njeans/fabric-test/fabric-test/channel-artifacts/genesis.block /opt/crypto/orderer.example.com/orderer.genesis.block
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/* /opt/crypto/orderer.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/* /opt/crypto/peer0.org1.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/* /opt/crypto/peer1.org1.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/* /opt/crypto/peer0.org2.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/* /opt/crypto/peer1.org2.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/chaincode/rockpaperscissors/ /opt/chaincode/rockpaperscissors/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/* /opt/crypto/cli/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config /opt/crypto/cli/crypto
cp $GOPATH/src/github.com/hyperledger/fabric/.build/bin/peer /opt/peer-bin

#install and instantiate rockpaperscissors chaincode
./byfn.sh up
