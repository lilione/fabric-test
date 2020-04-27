cd /opt/gopath/src/github.com/hyperledger/fabric/peer
peer node start &

cd /usr/src/HoneyBadgerMPC
python3.7 apps/fabric/test/test_server.py &> apps/fabric/test/log.txt

#cd /usr/src/HoneyBadgerMPC/apps/fabric/test

#cd /usr/src/HoneyBadgerMPC
#python3.7 apps/fabric/test/test_client.py
