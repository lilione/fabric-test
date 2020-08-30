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

	if fn == "verify" {
		prevProvider := args[0]
		sucProvider := args[1]
		proofProvider := args[2]
		prevAmt := args[3]
		sucAmt := args[4]
		proofAmt := args[5]

		if verify(prevProvider, sucProvider, proofProvider, prevAmt, sucAmt, proofAmt) {
			return shim.Success(nil)
		}

		return shim.Error("")

	}

	return shim.Error("invalid function name.")
}

func main() {}
