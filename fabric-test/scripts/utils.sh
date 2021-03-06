#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#

# This is a collection of bash functions used by different scripts

ORDERER_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
PEER0_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
PEER1_ORG1_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org1.example.com/peers/peer1.org1.example.com/tls/ca.crt
PEER0_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt
PEER1_ORG2_CA=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org2.example.com/peers/peer1.org2.example.com/tls/ca.crt

# verify the result of the end-to-end test
verifyResult() {
  if [ $1 -ne 0 ]; then
    echo "!!!!!!!!!!!!!!! "$2" !!!!!!!!!!!!!!!!"
    echo "========= ERROR !!! FAILED to execute End-2-End Scenario ==========="
    echo
    exit 1
  fi
}

# Set OrdererOrg.Admin globals
setOrdererGlobals() {
  CORE_PEER_LOCALMSPID="OrdererMSP"
  CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
  CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/ordererOrganizations/example.com/users/Admin@example.com/msp
}

setGlobals() {
  PEER=$1
  ORG=$2
  if [ $ORG -eq 1 ]; then
    CORE_PEER_LOCALMSPID="Org1MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG1_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org1.example.com:7051
    else
      CORE_PEER_ADDRESS=peer1.org1.example.com:7051
    fi
  elif [ $ORG -eq 2 ]; then
    CORE_PEER_LOCALMSPID="Org2MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG2_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org2.example.com:7051
    else
      CORE_PEER_ADDRESS=peer1.org2.example.com:7051
    fi

  elif [ $ORG -eq 3 ]; then
    CORE_PEER_LOCALMSPID="Org3MSP"
    CORE_PEER_TLS_ROOTCERT_FILE=$PEER0_ORG3_CA
    CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto-config/peerOrganizations/org3.example.com/users/Admin@org3.example.com/msp
    if [ $PEER -eq 0 ]; then
      CORE_PEER_ADDRESS=peer0.org3.example.com:7051
    else
      CORE_PEER_ADDRESS=peer1.org3.example.com:7051
    fi
  else
    echo "================== ERROR !!! ORG Unknown =================="
  fi

  if [ "$VERBOSE" == "true" ]; then
    env | grep CORE
  fi
}

updateAnchorPeers() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG

  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx >&log.txt
    res=$?
    set +x
  else
    set -x
    peer channel update -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/${CORE_PEER_LOCALMSPID}anchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA >&log.txt
    res=$?
    set +x
  fi
  cat log.txt
  verifyResult $res "Anchor peer update failed"
  echo "===================== Anchor peers updated for org '$CORE_PEER_LOCALMSPID' on channel '$CHANNEL_NAME' ===================== "
  sleep $DELAY
  echo
}

## Sometimes Join takes time hence RETRY at least 5 times
joinChannelWithRetry() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG

  set -x
  peer channel join -b $CHANNEL_NAME.block >&log.txt
  res=$?
  set +x
  cat log.txt
  if [ $res -ne 0 -a $COUNTER -lt $MAX_RETRY ]; then
    COUNTER=$(expr $COUNTER + 1)
    echo "peer${PEER}.org${ORG} failed to join the channel, Retry after $DELAY seconds"
    sleep $DELAY
    joinChannelWithRetry $PEER $ORG
  else
    COUNTER=1
  fi
  verifyResult $res "After $MAX_RETRY attempts, peer${PEER}.org${ORG} has failed to join channel '$CHANNEL_NAME' "
}

installChaincode() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  VERSION=${3:-1.0}
  set -x
  peer chaincode install -n $CC_NAME -v ${VERSION} -l ${LANGUAGE} -p ${CC_SRC_PATH} >&log.txt
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Chaincode installation on peer${PEER}.org${ORG} has failed"
  echo "===================== Chaincode is installed on peer${PEER}.org${ORG} ===================== "
  echo
}

instantiateChaincode() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG
  VERSION=${3:-1.0}

  # while 'peer chaincode' command can get the orderer endpoint from the peer
  # (if join was successful), let's supply it directly as we know it using
  # the "-o" option
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer chaincode instantiate -o orderer.example.com:7050 -C $CHANNEL_NAME -n $CC_NAME -l ${LANGUAGE} -v ${VERSION} -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" >&log.txt
    res=$?
    set +x
  else
    set -x
    peer chaincode instantiate -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME -l ${LANGUAGE} -v 1.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')" >&log.txt
    res=$?
    set +x
  fi
  cat log.txt
  verifyResult $res "Chaincode instantiation on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' failed"
  echo "===================== Chaincode is instantiated on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' ===================== "
  echo
}

upgradeChaincode() {
  PEER=$1
  ORG=$2
  setGlobals $PEER $ORG

  set -x
  peer chaincode upgrade -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME -v 2.0 -c '{"Args":["init"]}' -P "OR ('Org1MSP.peer','Org2MSP.peer')"
  res=$?
  set +x
  cat log.txt
  verifyResult $res "Chaincode upgrade on peer${PEER}.org${ORG} has failed"
  echo "===================== Chaincode is upgraded on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' ===================== "
  echo
}

chaincodeQuery() {
  ARG=("$@")
  PEER="${ARG[1]}"
  ORG="${ARG[2]}"
  setGlobals $PEER $ORG
  EXPECTED_RESULT=$3

  #formating the query functions inputs
  INPUTS=("${ARG[@]:3}")
  len=${#INPUTS[@]}
  s=""
  if [ "$len" -gt "0" ]; then
    s=",\"${INPUTS[0]}\""
    for i in "${INPUTS[@]:1}"
    do
       s="$s,\"$i\""
    done
  fi
  echo $s

  echo "===================== Querying on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME'... ===================== "
  local rc=1
  local starttime=$(date +%s)

  # continue to poll
  # we either get a successful response, or reach TIMEOUT
  while
    test "$(($(date +%s) - starttime))" -lt "$TIMEOUT" -a $rc -ne 0
  do
    sleep $DELAY
    echo "Attempting to Query peer${PEER}.org${ORG} ...$(($(date +%s) - starttime)) secs"
    set -x
    peer chaincode query -C $CHANNEL_NAME -n $CC_NAME -c '{"Args":["query"'$INPUTS']}' >&log.txt
    res=$?
    set +x
    test $res -eq 0 && VALUE=$(cat log.txt | awk '/Query Result/ {print $NF}')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
    # removed the string "Query Result" from peer chaincode query command
    # result. as a result, have to support both options until the change
    # is merged.
    test $rc -ne 0 && VALUE=$(cat log.txt | egrep '^[0-9]+$')
    test "$VALUE" = "$EXPECTED_RESULT" && let rc=0
  done
  echo
  cat log.txt
  if test $rc -eq 0; then
    echo "===================== Query successful on peer${PEER}.org${ORG} on channel '$CHANNEL_NAME' ===================== "
  else
    echo "!!!!!!!!!!!!!!! Query result on peer${PEER}.org${ORG} is INVALID !!!!!!!!!!!!!!!!"
    echo "================== ERROR !!! FAILED to execute End-2-End Scenario =================="
    echo
    exit 1
  fi
}

# fetchChannelConfig <channel_id> <output_json>
# Writes the current channel config for a given channel to a JSON file
fetchChannelConfig() {
  CHANNEL=$1
  OUTPUT=$2

  setOrdererGlobals

  echo "Fetching the most recent configuration block for the channel"
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer channel fetch config config_block.pb -o orderer.example.com:7050 -c $CHANNEL --cafile $ORDERER_CA
    set +x
  else
    set -x
    peer channel fetch config config_block.pb -o orderer.example.com:7050 -c $CHANNEL --tls --cafile $ORDERER_CA
    set +x
  fi

  echo "Decoding config block to JSON and isolating config to ${OUTPUT}"
  set -x
  configtxlator proto_decode --input config_block.pb --type common.Block | jq .data.data[0].payload.data.config >"${OUTPUT}"
  set +x
}

# signConfigtxAsPeerOrg <org> <configtx.pb>
# Set the peerOrg admin of an org and signing the config update
signConfigtxAsPeerOrg() {
  PEERORG=$1
  TX=$2
  setGlobals 0 $PEERORG
  set -x
  peer channel signconfigtx -f "${TX}"
  set +x
}

# createConfigUpdate <channel_id> <original_config.json> <modified_config.json> <output.pb>
# Takes an original and modified config, and produces the config update tx
# which transitions between the two
createConfigUpdate() {
  CHANNEL=$1
  ORIGINAL=$2
  MODIFIED=$3
  OUTPUT=$4

  set -x
  configtxlator proto_encode --input "${ORIGINAL}" --type common.Config >original_config.pb
  configtxlator proto_encode --input "${MODIFIED}" --type common.Config >modified_config.pb
  configtxlator compute_update --channel_id "${CHANNEL}" --original original_config.pb --updated modified_config.pb >config_update.pb
  configtxlator proto_decode --input config_update.pb --type common.ConfigUpdate >config_update.json
  echo '{"payload":{"header":{"channel_header":{"channel_id":"'$CHANNEL'", "type":2}},"data":{"config_update":'$(cat config_update.json)'}}}' | jq . >config_update_in_envelope.json
  configtxlator proto_encode --input config_update_in_envelope.json --type common.Envelope >"${OUTPUT}"
  set +x
}

# parsePeerConnectionParameters $@
# Helper function that takes the parameters from a chaincode operation
# (e.g. invoke, query, instantiate) and checks for an even number of
# peers and associated org, then sets $PEER_CONN_PARMS and $PEERS
parsePeerConnectionParameters() {
  # check for uneven number of peer and org parameters
  if [ $(($# % 2)) -ne 0 ]; then
    exit 1
  fi

  PEER_CONN_PARMS=""
  PEERS=""
  while [ "$#" -gt 0 ]; do
    PEER="peer$1.org$2"
    PEERS="$PEERS $PEER"
    PEER_CONN_PARMS="$PEER_CONN_PARMS --peerAddresses $PEER.example.com:7051"
#    echo ${CORE_PEER_TLS_ENABLED}
#    echo ${PWD}
#    echo eval echo "--tlsRootCertFiles \$PEER$1_ORG$2_CA"
    if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "true" ]; then
      TLSINFO=$(eval echo "--tlsRootCertFiles \$PEER$1_ORG$2_CA")
      PEER_CONN_PARMS="$PEER_CONN_PARMS $TLSINFO"
    fi
    # shift by two to get the next pair of peer/org parameters
    shift
    shift
  done
  # remove leading space for output
  PEERS="$(echo -e "$PEERS" | sed -e 's/^[[:space:]]*//')"
}

invoke() {
  echo $1 $2 $3
  parsePeerConnectionParameters $1 $2
  res=$?
  verifyResult $res "Invoke transaction failed on channel '$CHANNEL_NAME' due to uneven number of peer and org parameters "

  # while 'peer chaincode' command can get the orderer endpoint from the
  # peer (if join was successful), let's supply it directly as we know
  # it using the "-o" option
#  echo ${CORE_PEER_TLS_ENABLED}
  if [ -z "$CORE_PEER_TLS_ENABLED" -o "$CORE_PEER_TLS_ENABLED" = "false" ]; then
    set -x
    peer chaincode invoke -o orderer.example.com:7050 -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c $3 >& log/chaincode/$4_peer$1org$2.txt
    res=$?
    set +x
  else
    set -x
    peer chaincode invoke -o orderer.example.com:7050 --tls $CORE_PEER_TLS_ENABLED --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c $3 >& log/chaincode/$4_peer$1org$2.txt
    res=$?
    set +x
  fi
  cat log/chaincode/$4_peer$1org$2.txt
  verifyResult $res "Invoke execution on $PEERS failed "
  echo "===================== Invoke transaction successful on $PEERS on channel '$CHANNEL_NAME' ===================== "
  echo
}

# for supplychain version 1

1_queryShipment() {
  CC_NAME=supplychain_cc_1
  echo "queryShipment"
  arg="{\"Args\":[\"queryShipment\",\"$3\",\"$4\"]}"
  invoke $1 $2 $arg "queryShipment"
}

1_queryInquiry() {
  CC_NAME=supplychain_cc_1
  echo "queryInquiry"
  arg="{\"Args\":[\"queryInquiry\",\"$3\",\"$4\"]}"
  invoke $1 $2 $arg "queryInquiry"
}

1_getInputmaskIdx() {
  CC_NAME=supplychain_cc_1
  echo "getInputmaskIdx"
  arg="{\"Args\":[\"getInputmaskIdx\",\"$3\"]}"
  invoke $1 $2 $arg "getInputmaskIdx"
}

1_registerItemClientGlobal() {
  CC_NAME=supplychain_cc_1
  echo "registerItemClientGlobal"
  arg="{\"Args\":[\"registerItemClientGlobal\",\"$3\"]}"
  invoke $1 $2 $arg "registerItemClientGlobal"
}

1_registerItemClientLocal() {
  CC_NAME=supplychain_cc_1
  echo "registerItemClientLocal"
  arg="{\"Args\":[\"registerItemClientLocal\",\"$3\"]}"
  invoke $1 $2 $arg "registerItemClientLocal"
}

1_handOffItemClientGlobal() {
  CC_NAME=supplychain_cc_1
  echo "handOffItemClientGlobal"
  arg="{\"Args\":[\"handOffItemClientGlobal\",\"$3\"]}"
  invoke $1 $2 $arg "handOffItemClientGlobal"
}

1_handOffItemClientLocal() {
  CC_NAME=supplychain_cc_1
  echo "handOffItemClientLocal"
  arg="{\"Args\":[\"handOffItemClientLocal\",\"$3\"]}"
  invoke $1 $2 $arg "handOffItemClientLocal"
}

1_handOffItemServerGlobal() {
  CC_NAME=supplychain_cc_1
  echo "handOffItemServerGlobal"
  arg="{\"Args\":[\"handOffItemServerGlobal\",\"$3\"]}"
  invoke $1 $2 $arg "handOffItemServerGlobal"
}

1_sourceItemClientLocal() {
  CC_NAME=supplychain_cc_1
  echo "sourceItemClientLocal"
  arg="{\"Args\":[\"sourceItemClientLocal\",\"$3\"]}"
  invoke $1 $2 $arg "sourceItemClientLocal"
}

1_sourceItemServerGlobal() {
  CC_NAME=supplychain_cc_1
  echo "sourceItemServerGlobal"
  arg="{\"Args\":[\"sourceItemServerGlobal\",\"$3\"]}"
  invoke $1 $2 $arg "sourceItemServerGlobal"
}

# for supplychain version 2

2_queryShipment() {
  CC_NAME=supplychain_cc_2
  echo "queryShipment"
  arg="{\"Args\":[\"queryShipment\",\"$3\",\"$4\"]}"
  invoke $1 $2 $arg "queryShipment"
}

2_registerItem() {
  CC_NAME=supplychain_cc_2
  echo "2_registerItem"
  arg="{\"Args\":[\"registerItem\",\"$3\",\"$4\"]}"
  invoke $1 $2 $arg "registerItem"
}

2_handOffItemToNextProvider() {
  CC_NAME=supplychain_cc_2
  echo "2_handOffItemToNextProvider"
  arg="{\"Args\":[\"handOffItemToNextProvider\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\",\"$8\",\"$9\"]}"
  invoke $1 $2 $arg "handOffItemToNextProvider"
}

# for supplychain version 3
3_getInputmaskIdx() {
  CC_NAME=supplychain_cc_3
  echo "getInputmaskIdx"
  arg="{\"Args\":[\"getInputmaskIdx\",\"$3\"]}"
  invoke $1 $2 $arg "getInputmaskIdx"
}

3_createTruckGlobal() {
  CC_NAME=supplychain_cc_3
  echo "createTruckGlobal"
  arg="{\"Args\":[\"createTruckGlobal\"]}"
  invoke $1 $2 $arg "createTruckGlobal"
}

3_recordShipmentStartLocal() {
  CC_NAME=supplychain_cc_3
  echo "recordShipmentStartLocal"
  arg="{\"Args\":[\"recordShipmentStartLocal\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\"]}"
  invoke $1 $2 $arg "recordShipmentStartLocal"
}

3_recordShipmentFinalizeGlobal() {
  CC_NAME=supplychain_cc_3
  echo "recordShipmentFinalizeGlobal"
  arg="{\"Args\":[\"recordShipmentFinalizeGlobal\",\"$3\",\"$4\",\"$5\"]}"
  invoke $1 $2 $arg "recordShipmentFinalizeGlobal"
}

3_recordShipmentFinalizeLocal() {
  CC_NAME=supplychain_cc_3
  echo "recordShipmentFinalizeLocal"
  arg="{\"Args\":[\"recordShipmentFinalizeLocal\",\"$3\",\"$4\",\"$5\",\"$6\"]}"
  invoke $1 $2 $arg "recordShipmentFinalizeLocal"
}

3_queryPositionsStartLocal() {
  CC_NAME=supplychain_cc_3
  echo "queryPositionsStartLocal"
  arg="{\"Args\":[\"queryPositionsStartLocal\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\"]}"
  invoke $1 $2 $arg "queryPositionsStartLocal"
}

3_queryPositionsFinalizeGlobal() {
  CC_NAME=supplychain_cc_3
  echo "queryPositionsFinalizeGlobal"
  arg="{\"Args\":[\"queryPositionsFinalizeGlobal\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\",\"$8\"]}"
  invoke $1 $2 $arg "queryPositionsFinalizeGlobal"
}

3_queryNumberStartLocal() {
  CC_NAME=supplychain_cc_3
  echo "queryNumberStartLocal"
  arg="{\"Args\":[\"queryNumberStartLocal\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\"]}"
  invoke $1 $2 $arg "queryNumberStartLocal"
}

3_queryNumberFinalizeGlobal() {
  CC_NAME=supplychain_cc_3
  echo "queryNumberFinalizeGlobal"
  arg="{\"Args\":[\"queryNumberFinalizeGlobal\",\"$3\",\"$4\",\"$5\",\"$6\",\"$7\",\"$8\"]}"
  invoke $1 $2 $arg "queryNumberFinalizeGlobal"
}