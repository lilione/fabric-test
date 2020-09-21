/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	//"github.com/hyperledger/fabric-chaincode-go/shim"
	//"github.com/hyperledger/fabric-protos-go/peer"
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

	if fn == "dbPut" {
		key := args[0]
		value := args[1]

		dbPut(key, value)

		return shim.Success([]byte(""))

	} else if fn == "dbGet" {
		key := args[0]

		value := dbGet(key)

		return shim.Success(value)

	} else if fn == "calcShare" {
		idx := args[0]
		maskedShare := args[1]

		value := calcShare(idx, maskedShare)

		return shim.Success([]byte(value))

	} else if fn == "handOffItem" {
		handOffItem(args[0])

	} else if fn == "sourceItem" {
		sourceItem(args[0])
	}

	return shim.Error("invalid function name.")
}

func main() {}
