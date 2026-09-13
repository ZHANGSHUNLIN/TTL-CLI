package config

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testHome, err := os.MkdirTemp("", "ttl-conf-test-home-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.Setenv("HOME", testHome); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := os.Setenv("USERPROFILE", testHome); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	code := m.Run()
	if err := os.RemoveAll(testHome); err != nil && code == 0 {
		fmt.Fprintln(os.Stderr, err)
		code = 1
	}
	os.Exit(code)
}
