package test

import (
	"fmt"
	"github.com/taoey/go-log-listener/listener"
	"runtime"
	"testing"
	"time"
)

// 测试默认日志监听器
func TestDefaultLogListener(t *testing.T) {
	filePath := "E:\\projects\\go-mod\\log-listener\\go.sum"
	logListener := listener.NewDefaultLogListener(filePath, 3)

	varLogHandler := func(logStr string) interface{} {
		fmt.Println("日志处理：", logStr)
		return logStr
	}

	varStorageHandler := func(logObject interface{}) {
		fmt.Println("日志存储:", logObject)
	}

	logListener.SetHandler(varLogHandler, varStorageHandler)

	logListener.Run()

	time.Sleep(time.Second * 8)
	fmt.Println("runtime goroutine number:", runtime.NumGoroutine())
	logListener.Stop()
	fmt.Println("the logListener has been closed")
	time.Sleep(time.Second * 5)
	fmt.Println("runtime goroutine number:", runtime.NumGoroutine())
	time.Sleep(time.Hour * 1)
}
