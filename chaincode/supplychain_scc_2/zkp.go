
package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func writeToFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func verify(prevProvider string, sucProvider string, proofProvider string, prevAmt string, sucAmt string, proofAmt string) bool {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/verify.py", prevProvider, sucProvider, proofProvider, prevAmt, sucAmt, proofAmt)
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
				return result > 0
			}
		} else if strings.Contains(line, "exe_time") {
			err := writeToFile("/usr/src/HoneyBadgerMPC/time.log", strings.Split(line, " ")[1] + "\n")
			if err != nil {
				log.Fatalf("Write to time.log failed with %s\n", err)
			}
		}
	}
	return false
}
