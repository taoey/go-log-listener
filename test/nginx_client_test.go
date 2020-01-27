package test

import (
	"github.com/taoey/go-log-listener/listener"
	"testing"
	"time"
)

// 日志监听测试
func Test01(t *testing.T) {
	filePath := "E:\\projects\\go-mod\\log-listener\\go.sum"
	logListener := listener.DefaultLogListener{
		RefreshTime: 1,
		FilePath:    filePath,
	}

	logListener.Run()

	time.Sleep(time.Hour * 1)
}

func Test02(t *testing.T) {
	filePath := "E:\\projects\\go-mod\\log-listener\\go.sum"
	logListener := listener.NewNginxLogListener(filePath, 3)

	logListener.Run()

	time.Sleep(time.Hour * 1)
}
