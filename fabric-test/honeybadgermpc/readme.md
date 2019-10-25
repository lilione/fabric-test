Build peer and peer docker images and run install and instantiate rockpaperscissors chaincode

```docker run -it -v /var/run/docker.sock:/var/run/docker.sock -v ${PWD}/setup.sh:/opt/setup.sh -w /opt hyperledger/fabric-baseimage:amd64-0.4.14 sh setup.sh```
