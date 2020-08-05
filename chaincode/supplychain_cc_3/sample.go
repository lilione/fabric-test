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
	IdxLoadTime 	string
	IdxUnloadTime 	string
}

type truck struct {
	ShipmentList	[]shipment
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

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
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

	} else if fn == "createTruckGlobal" {
		truckID := getCounter(stub, "truckID", 1)

		var truckInstance truck
		putTruck(stub, ("truckRegistry" + truckID), truckInstance)

		return shim.Success([]byte(truckID))

	} else if fn == "recordShipmentStartLocal" {
		truckID := args[0]
		idxLoadTime := args[1]
		maskedLoadTime := args[2]
		idxUnloadTime := args[3]
		maskedUnloadTime := args[4]

		recordShipment(stub, truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)

		return shim.Success([]byte(fmt.Sprintf("recordShipmentStartLocal succeed")))

	} else if fn == "recordShipmentFinalizeGlobal" {
		truckID := args[0]
		idxLoadTime := args[1]
		idxUnloadTime := args[2]

		shipmentInstance := shipment{
			IdxLoadTime: idxLoadTime,
			IdxUnloadTime: idxUnloadTime,
		}

		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		truckInstance.ShipmentList = append(truckInstance.ShipmentList, shipmentInstance)
		putTruck(stub, ("truckRegistry" + truckID), truckInstance)

		return shim.Success([]byte("recordShipmentFinalizeGlobal succeed"))

	} else if fn == "recordShipmentFinalizeLocal" {
		idxLoadTime := args[0]
		maskedLoadTime := args[1]
		idxUnloadTime := args[2]
		maskedUnloadTime := args[3]

		shareLoadTime := calcShare(stub, idxLoadTime, maskedLoadTime)
		shareUnloadTime := calcShare(stub, idxUnloadTime, maskedUnloadTime)

		dbPut(stub, idxLoadTime, shareLoadTime)
		dbPut(stub, idxUnloadTime, shareUnloadTime)

		return shim.Success([]byte("recordShipmentFinalizeLocal succeed"))

	} else if fn == "queryPositionsStartLocal" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		shares := ""
		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		for _, shipmentInstance := range truckInstance.ShipmentList {
			shareLoadTime := dbGet(stub, shipmentInstance.IdxLoadTime)
			shareUnloadTime := dbGet(stub, shipmentInstance.IdxUnloadTime)

			shares += fmt.Sprintf("%s,%s;", shareLoadTime, shareUnloadTime)
		}
		shares = shares[:len(shares) - 1]

		queryPositions(stub, truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)

		return shim.Success([]byte("queryPositionsStartLocal succeed"))

	} else if fn == "queryPositionsFinalizeGlobal" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]
		positions := args[5]

		stub.PutState("queryPositions" + truckID + idxInitTime + maskedInitTime + idxEndTime + maskedEndTime, []byte(positions))

		return shim.Success([]byte("queryPositionsFinalizeGlobal succeed"))

	} else if fn == "queryNumberStartLocal" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]

		shares := ""
		truckInstance := getTruck(stub, ("truckRegistry" + truckID))
		for _, shipmentInstance := range truckInstance.ShipmentList {
			shareLoadTime := dbGet(stub, shipmentInstance.IdxLoadTime)
			shareUnloadTime := dbGet(stub, shipmentInstance.IdxUnloadTime)

			shares += fmt.Sprintf("%s,%s;", shareLoadTime, shareUnloadTime)
		}
		shares = shares[:len(shares) - 1]

		queryNumber(stub, truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)

		return shim.Success([]byte("queryNumberStartLocal succeed"))

	} else if fn == "queryNumberFinalizeGlobal" {
		truckID := args[0]
		idxInitTime := args[1]
		maskedInitTime := args[2]
		idxEndTime := args[3]
		maskedEndTime := args[4]
		number := args[5]

		stub.PutState("queryNumber" + truckID + idxInitTime + maskedInitTime + idxEndTime + maskedEndTime, []byte(number))

		return shim.Success([]byte("queryNumberFinalizeGlobal succeed"))

	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
