package main

import (
	"log"
	"os"
	"os/exec"
)

func handOffItem(args string) {
	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v1/hand_off_item.py", args)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/apps/fabric/log/exec/error.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}

}

func sourceItem(args string) {
	cmd := exec.Command("python3.7", "-u", "apps/fabric/src/supplychain/v1/source_item.py", args)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	file, err := os.Create("/usr/src/HoneyBadgerMPC/apps/fabric/log/exec/error.txt")
	if err != nil {
		log.Fatal(err)
	}
	cmd.Stdout = file
	cmd.Stderr = file
	errmsg := cmd.Start()
	if errmsg != nil {
		log.Fatalf("cmd.Start() failed with %s\n", errmsg)
	}

}