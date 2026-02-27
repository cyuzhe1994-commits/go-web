package test

import (
	"testing"
	"time"

	go_web "github.com/go-web"
)

// go test -run TestKernal -v

func TestKernalDefault(t *testing.T) {
	engin := go_web.NewEngine(nil)
	go engin.Run(":8080")
}

func TestKernalWithLogger(t *testing.T) {
	logger := go_web.NewDefaultFrameWorkLog(go_web.DefaultFrameWorkLogLevelWarn, time.Local)
	engin := go_web.NewEngine(logger)
	go engin.Run(":8080")
}
