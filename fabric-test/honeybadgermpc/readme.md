# Setup Fabric network
* Build peer binary
* Build peer docker images
* Run install and instantiate rockpaperscissors chaincode
* start network with 4 peers from 2 orgs, 1 orderer, 1 kafka, 1 zookeeper node
```
docker volume create orderer.example.com-crypto
docker volume create peer0.org1.example.com-msp
docker volume create peer0.org2.example.com-msp
docker volume create peer1.org1.example.com-msp
docker volume create peer1.org2.example.com-msp
docker volume create peer0.org1.example.com-tls
docker volume create peer0.org2.example.com-tls
docker volume create peer1.org1.example.com-tls
docker volume create peer1.org2.example.com-tls
docker volume create honeybadgerscc
docker volume create chaincode
docker volume create cli-config
docker volume create peer-bin
docker build -t setup-honeybadger .
docker run -it \
	-v /var/run/docker.sock:/var/run/docker.sock \
	-v ${PWD}/setup.sh:/opt/setup.sh \
	-v peer-bin:/opt/peer-bin \
	-v orderer.example.com-crypto:/opt/crypto/orderer.example.com \
	-v peer0.org1.example.com-tls:/opt/crypto/peer0.org1.example.com/tls \
	-v peer0.org1.example.com-msp:/opt/crypto/peer0.org1.example.com/msp \
	-v peer1.org1.example.com-tls:/opt/crypto/peer1.org1.example.com/tls \
	-v peer1.org1.example.com-msp:/opt/crypto/peer1.org1.example.com/msp \
	-v peer0.org2.example.com-tls:/opt/crypto/peer0.org2.example.com/tls \
	-v peer0.org2.example.com-msp:/opt/crypto/peer0.org2.example.com/msp \
	-v peer1.org2.example.com-tls:/opt/crypto/peer1.org2.example.com/tls \
	-v peer1.org2.example.com-msp:/opt/crypto/peer1.org2.example.com/msp \
	-v honeybadgerscc:/opt/gopath/src/github.com/njeans/honeybadgerscc \
	-v chaincode:/opt/chaincode \
	-v cli-config:/opt/crypto/cli \
	--name setup-honeybadger \
	-w /opt setup-honeybadger \
	bash setup.sh
```

# Send commands to rockpaperscissors chaincode
First must ssh into cli container to send commands to the peer nodes
```
docker exec -it cli bash
export CHANNEL_NAME=mychannel
```

* `createGame <peer id> <peer org id> <game name> <timeout> <user name>`

Example:

```
./scripts/run_cmd.sh createGame 0 1 "game1" 100 "user1"
```

* `chaincodeQuery <peer id> <peer org id> <arg 1> <arg 2> ... <arg n>`

Example:
```
./scripts/run_cmd.sh chaincodeQuery 0 1 "game1"
```

* `getActiveGames <peer id> <peer org id>`
* `getCompletedGames <peer id> <peer org id>`
* `joinGame <peer id> <peer org id> <game name> <user name>`
* `openMoves <peer id> <peer org id> <game name>`
* `endGame <peer id> <peer org id> <game name>`

#Restart network

must be done from fabric-test/fabric-test directory

```
cd ..
mkdir crypto-config
./byfn.sh restart -f docker-compose-from-docker.yaml -o kafka
```

#Cleanup
must be done from fabric-test/fabric-test directory

```
cd ..
./byfn.sh down -f docker-compose-from-docker.yaml -o kafka
docker container rm setup-honeybadger
docker volume rm orderer.example.com-crypto peer0.org1.example.com-msp peer0.org1.example.com-tls peer0.org2.example.com-tls peer0.org2.example.com-msp peer1.org1.example.com-tls peer1.org1.example.com-msp peer1.org2.example.com-tls peer1.org2.example.com-msp honeybadgerscc chaincode cli-config peer-bin
docker rmi setup-honeybadger
```
