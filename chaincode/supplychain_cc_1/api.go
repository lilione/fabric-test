package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
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
	sharePrevAmt string) {

	chainCodeArgs := toChaincodeArgs("handOffItem", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)

}

func sourceItem(stub shim.ChaincodeStubInterface, itemID string, seq string, shares string) {

	chainCodeArgs := toChaincodeArgs("sourceItem", itemID, seq, shares)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)

}