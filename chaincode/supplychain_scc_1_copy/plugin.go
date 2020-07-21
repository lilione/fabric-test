/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"

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
	IdxInputProvider  string
	IdxOutputProvider string
	IdxAmount         string
	Prev              string
	Succs             []string
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

func getShipment(stub shim.ChaincodeStubInterface, key string) shipment {
	shipmentJSON, _ := stub.GetState(key)
	var shipmentInstance shipment
	json.Unmarshal(shipmentJSON, &shipmentInstance)
	fmt.Println("getShipment", key, shipmentInstance, shipmentJSON)
	return shipmentInstance
}

func putShipment(stub shim.ChaincodeStubInterface, key string, shipmentInstance shipment) {
	shipmentJSON, _ := json.Marshal(&shipmentInstance)
	fmt.Println("putShipment", key, shipmentInstance, shipmentJSON)
	stub.PutState(key, shipmentJSON)
}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "getInputmaskIdx" {
		num, _ := strconv.Atoi(args[0])

		var inputmaskIdx string
		value, _ := strconv.Atoi(getCounter(stub,"inputmaskCnt", num))
		for num > 0 {
			inputmaskIdx += strconv.Itoa(value)
			value += 1
			num -= 1
			if num > 0 {
				inputmaskIdx += " "
			}
		}

		return shim.Success([]byte(inputmaskIdx))

	} else if fn == "registerItem" {
		idxRegistrant := args[0]
		maskedRegistrant := args[1]
		idxAmt := args[2]
		maskedAmt := args[3]

		itemID := getCounter(stub,"itemID", 1)
		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)

		dbPut(idxRegistrant, calcShare(idxRegistrant, maskedRegistrant))
		dbPut(idxAmt, calcShare(idxAmt, maskedAmt))

		var shipmentInstance shipment
		shipmentInstance.IdxInputProvider = ""
		shipmentInstance.IdxOutputProvider = idxRegistrant
		shipmentInstance.IdxAmount = idxAmt
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(itemID + " " + seq))

	} else if fn == "handOffItemToNextProvider" {
		idxInputProvider := args[0]
		maskedInputProvider := args[1]
		idxOutputProvider := args[2]
		maskedOutputProvider := args[3]
		idxAmt := args[4]
		maskedAmt := args[5]
		itemID := args[6]
		prevSeq := args[7]

		prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))

		inputProvider := calcShare(idxInputProvider, maskedInputProvider)
		fmt.Println("before")
		fmt.Println(prevShipmentInstance)
		prevOutputProvider := dbGet(prevShipmentInstance.IdxOutputProvider)
		fmt.Println("after")
		if !eq(inputProvider, prevOutputProvider) {
			return shim.Error("Invalid input provider")
		}
		amt := calcShare(idxAmt, maskedAmt)
		prevAmt := dbGet(prevShipmentInstance.IdxAmount)
		if cmp(prevAmt, amt) {
			return shim.Error("Invalid amount")
		}

		_prevAmt, _ := strconv.Atoi(prevAmt)
		_amt, _ := strconv.Atoi(amt)
		dbPut(prevShipmentInstance.IdxAmount, strconv.Itoa(_prevAmt - _amt))
		seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)
		prevShipmentInstance.Succs = append(prevShipmentInstance.Succs, seq)
		putShipment(stub, ("itemInfo" + itemID + prevSeq), prevShipmentInstance)

		dbPut(idxInputProvider, inputProvider)
		outputProvider := calcShare(idxOutputProvider, maskedOutputProvider)
		dbPut(idxOutputProvider, outputProvider)
		dbPut(idxAmt, amt)

		var shipmentInstance shipment
		shipmentInstance.IdxInputProvider = idxInputProvider
		shipmentInstance.IdxOutputProvider = idxOutputProvider
		shipmentInstance.IdxAmount = idxAmt
		shipmentInstance.Prev = prevSeq
		putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)

		return shim.Success([]byte(seq))

	} else if fn == "sourceItem" {
		itemID := args[0]
		seq := args[1]

		var providers string
		for true {
			shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
			providers += reconstruct(dbGet(shipmentInstance.IdxOutputProvider))
			seq = shipmentInstance.Prev
			if seq == "" {
				break
			} else {
				providers += " "
			}
		}
		return shim.Success([]byte(providers))

	}

	return shim.Error("invalid function name.")
}

func main() {}
