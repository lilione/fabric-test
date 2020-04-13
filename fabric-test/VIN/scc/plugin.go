/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct{}

// Init implements the chaincode shim interface
func (s *scc) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

// Invoke implements the chaincode shim interface
func (s *scc) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	if fn == "query" {
		if len(args) != 1 {
			return shim.Error("invalid number of args for query function")
		}
		key := args[0]
		value, err := stub.GetState(key)
		if err != nil {
			return shim.Error(fmt.Sprint("failed to get value for key %s", key))
		}
		return shim.Success([]byte(value))
	} else if fn == "update" {
		if len(args) != 2 {
			return shim.Error("invalid number of args for update function")
		}
		key, value := args[0], []byte(args[1])
		err := stub.PutState(key, value)
		if err != nil {
			return shim.Error(fmt.Sprint("failed to update for key %s value %s", key, value))
		}
		return shim.Success(nil)
	}
	return shim.Error("invalid function name.")
}

func main() {}
