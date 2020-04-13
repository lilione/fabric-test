#!/usr/bin/env bash
set -e

#build system chaincode
#cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test/VIN/scc
#go build -buildmode=plugin -tags ""

#build fabric peer binary
#export GO_TAGS="$GO_TAGS pluginsenabled"
#cd $GOPATH/src/github.com/hyperledger/fabric
#make peer

#build fabric peer docker image
#export DOCKER_DYNAMIC_LINK=true
#make peer-docker IN_DOCKER=true

cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh generate

#cleanup running docker containers
#cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
#./byfn.sh down -f docker-compose-from-docker.yaml -o kafka

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
cp $GOPATH/src/github.com/hyperledger/fabric/.build/bin/peer /opt/peer-bin
#mkdir -p /opt/chaincode/cc/
#cp -r $GOPATH/src/github.com/lilione/fabric-test/fabric-test/VIN/cc/* /opt/chaincode/cc/
mkdir -p /opt/chaincode/cc/
cp -r $GOPATH/src/github.com/lilione/fabric-test/chaincode/cc/* /opt/chaincode/cc/

#install and instantiate application chaincode
cd $GOPATH/src/github.com/lilione/fabric-test/fabric-test
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka

#docker exec -it cli bash
#export CHANNEL_NAME=mychannel
#bash scripts/run_cmd.sh update 0 1 "1" "2"
#bash scripts/run_cmd.sh query 0 1 "1"
#bash scripts/run_cmd.sh createGame 0 1 "game2" 100 "user1"
#bash scripts/run_cmd.sh getActiveGames 0 1

#peer chaincode install -n rpscc -v 1.0 -l golang -p github.com/chaincode/cc
#peer chaincode instantiate -o orderer.example.com:7050 --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem -C mychannel -n rpscc -l golang -v 1.0 -c '{"Args":["init"]}' -P 'OR ('\''Org1MSP.peer'\'','\''Org2MSP.peer'\'')'




