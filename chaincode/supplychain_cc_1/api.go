package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	channelName = "mychannel"
	sccName = "supplychain_scc_1"
)

func toChaincodeArgs(args ...string) [][]byte {
	bargs := make([][]byte, len(args))
	for i, arg := range args {
		bargs[i] = []byte(arg)
	}
	return bargs
}

func handOffItem(
	stub shim.ChaincodeStubInterface,
	idxInputProvider string,
	maskedInputProvider string,
	idxOutputProvider string,
	maskedOutputProvider string,
	idxAmt string,
	maskedAmt string,
	itemID string,
	prevSeq string,
	seq string,
	sharePrevOutputProvider string,
	sharePrevAmt string) peer.Response {

	chainCodeArgs := toChaincodeArgs("handOffItem", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)

}

func dbPut(stub shim.ChaincodeStubInterface, key string, value string) {

	chainCodeArgs := toChaincodeArgs("dbPut", key, value)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func dbGet(stub shim.ChaincodeStubInterface, key string) string {

	chainCodeArgs := toChaincodeArgs("dbGet", key)
	res := stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
	if res.Status == 200 {
		return string(res.Payload)
	}
	return ""

}

func calcShare(stub shim.ChaincodeStubInterface, idx string, maskeShare string) string {

	chainCodeArgs := toChaincodeArgs("calcShare", idx, maskeShare)
	res := stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
	if res.Status == 200 {
		return string(res.Payload)
	}
	return ""

}

//func sourceItem(stub shim.ChaincodeStubInterface, itemID string, seq string) peer.Response {
//	chainCodeArgs := toChaincodeArgs("sourceItem", itemID, seq)
//	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
//}