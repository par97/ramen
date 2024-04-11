package util

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/go-logr/logr"
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

func CurrentFuncName() string {
	counter, _, _, success := runtime.Caller(2)
	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	fullPathName := runtime.FuncForPC(counter).Name()
	split := strings.Split(fullPathName, "/")
	return split[len(split)-1]
}

func LogEnter(log *logr.Logger) {
	log.Info("🚀 " + CurrentFuncName())
}

func LogExit(log *logr.Logger) {
	log.Info("👍 " + CurrentFuncName())
}
