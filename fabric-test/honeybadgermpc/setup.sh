#!/bin/bash
set -e

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
# cp /opt/docker-compose-from-docker.yaml $GOPATH/src/github.com/njeans/fabric-test/fabric-test/docker-compose-from-docker.yaml

./byfn.sh generate

#copy releavant crypto data and binaries to named volumes
cp $GOPATH/src/github.com/njeans/fabric-test/fabric-test/channel-artifacts/genesis.block /opt/crypto/orderer.example.com/orderer.genesis.block
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/* /opt/crypto/orderer.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/* /opt/crypto/peer0.org1.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/* /opt/crypto/peer1.org1.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/* /opt/crypto/peer0.org2.example.com/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/* /opt/crypto/peer1.org2.example.com/
mkdir -p /opt/chaincode/rockpaperscissors/
cp -r $GOPATH/src/github.com/njeans/fabric-test/chaincode/rockpaperscissors/* /opt/chaincode/rockpaperscissors/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/* /opt/crypto/cli/
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/crypto-config /opt/crypto/cli/crypto
cp -r $GOPATH/src/github.com/njeans/fabric-test/fabric-test/scripts/ /opt/crypto/cli/crypto/scripts
cp $GOPATH/src/github.com/hyperledger/fabric/.build/bin/peer /opt/peer-bin

#install and instantiate rockpaperscissors chaincode
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka
