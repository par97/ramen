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
		fmt.Println("============")
		fmt.Println("error of: " + cmd.String())
		fmt.Println(errStr)
		fmt.Println("============")
		return fmt.Errorf("command failed"), errStr
	}
	if os.Getenv("e2e_debug") == "true" {
		fmt.Println("============")
		fmt.Println("output of: " + cmd.String())
		fmt.Println(outStr)
		fmt.Println("============")
	}
	return nil, outStr
}
