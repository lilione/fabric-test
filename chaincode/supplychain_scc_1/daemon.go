package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func handOffItem(
	idxInputProvider string,
	maskedInputProvider string,
	idxOutputProvider string,
	maskedOutputProvider string,
	idxAmt string,
	maskedAmt string,
	itemID string,
	prevSeq string,
	seq string,
	sharePrevOutputProvider string,
	sharePrevAmt string) {

	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v1/hand_off_item.py", idxInputProvider, maskedInputProvider, idxOutputProvider, maskedOutputProvider, idxAmt, maskedAmt, itemID, prevSeq, seq, sharePrevOutputProvider, sharePrevAmt)
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

func sourceItem(itemID string, seq string, shareInputProvider string) {

	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v1/source_item.py", itemID, seq, shareInputProvider)
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