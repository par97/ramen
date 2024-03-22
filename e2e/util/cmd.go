package util

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(cmd *exec.Cmd) (error, string) {

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println("====== cmd start ======")
		fmt.Println("cmd error: " + cmd.String())
		fmt.Println(errStr)
		fmt.Println("====== cmd end ======")
		return fmt.Errorf("command failed"), errStr
	}
	if os.Getenv("e2e_debug") == "true" {
		fmt.Println("====== cmd start ======")
		fmt.Println("cmd out: " + cmd.String())
		fmt.Println(outStr)
		fmt.Println("====== cmd end ======")
	}
	return nil, outStr
}
