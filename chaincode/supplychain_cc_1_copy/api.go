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

func getInputmaskIdx(stub shim.ChaincodeStubInterface, num string) peer.Response {
	chainCodeArgs := toChaincodeArgs("getInputmaskIdx", num)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func registerItem(stub shim.ChaincodeStubInterface, idxRegistrant string, maskedRegistrant string, idxAmt string, maskedAmt string) peer.Response {
	chainCodeArgs := toChaincodeArgs("registerItem", idxRegistrant, maskedRegistrant, idxAmt, maskedAmt)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func handOffItemToNextProvider(stub shim.ChaincodeStubInterface, idxInputProvider string, maskedInputProvider string, idxOutputProvider string, maskedOutputProvider string, idxAmt string, maskedAmt string, itemID string, prevSeq string) peer.Response {
	chainCodeArgs := toChaincodeArgs("handOffItemToNextProvider", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func sourceItem(stub shim.ChaincodeStubInterface, itemID string, seq string) peer.Response {
	chainCodeArgs := toChaincodeArgs("sourceItem", itemID, seq)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}