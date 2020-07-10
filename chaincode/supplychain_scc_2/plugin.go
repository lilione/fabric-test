/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"

	"strconv"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct {
}

type shipment struct {
	CommitInputProvider  string
	CommitOutputProvider string
	Prev                 string
	Succ                 string
}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
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

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "registerItem" {
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

		if verify_eq(prevCommitOutputProvider, commitInputProvider, proof) {
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

	return shim.Error("invalid function name.")
}

func main() {}
