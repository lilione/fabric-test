/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// New returns an implementation of the chaincode interface
func New() shim.Chaincode {
	return &scc{}
}

type scc struct {
}

type shipment struct {
	IdxLoadTime 	string
	IdxUnloadTime 	string
}

type truck struct {
	ShipmentList	[]shipment
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

func getTruck(stub shim.ChaincodeStubInterface, key string) truck {
	truckJSON, _ := stub.GetState(key)
	var truckInstance truck
	json.Unmarshal(truckJSON, &truckInstance)
	return truckInstance
}

func putTruck(stub shim.ChaincodeStubInterface, key string, truckInstance truck) {
	truckJSON, _ := json.Marshal(&truckInstance)
	stub.PutState(key, truckJSON)
}

func inRange(x string, l string, r string) bool {
	fmt.Println("inRange")
	if cmp(x, l) {
		fmt.Println("cmp(x, l)")
		return false
	}
	if cmp(r, x) {
		fmt.Println("cmp(r, x)")
		return false
	}
	fmt.Println("pass")
	return true
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

	} else if fn == "createTruck" {

		truckID := getCounter(stub, "truckID", 1)

		var truckInstance truck
		putTruck(stub, ("truckRegistry" + truckID), truckInstance)

		return shim.Success([]byte(truckID))

	} else if fn == "recordShipment" {

		truckID := args[0]
		idxLoadTime := args[1]
		maskedLoadTime := args[2]
		idxUnloadTime := args[3]
		maskedUnloadTime := args[4]

		shareLoadTime := calcShare(idxLoadTime, maskedLoadTime)
		shareUnloadTime := calcShare(idxUnloadTime, maskedUnloadTime)

		dbPut(idxLoadTime, shareLoadTime)
		dbPut(idxUnloadTime, shareUnloadTime)

		if !cmp(shareUnloadTime, shareLoadTime) {
			shipmentInstance := shipment{
				IdxLoadTime: idxLoadTime,
				IdxUnloadTime: idxUnloadTime,
			}

			truckInstance := getTruck(stub, ("truckRegistry" + truckID))
			truckInstance.ShipmentList = append(truckInstance.ShipmentList, shipmentInstance)
			putTruck(stub, ("truckRegistry" + truckID), truckInstance)

			return shim.Success([]byte("recorded successfully"))
		}

		return shim.Error("invalid load and unlaod time")

	} else if fn == "queryPositions" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		shareInitTime := calcShare(idxInitTime, maskedInitTime)
		shareEndTime := calcShare(idxEndTime, maskedEndTime)

		var positions string
		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		fmt.Println("len", len(truckInstance.ShipmentList))
		for index, shipmentInstance := range truckInstance.ShipmentList {
			fmt.Println("index", index)
			shareLoadTime := dbGet(shipmentInstance.IdxLoadTime)
			shareUnloadTime := dbGet(shipmentInstance.IdxUnloadTime)

			if inRange(shareLoadTime, shareInitTime, shareEndTime) || inRange(shareUnloadTime, shareInitTime, shareEndTime) || inRange(shareInitTime, shareLoadTime, shareUnloadTime) {
				fmt.Println("in")
				if positions != "" {
					fmt.Println("add blank")
					positions += " "
				}
				positions += strconv.Itoa(index)
			}
			fmt.Println("positions", positions)
		}
		fmt.Println("finished")
		fmt.Println([]byte(positions))
		return shim.Success([]byte(positions))

	}

	return shim.Error("invalid function name.")
}

func main() {
}
