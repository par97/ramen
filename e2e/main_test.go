package main

import (
	"os"
	"testing"

	"github.com/ramendr/ramen/e2e/suites"
)

func TestMain(m *testing.M) {
	setup()

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestPrecheckSuite(t *testing.T) {
	err := RunSuite(&suites.PrecheckSuite{}, &ctx)
	if err != nil {
		panic(err)
	}
}

func TestBasicSuite(t *testing.T) {
	err := RunSuite(&suites.BasicSuite{}, &ctx)
	if err != nil {
		panic(err)
	}
}

func TestAppSetSuite(t *testing.T) {
	err := RunSuite(&suites.AppSetSuite{}, &ctx)
	if err != nil {
		panic(err)
	}
}
