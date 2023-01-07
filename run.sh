#!/bin/sh
LISTEN_ADDR=0.0.0.0:5300 go run `ls *.go | grep -v _test.go`