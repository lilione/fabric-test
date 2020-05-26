package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

func storeInput(idx string, maskedInput string) {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/calc_share.py", idx, maskedInput)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()
	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
	lines := strings.Split(outb.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "share") {
			// Very hacky way of doing this.
			shareParts := strings.Split(line, " ")
			if len(shareParts) >= 2 {
				share := shareParts[1]
				fmt.Println("The share is ", share)
				dbPut(idx, share)
			}
		}
	}
}

func reconstruct(idx string) string {
	share := dbGet(idx)
	cmd := exec.Command("python3.7", "apps/fabric/src/server/reconstruct.py", share)
	cmd.Dir = "/usr/src/HoneyBadgerMPC"
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	errmsg := cmd.Run()
	if errmsg != nil {
		log.Fatalf("cmd.Run() failed with %s\n", errmsg)
	}
	lines := strings.Split(outb.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "result") {
			resultParts := strings.Split(line, " ")
			if len(resultParts) >= 2 {
				result := resultParts[1]
				fmt.Println("The result is ", result)
				return result
			}
		}
	}
	return "None"
}