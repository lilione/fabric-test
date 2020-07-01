#!/usr/bin/env bash
set -e

# build system chaincode
echo "building system chaincode"
cd $GOPATH/src/github.com/lilione/fabric-test/chaincode/myscc
go build -buildmode=plugin

# build fabric peer binary
#echo "building fabric-peer binary"
#cd $GOPATH/src/github.com/hyperledger/fabric
#GO_TAGS+=" pluginsenabled" make peer

# build fabric peer docker image
#echo "building fabric-peer docker image"
#DOCKER_DYNAMIC_LINK=true GO_TAGS+=" pluginsenabled" make peer-docker IN_DOCKER=true

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh generate

# copy releavant crypto data and binaries to named volumes
cp $GOPATH/src/github.com/lilione/fabric-test/fabric-test/channel-artifacts/genesis.block /opt/crypto/orderer.example.com/orderer.genesis.block
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/* /opt/crypto/orderer.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/* /opt/crypto/peer0.org1.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/* /opt/crypto/peer1.org1.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/* /opt/crypto/peer0.org2.example.com/
cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/* /opt/crypto/peer1.org2.example.com/

# install and instantiate application chaincode
cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka

# start mpc servers
echo "starting mpc servers"
cd $GOPATH/src/github.com/lilione/HoneyBadgerMPC
bash apps/fabric/scripts/start_server.sh

# start client
echo "staring client"
docker run -d \
 -v /Users/lilione/go/src/github.com/lilione/HoneyBadgerMPC:/usr/src/HoneyBadgerMPC \
 -v /Users/lilione/go/src/github.com/lilione/fabric-test/fabric-test/log/chaincode:/usr/src/HoneyBadgerMPC/apps/fabric/log/chaincode \
 -v /var/run/docker.sock:/var/run/docker.sock \
 --name client -it hyperledger/fabric-peer:latest
docker network connect net_byfn client
docker exec -it client bash

# docker exec -it cli bash
# export CHANNEL_NAME=mychannel
# bash scripts/run_cmd.sh update 0 1 "1" "2"
# bash scripts/run_cmd.sh query 0 1 "1"
# bash scripts/run_cmd.sh getInputmaskIdx 0 1
# bash scripts/run_cmd.sh sendMaskedInput 0 1 0 12364027589777478026535591640119302239196130422776730929764185853585614015560
#bash scripts/run_cmd.sh reconstruct 0 1 0

# docker exec -it peer0.org1.example.com bash
# docker exec -it peer1.org1.example.com bash
# docker exec -it peer0.org2.example.com bash
# docker exec -it peer1.org2.example.com bash

# python3.7 apps/fabric/src/client/start_client.py

#docker run \
# -v /Users/lilione/go/src/github.com/lilione/HoneyBadgerMPC:/usr/src/HoneyBadgerMPC \
# -v /Users/lilione/go/src/github.com/lilione/fabric-test/fabric-test/log/chaincode:/usr/src/HoneyBadgerMPC/apps/fabric/log/chaincode \
# -v /var/run/docker.sock:/var/run/docker.sock \
# -it hyperledger/fabric-peer:latest