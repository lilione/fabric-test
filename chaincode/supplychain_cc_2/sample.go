package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

type SmartContract struct {
}

type shipment struct {
	CommitInputProvider  string
	CommitOutputProvider string
	Prev                 string
	Succ                 string
}

func getCounter(stub shim.ChaincodeStubInterface, key string, inc int) string {
	cnt, _ := stub.GetState(key)
	if cnt == nil {
		cnt = []byte("0")
	}
	_cnt, _ := strconv.Atoi(string(cnt))
	stub.PutState(key, []byte(strconv.Itoa(_cnt + inc)))
	return string(cnt)
}

func putShipment(stub shim.ChaincodeStubInterface, key string, shipmentInstance shipment) {
	shipmentJSON, _ := json.Marshal(shipmentInstance)
	stub.PutState(key, shipmentJSON)
}

func getShipment(stub shim.ChaincodeStubInterface, key string) shipment {
	shipmentJSON, _ := stub.GetState(key)
	var shipmentInstance shipment
	json.Unmarshal(shipmentJSON, &shipmentInstance)
	return shipmentInstance
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "queryShipment" {
		itemID := args[0]
		seq := args[1]

		shipmentJSON, _ := stub.GetState("itemInfo" + itemID + seq)

		return shim.Success(shipmentJSON)

	} else if fn == "registerItem" {
		commitRegistrant := args[0]

		itemID := getCounter(stub,"itemID", 1)
		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		var shipmentInstance shipment
		shipmentInstance.CommitInputProvider = ""
		shipmentInstance.CommitOutputProvider = commitRegistrant
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(itemID + " " + seq))

	} else if fn == "handOffItemToNextProvider" {
		commitInputProvider := args[0]
		commitOutputProvider := args[1]
		proof := args[2]
		itemID := args[3]
		prevSeq := args[4]

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
		prevCommitOutputProvider := prevShipmentInstance.CommitOutputProvider

		if verifyEq(stub, prevCommitOutputProvider, commitInputProvider, proof) {
			seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)
			prevShipmentInstance.Succ = seq
			putShipment(stub, ("itemInfo" + itemID + prevSeq), prevShipmentInstance)

			var shipmentInstance shipment
			shipmentInstance.CommitInputProvider = commitInputProvider
			shipmentInstance.CommitOutputProvider = commitOutputProvider
			shipmentInstance.Prev = prevSeq
			putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

			return shim.Success([]byte(seq))
		}

		return shim.Error("invalid input provider")

	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
