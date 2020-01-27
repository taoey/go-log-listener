package test

import (
	"fmt"
	"github.com/taoey/go-log-listener/listener"
)

// 自定义Nginx实现类，能够解析Nginx日志
// 需要重写三个接口

type NginxLogListener struct {
	listener.DefaultLogListener
}

// 构造函数
func NewNginxLogListener(filePath string, refreshTime int64) *NginxLogListener {
	return &NginxLogListener{
		listener.DefaultLogListener{
			FilePath:    filePath,
			RefreshTime: refreshTime,
		},
	}
}

// 重写日志分析模块:仅提供日志打印输出功能，具体功能需要覆盖该方法
func (this *NginxLogListener) logHandler(logChannel chan string, storageChannel chan interface{}) {
	for logStr := range logChannel {
		if logStr != "" {
			fmt.Print("日志处理：", logStr)
		}
	}
}

// 重写启动模块，将logHandler变更为NginxLogListener中的logHandler
func (this *NginxLogListener) Run() {
	var logChannel = make(chan string, 15)
	var storageChannel = make(chan interface{}, 15)

	go this.ReadFileLineByLine(this.FilePath, logChannel)

	go this.logHandler(logChannel, storageChannel)
}
