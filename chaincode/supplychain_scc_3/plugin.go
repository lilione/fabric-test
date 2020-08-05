/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct {
}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "calcShare" {
		idx := args[0]
		maskedShare := args[1]

		value := calcShare(idx, maskedShare)

		return shim.Success([]byte(value))
	} else if fn == "dbPut" {
		key := args[0]
		value := args[1]

		dbPut(key, value)

		return shim.Success([]byte(""))

	} else if fn == "dbGet" {
		key := args[0]

		value := dbGet(key)

		return shim.Success([]byte(value))

	} else if fn == "recordShipment" {
		truckID := args[0]
		idxLoadTime := args[1]
		maskedLoadTime := args[2]
		idxUnloadTime := args[3]
		maskedUnloadTime := args[4]

		recordShipment(truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)

	} else if fn == "queryPositions" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]
		shares := args[5]

		queryPositions(truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)

	} else if fn == "queryNumber" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]
		shares := args[5]

		queryNumber(truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)
	}

	return shim.Error("invalid function name.")
}

func main() {
}
