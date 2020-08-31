package main

import (
	"bytes"
	"log"
	"os/exec"
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
				return share
			}
		}
	}
	return ""
}