package test

import (
	"github.com/taoey/go-log-listener/listener"
	"testing"
	"time"
)

// 测试默认日志监听器
func TestDefaultLogListener(t *testing.T) {
	filePath := "E:\\projects\\go-mod\\log-listener\\go.sum"
	logListener := listener.NewDefaultLogListener(filePath, 3)
	logListener.Run()

	time.Sleep(time.Second * 8)

	logListener.Stop()
	time.Sleep(time.Hour * 1)
}
