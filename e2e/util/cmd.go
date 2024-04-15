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

func RunCommand(log *logr.Logger, cmd *exec.Cmd) (string, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := stderr.String()
		log.Error(err, "====== cmd start ======\n"+
			"cmd error: "+cmd.String()+"\n"+
			errStr+"\n"+
			"====== cmd end ======")

		return errStr, fmt.Errorf("command failed")
	}

	outStr := stdout.String()
	if os.Getenv("e2e_debug") == "true" {
		log.Info("====== cmd start ======\n" +
			"cmd out: " + cmd.String() + "\n" +
			outStr + "\n" +
			"====== cmd end ======")
	}

	return outStr, nil
}

func Pause(log *logr.Logger) {
	log.Info("Paused. Press any key to continue.")

	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	// fmt.Println(input.Text())
}

func CurrentFuncName() string {
	countback := 2
	counter, _, _, success := runtime.Caller(countback)

	if !success {
		// println("functionName: runtime.Caller: failed")
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
