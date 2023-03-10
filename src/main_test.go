package main

import (
	"os"
	"runtime/debug"
	"testing"
)

var config = `
listen: 0.0.0.0:5300
upstream:
  - 8.8.8.8:53
  - 8.8.4.4:53
blacklist:
  - https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
whitelist:
  - googleadservices.com
  - iadsdk.apple.com
local:
  - name: service1.local
    target:
    - address: 192.168.178.1
      type: A
    - address: 192.168.179.1
      type: A
    - address: fe80::9656:d028:8652:1111
      type: AAAA
  - name: service2.local
    target:
    - address: 192.168.178.2
      type: A
`

func TestMain(m *testing.M) {
	GetConfig().ReadConfigData([]byte(config))
	initServer()
	code := m.Run()
	os.Exit(code)
}

func checkTestBool(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Fatalf("Expected '%t', but got '%t' at:\n%s", expected, actual, debug.Stack())
	}
}

func checkTestInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected '%d', but got '%d' at:\n%s", expected, actual, debug.Stack())
	}
}

func checkTestString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%s', but got '%s' at:\n%s", expected, actual, debug.Stack())
	}
}
