package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

func calcShare(idx string, maskedShare string) string {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/calc_share.py", idx, maskedShare)
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
			shareParts := strings.Split(line, " ")
			if len(shareParts) >= 2 {
				share := shareParts[1]
				fmt.Println("The share is ", share)
				return share
			}
		}
	}
	return ""
}

func reconstruct(share string) string {
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
		if strings.Contains(line, "value") {
			valueParts := strings.Split(line, " ")
			if len(valueParts) >= 2 {
				value := valueParts[1]
				fmt.Println("The value is ", value)
				return value
			}
		}
	}
	return ""
}

func cmp(share_a string, share_b string) bool {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/cmp.py", share_a, share_b)
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
				res, _ := strconv.Atoi(result)
				return (res > 0)
			}
		}
	}
	return false
}

func eq(share_a string, share_b string) bool {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/eq.py", share_a, share_b)
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
				res, _ := strconv.Atoi(result)
				return (res > 0)
			}
		}
	}
	return false
}