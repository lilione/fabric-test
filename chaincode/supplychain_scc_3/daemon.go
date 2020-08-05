package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func recordShipment(truckID string, idxLoadTime string, maskedLoadTime string, idxUnloadTime string, maskedUnloadTime string) {
	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v3/record_shipment.py", truckID, idxLoadTime, maskedLoadTime, idxUnloadTime, maskedUnloadTime)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	fmt.Println(cmd)
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}
}

func queryPositions(truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string, shares string) {
	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v3/query_positions.py", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	fmt.Println(cmd)
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}
}

func queryNumber(truckID string, idxInitTime string, maskedInitTime string, idxEndTime string, maskedEndTime string, shares string) {
	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v3/query_number.py", truckID, idxInitTime, maskedInitTime, idxEndTime, maskedEndTime, shares)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/log.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	fmt.Println(cmd)
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}
}