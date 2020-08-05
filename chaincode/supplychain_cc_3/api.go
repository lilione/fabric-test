package main

import (
	"github.com/hyperledger/fabric-chaincode-go/shim"
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

func recordShipment(stub shim.ChaincodeStubInterface, truckID string, idxLoadTime string, maskedLoadTime string, idxUnloadTime string, maskedUnloadTime string) {
	chainCodeArgs := toChaincodeArgs("recordShipment", truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func queryPositions(stub shim.ChaincodeStubInterface, truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string, shares string) {
	chainCodeArgs := toChaincodeArgs("queryPositions", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}

func queryNumber(stub shim.ChaincodeStubInterface, truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string, shares string) {
	chainCodeArgs := toChaincodeArgs("queryNumber", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)
	stub.InvokeChaincode(sccName, chainCodeArgs, channelName)
}