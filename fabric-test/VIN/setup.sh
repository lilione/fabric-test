#!/usr/bin/env bash
set -e

#build system chaincode
cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/myscc
go build -buildmode=plugin

#build fabric peer binary
cd $GOPATH/src/github.com/hyperledger/fabric
GO_TAGS+=" pluginsenabled" make peer

#build fabric peer docker image
DOCKER_DYNAMIC_LINK=true GO_TAGS+=" pluginsenabled" make peer-docker IN_DOCKER=true

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh generate

#copy releavant crypto data and binaries to named volumes
cp $GOPATH/src/github.com/lilione/fabric-test/fabric-test/channel-artifacts/genesis.block /opt/crypto/orderer.example.com/orderer.genesis.block
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/* /opt/crypto/orderer.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/* /opt/crypto/peer0.org1.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/* /opt/crypto/peer1.org1.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/* /opt/crypto/peer0.org2.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/* /opt/crypto/peer1.org2.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/* /opt/crypto/cli/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config /opt/crypto/cli/crypto
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/scripts/ /opt/crypto/cli/crypto/scripts
mkdir -p /opt/chaincode/cc/
cp -r $GOPATH/src/github.com/lilione/fabric-test/chaincode/cc/* /opt/chaincode/cc/

#install and instantiate application chaincode
cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka


#docker exec -it cli bash
#export CHANNEL_NAME=mychannel
#bash scripts/run_cmd.sh update 0 1 "1" "2"
#bash scripts/run_cmd.sh query 0 1 "1"

#docker exec -it peer0.org1.example.com bash

