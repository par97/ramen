package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func RunCommand(cmd *exec.Cmd) (string, error) {

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
		return errStr, fmt.Errorf("command failed")
	}
	if os.Getenv("e2e_debug") == "true" {
		fmt.Println("====== cmd start ======")
		fmt.Println("cmd out: " + cmd.String())
		fmt.Println(outStr)
		fmt.Println("====== cmd end ======")
	}
	return outStr, nil
}

func Pause() {
	fmt.Print("Paused. Press any key to continue. ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	fmt.Println(input.Text())
}
