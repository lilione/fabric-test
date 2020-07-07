package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func verify_eq(commitPrev string, commitSuc string, proof string) bool {
	fmt.Println(commitPrev, commitSuc, proof)
	cmd := exec.Command("python3.7", "apps/fabric/src/server/verify_eq.py", commitPrev, commitSuc, proof)
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
				result, _ := strconv.Atoi(resultParts[1])
				fmt.Println("The result is ", result)
				return result > 0
			}
		}
	}
	return false
}
