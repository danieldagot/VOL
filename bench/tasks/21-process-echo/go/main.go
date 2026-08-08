package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	out, err := exec.Command("echo", "vol").Output()
	status := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			status = ee.ExitCode()
		} else {
			panic(err)
		}
	}
	fmt.Println(status)
	fmt.Println(strings.TrimSpace(string(out)))
}
