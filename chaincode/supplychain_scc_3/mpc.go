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

func cmpShare(share_a string, share_b string) string {
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
				result = result
				fmt.Println("The result is ", result)
				return result
			}
		}
	}
	return ""
}

func cmp(share_a string, share_b string) bool {
	result_share := cmpShare(share_a, share_b)
	result, _ := strconv.Atoi(reconstruct(result_share))
	return result > 0
}

func mulShare(share_a string, share_b string) string {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/mul.py", share_a, share_b)
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
				result = result
				fmt.Println("The result is ", result)
				return result
			}
		}
	}
	return ""
}

func mul(share_a string, share_b string) int {
	result_share := mulShare(share_a, share_b)
	result, _ := strconv.Atoi(reconstruct(result_share))
	return result
}

func oneMinusShare(share string) string {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/one_minus_share.py", share)
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
				result = result
				fmt.Println("The result is ", result)
				return result
			}
		}
	}
	return ""
}

func addShare(share_a string, share_b string) string {
	cmd := exec.Command("python3.7", "apps/fabric/src/server/add.py", share_a, share_b)
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
				result = result
				fmt.Println("The result is ", result)
				return result
			}
		}
	}
	return ""
}