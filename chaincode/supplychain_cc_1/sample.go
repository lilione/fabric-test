package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

const (
	READY = "READY"
	SETTLE = "SETTLE"
)

type SmartContract struct {
}

type shipment struct {
	IdxInputProvider  	string
	IdxOutputProvider 	string
	IdxAmount         	string
	Prev              	string
	Succs             	string
	State				string
}

type inquiry struct {
	Value 	string
	State	string
}

func (s *SmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
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
	return shipmentInstance
}

func putShipment(stub shim.ChaincodeStubInterface, key string, shipmentInstance shipment) {
	shipmentJSON, _ := json.Marshal(&shipmentInstance)
	stub.PutState(key, shipmentJSON)
}

func delShipment(stub shim.ChaincodeStubInterface, key string) {
	stub.DelState(key)
}

func getInquiry(stub shim.ChaincodeStubInterface, key string) inquiry {
	inquiryJSON, _ := stub.GetState(key)
	var inquiryInstance inquiry
	json.Unmarshal(inquiryJSON, &inquiryInstance)
	return inquiryInstance
}

func putInquiry(stub shim.ChaincodeStubInterface, key string, inquiryInstance inquiry) {
	inquiryJSON, _ := json.Marshal(&inquiryInstance)
	stub.PutState(key, inquiryJSON)
}
func (s *SmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()

	if fn == "queryShipment" {
		itemID := args[0]
		seq := args[1]

		shipmentJSON, _ := stub.GetState(("itemInfo" + itemID + seq))

		return shim.Success(shipmentJSON)

	} else if fn == "queryInquiry" {
		itemID := args[0]
		seq := args[1]

		inquiryJSON, _ := stub.GetState(("sourceItem" + itemID + seq))

		return shim.Success(inquiryJSON)

	} else if fn == "getInputmaskIdx" {
		num, _ := strconv.Atoi(args[0])

		var inputmaskIdx string
		value, _ := strconv.Atoi(getCounter(stub, "inputmaskCnt", num))
		for num > 0 {
			inputmaskIdx += strconv.Itoa(value)
			value += 1
			num -= 1
			if num > 0 {
				inputmaskIdx += " "
			}
		}

		return shim.Success([]byte(inputmaskIdx))

	} else if fn == "registerItemClientGlobal" {
		data := strings.Split(args[0], ",")
		paraNum := 2
		if len(data) % paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		itemID, _ := strconv.Atoi(getCounter(stub, "itemID", itemNum))
		ret := ""
		for i := 0; i < itemNum; i++ {
			idxRegistrant := data[paraNum*i]
			idxAmt := data[paraNum*i+1]

			strItemID := strconv.Itoa(itemID)
			seq := getCounter(stub, ("itemShipmentCnt" + strItemID), 1)

			var shipmentInstance shipment
			shipmentInstance.IdxOutputProvider = idxRegistrant
			shipmentInstance.IdxAmount = idxAmt
			shipmentInstance.State = SETTLE
			putShipment(stub, ("itemInfo" + strItemID + seq), shipmentInstance)

			if i > 0 {
				ret += " "
			}
			ret += fmt.Sprintf("%s %s", strItemID, seq)

			itemID += 1
		}

		return shim.Success([]byte(ret))

	} else if fn == "registerItemClientLocal" {
		data := strings.Split(args[0], ",")
		paraNum := 4
		if len(data) % paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		args := ""
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum*i]
			seq := data[paraNum*i+1]
			maskedRegistrant := data[paraNum*i+2]
			maskedAmt := data[paraNum*i+3]

			shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
			if shipmentInstance.State != SETTLE {
				return shim.Error("Shipment hasn't settled yet")
			}

			idxRegistrant := shipmentInstance.IdxOutputProvider
			idxAmt := shipmentInstance.IdxAmount

			if i > 0 {
				args += " "
			}
			args += fmt.Sprintf("%s %s %s %s", idxRegistrant, maskedRegistrant, idxAmt, maskedAmt)
		}

		registerItem(stub, args)

		return shim.Success([]byte("registerItemClientLocal finished"))

	} else if fn == "handOffItemClientGlobal" {
		data := strings.Split(args[0], ",")
		paraNum := 5
		if len(data)%paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		set := make(map[string]bool)

		itemNum := len(data) / paraNum
		ret := ""
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum*i]
			prevSeq := data[paraNum*i+1]
			idxInputProvider := data[paraNum*i+2]
			idxOutputProvider := data[paraNum*i+3]
			idxAmt := data[paraNum*i+4]

			_, ok := set[itemID]
			if ok {
				return shim.Error(fmt.Sprintf("Repetitive itemID %s", itemID))
			}
			set[itemID] = true

			prevKey := "itemInfo" + itemID + prevSeq
			prevShipmentInstance := getShipment(stub, prevKey)
			if prevShipmentInstance.State != SETTLE {
				return shim.Error("Previous shipment hasn't settled yet")
			}
			if prevShipmentInstance.Succs != "" {
				return shim.Error("Previous shipment already has a successor")
			}

			seq := getCounter(stub, ("itemShipmentCnt" + itemID), 1)
			if i > 0 {
				ret += " "
			}
			ret += seq

			prevShipmentInstance.Succs = seq
			putShipment(stub, prevKey, prevShipmentInstance)

			var shipmentInstance shipment
			shipmentInstance.IdxInputProvider = idxInputProvider
			shipmentInstance.IdxOutputProvider = idxOutputProvider
			shipmentInstance.IdxAmount = idxAmt
			shipmentInstance.Prev = prevSeq
			shipmentInstance.State = READY
			putShipment(stub, ("itemInfo" + itemID + seq), shipmentInstance)
		}

		return shim.Success([]byte(ret))

	} else if fn == "handOffItemClientLocal" {
		data := strings.Split(args[0], ",")
		paraNum := 5
		if len(data)%paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		args := ""
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum*i]
			seq := data[paraNum*i+1]
			maskedInputProvider := data[paraNum*i+2]
			maskedOutputProvider := data[paraNum*i+3]
			maskedAmt := data[paraNum*i+4]

			shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
			if shipmentInstance.State != READY {
				return shim.Error("Shipment is not ready")
			}

			idxInputProvider := shipmentInstance.IdxInputProvider
			idxOutputProvider := shipmentInstance.IdxOutputProvider
			idxAmt := shipmentInstance.IdxAmount

			prevSeq := shipmentInstance.Prev
			prevShipmentInstance := getShipment(stub, ("itemInfo" + itemID + prevSeq))
			prevIdxOutputProvider := prevShipmentInstance.IdxOutputProvider
			prevIdxAmt := prevShipmentInstance.IdxAmount

			if i > 0 {
				args += " "
			}
			args += fmt.Sprintf("%s %s %s %s %s %s %s %s %s %s", itemID, seq, idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, prevIdxOutputProvider, prevIdxAmt)
		}

		handOffItem(stub, args)

		return shim.Success([]byte("handOffItemClientLocal finished"))

	} else if fn == "handOffItemServerGlobal" {
		data := strings.Split(args[0], ",")
		paraNum := 3
		if len(data)%paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum*i]
			seq := data[paraNum*i+1]
			result := data[paraNum*i+2]

			key := "itemInfo" + itemID + seq
			shipmentInstance := getShipment(stub, key)
			if result == "pass" {
				shipmentInstance.State = SETTLE
				putShipment(stub, key, shipmentInstance)
			} else if result == "fail" {
				prevSeq := shipmentInstance.Prev
				delShipment(stub, key)

				prevKey := "itemInfo" + itemID + prevSeq
				prevShipmentInstance := getShipment(stub, prevKey)
				prevShipmentInstance.Succs = ""
				putShipment(stub, prevSeq, prevShipmentInstance)
			}
		}

		return shim.Success([]byte("handOffItemServerGlobal finished"))

	} else if fn == "sourceItemClientLocal" {
		data := strings.Split(args[0], ",")
		paraNum := 2
		if len(data)%paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		args := ""
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum*i]
			seq := data[paraNum*i+1]

			inquiryInstance := getInquiry(stub, ("sourceItem" + itemID + seq))
			if inquiryInstance.State == SETTLE {
				continue
			}

			if i > 0 {
				args += " "
			}
			args += fmt.Sprintf("%s %s ", itemID, seq)

			shareIdxes := ""
			if seq != "0" {
				for true {
					shipmentInstance := getShipment(stub, ("itemInfo" + itemID + seq))
					shareIdxes += shipmentInstance.IdxInputProvider
					seq = shipmentInstance.Prev
					if seq == "0" {
						break
					}
					shareIdxes += ","
				}
			}

			args += shareIdxes
		}

		sourceItem(stub, args)

		return shim.Success([]byte("sourceItemClientLocal finished"))

	} else if fn == "sourceItemServerGlobal" {
		data := strings.Split(args[0], ",")
		paraNum := 3
		if len(data)%paraNum != 0 {
			shim.Error("Invalid number of arguments")
		}

		itemNum := len(data) / paraNum
		for i := 0; i < itemNum; i++ {
			itemID := data[paraNum * i]
			seq := data[paraNum * i + 1]
			listInputProvider := data[paraNum * i + 2]

			key := "sourceItem" + itemID + seq
			inquiryInstance := getInquiry(stub, key)
			inquiryInstance.Value = listInputProvider
			inquiryInstance.State = SETTLE
			putInquiry(stub, key, inquiryInstance)
		}

		return shim.Success([]byte("sourceItemServerGlobal finished"))

	}

	return shim.Error("Invalid Smart Contract function name.")
}

func main() {
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
