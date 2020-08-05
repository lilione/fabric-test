package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	channelName = "mychannel"
	sccName = "supplychain_scc_3"
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

func createTruck(stub shim.ChaincodeStubInterface) peer.Response {
	chainCodeArgs := toChaincodeArgs("createTruck")
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func recordShipment(stub shim.ChaincodeStubInterface, truckID string, idxLoadTime string, maskedLoadTime string, idxUnloadTime string, maskedUnloadTime string) peer.Response {
	chainCodeArgs := toChaincodeArgs("recordShipment", truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func queryPositions(stub shim.ChaincodeStubInterface, truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string) peer.Response {
	chainCodeArgs := toChaincodeArgs("queryPositions", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func queryNumber(stub shim.ChaincodeStubInterface, truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string) peer.Response {
	chainCodeArgs := toChaincodeArgs("queryNumber", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime)
	return stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}