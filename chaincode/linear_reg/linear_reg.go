package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type eq struct{}

type LinReg struct {
	ObjectType   string   `json:"docType"`
	Name         string   `json:"name"`
	EndPoints    int      `json:endPoints`
	Participants []string `json:"participants"`
	X            []string `json:"X"`
	Y            []string `json:"Y"`
	M            int      `json:"M"`
	B            int      `json:"B"`
	Result       string   `json:"Result"`
}

func getInstance(key string, stub shim.ChaincodeStubInterface) LinReg {
	item, _ := stub.GetState(key)
	var LinRegInstance LinReg
	json.Unmarshal([]byte(item), &LinRegInstance)
	return LinRegInstance
}

func (instance LinReg) save(stub shim.ChaincodeStubInterface) {
	JSONasBytes, _ := json.Marshal(instance)
	stub.PutState(instance.Name, JSONasBytes)
}

// Init
func (t *eq) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func createCell(stub shim.ChaincodeStubInterface, key string, userID string, namespace string) {
	chainCodeArgs := ToChaincodeArgs("createCell", key, userID, namespace)
	response := stub.InvokeChaincode("honeybadgerscc", chainCodeArgs, "mychannel")

	if response.Status != shim.OK {
		fmt.Println(response.Message)
	}
}

func getActiveLinRegs(stub shim.ChaincodeStubInterface) peer.Response {
	keysIter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var items []LinReg
	for keysIter.HasNext() {
		key, _ := keysIter.Next()
		var LinRegInstance LinReg
		LinRegInstance = getInstance(key.Key, stub)
		if LinRegInstance.M == -1 {
			items = append(items, LinRegInstance)
		}
	}
	// LinRegInstance.save()
	JSONasBytes, _ := json.Marshal(items)
	return shim.Success([]byte(JSONasBytes))

}

func getCompletedLinRegs(stub shim.ChaincodeStubInterface) peer.Response {
	keysIter, err := stub.GetStateByRange("", "")
	if err != nil {
		return shim.Error(fmt.Sprintf("keys operation failed. Error accessing state: %s", err))
	}
	defer keysIter.Close()

	var items []LinReg
	for keysIter.HasNext() {
		key, _ := keysIter.Next()
		var LinRegInstance LinReg
		LinRegInstance = getInstance(key.Key, stub)
		if LinRegInstance.M != -1 {
			items = append(items, LinRegInstance)
		}
	}
	// LinRegInstance.save()
	JSONasBytes, _ := json.Marshal(items)
	return shim.Success([]byte(JSONasBytes))

}

// Invoke
func (t *eq) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Extract the function and args from the transaction proposal
	fn, args := stub.GetFunctionAndParameters()
	if fn == "createInstance" && len(args) >= 3 {
		fmt.Println("In createInstance endpoint ")
		// Create a LinReg with the parameters provided
		// args[0] - name of the LinReg
		// args[1] - time to end the LinReg in seconds
		// args[2] - name of the user
		// returns memory cell id for the user move
		endParams, err := strconv.Atoi(args[1])
		if err != nil {
			return shim.Error(err.Error())
		}

		objectType := "LinReg"

		// be a hash or something non-deterministic
		// createCell(stub, memcell, args[2], "eq") // create a memory cell for the user to enter data
		LinReg := &LinReg{objectType, args[0], endParams, []string{args[2]}, []string{}, []string{}, -1, -1, "None"}
		LinReg.save(stub)
		return shim.Success([]byte("Success"))
	} else if fn == "getActiveLinRegs" {

		return getActiveLinRegs(stub)

	} else if fn == "getCompletedLinRegs" {

		return getCompletedLinRegs(stub)

	} else if fn == "addData" && len(args) >= 2 {
		fmt.Println("In joinLinReg endpoint")
		// Add data to an instance
		// args[0] - name of the LinReg
		// args[1] - name of the user
		// returns memory cell id for the user move
		key := args[0]
		var LinRegInstance LinReg
		LinRegInstance = getInstance(key, stub)
		memcellX := "linReg" + args[0] + strconv.Itoa(len(LinRegInstance.X)) + "X_cell"
		memcellY := "linReg" + args[0] + strconv.Itoa(len(LinRegInstance.Y)) + "Y_cell"
		createCell(stub, memcellX, args[1], "linReg")
		createCell(stub, memcellY, args[1], "linReg")
		LinRegInstance.X = append(LinRegInstance.X, memcellX)
		LinRegInstance.Y = append(LinRegInstance.Y, memcellY)
		LinRegInstance.save(stub)
		return shim.Success([]byte(memcellX + ";" + memcellY))
	} else if fn == "runMPC" && len(args) >= 1 {
		var LinRegInstance LinReg
		fmt.Println("In  run mpc operation")
		LinRegAsBytes, _ := stub.GetState(args[0])
		json.Unmarshal([]byte(LinRegAsBytes), &LinRegInstance)
		fmt.Println(len(LinRegInstance.X))
		fmt.Println(len(LinRegInstance.Y))
		if len(LinRegInstance.X) < LinRegInstance.EndPoints && len(LinRegInstance.Y) < LinRegInstance.EndPoints {
			return shim.Success([]byte("None"))
		}
		fmt.Println("passed number check")
		for _, element := range LinRegInstance.X {
			var cell secretcell
			cell = getCellMetaData(stub, element)
			if cell.IsWritten == false {
				fmt.Println(cell.CellName + "not written")
				return shim.Success([]byte("None"))
			}
		}
		for _, element := range LinRegInstance.Y {
			var cell secretcell
			cell = getCellMetaData(stub, element)
			if cell.IsWritten == false {
				fmt.Println(cell.CellName + "not written")
				return shim.Success([]byte("None"))
			}
		}
		fmt.Println("Started MPC operation")
		cells := append(LinRegInstance.X, LinRegInstance.Y...)

		mpcOp(stub, "linear_regression_mpc", LinRegInstance.Name, cells...)
		return shim.Success([]byte("Started"))
	} else if fn == "endLinReg" && len(args) >= 2 {
		// end the LinReg
		// args[0] - name of the LinReg
		var LinRegInstance LinReg
		fmt.Println("In endLinReg ")
		LinRegAsBytes, err := stub.GetState(args[0])
		if err != nil {
			return shim.Error(err.Error())
		}

		err = json.Unmarshal([]byte(LinRegAsBytes), &LinRegInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		// check if both cells are ready to be opened
		if checkMPCResult(stub, LinRegInstance.Name, "LinReg") == "None" {
			return shim.Success([]byte("None"))
		}
		LinRegInstance.Result = checkMPCResult(stub, LinRegInstance.Name, "LinReg")
		LinRegJSONasBytes, err := json.Marshal(LinRegInstance)
		if err != nil {
			return shim.Error(err.Error())
		}
		err = stub.PutState(args[0], []byte(LinRegJSONasBytes))
		if err != nil {
			return shim.Error(err.Error())
		}
		return shim.Success([]byte(LinRegInstance.Result))
	}
	return shim.Success([]byte("Invalid endpoint"))
}

// main function starts up the chaincode in the container during instantiate
func main() {
	if err := shim.Start(new(eq)); err != nil {
		fmt.Printf("Error starting eq chaincode: %s", err)
	}
}
