#!/usr/bin/env bash

docker exec peer0.org1.example.com pkill -f hand_off_item
docker exec peer0.org2.example.com pkill -f hand_off_item
docker exec peer1.org1.example.com pkill -f hand_off_item
docker exec peer1.org2.example.com pkill -f hand_off_item

docker exec peer0.org1.example.com pkill -f start_server
docker exec peer0.org2.example.com pkill -f start_server
docker exec peer1.org1.example.com pkill -f start_server
docker exec peer1.org2.example.com pkill -f start_server

# start mpc servers
echo "starting mpc servers"
cd $GOPATH/src/github.com/lilione/HoneyBadgerMPC
bash apps/fabric/scripts/start_server.sh

docker exec -it client bash