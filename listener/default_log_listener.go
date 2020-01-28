package listener

import (
	"bufio"
	"io"
	"os"
	"time"
)

// 默认日志监听处理器
type DefaultLogListener struct {
	logChannel     chan string
	storageChannel chan interface{}
	watchChannel   chan int
	startLineNum   int //从特定的行开始监听
	//currentLineNum int //当前遍历的行，goroutine-ReadFileLineByLine进行维护

	RefreshTime int64
	FilePath    string
}

// 构造函数
func NewDefaultLogListener(filePath string, refreshTime int64) *DefaultLogListener {
	return NewDefaultLogListenerWithParams(filePath, refreshTime, 15, 15)
}

// 带参数的构造函数
func NewDefaultLogListenerWithParams(filePath string, refreshTime int64, logChannelSize, storageChannelSize int) *DefaultLogListener {
	return &DefaultLogListener{
		logChannelSize: logChannelSize,
		logChannel:     make(chan string, logChannelSize),
		storageChannel: make(chan interface{}, storageChannelSize),
		watchChannel:   make(chan int),

		RefreshTime: refreshTime,
		FilePath:    filePath,
	}
}

// 日志监听模块
func (this *DefaultLogListener) ReadFileLineByLine(filePath string, logChannel chan string) error {
	file, err := os.Open(filePath)
	if err != nil {
		LOG.Warnf("ReadFileLinebyLine can't open file:%s", filePath)
		return err
	}
	defer close(logChannel)
	defer file.Close()

	bufferRead := bufio.NewReader(file)

Loop:
	for {
		select {
		case <-this.watchChannel:
			// 读取当前行号并记录
			break Loop
		default:
			line, err := bufferRead.ReadString('\n')
			logChannel <- line
			if err != nil {
				if err == io.EOF {
					time.Sleep(time.Second * time.Duration(this.RefreshTime)) // 读取日志刷新时间
				} else {
					LOG.Warningf("ReadFileLinebyLine read log error")
				}
			}
		}
	}
	return nil
}

// 日志分析模块:仅提供日志打印输出功能，具体功能需要覆盖该方法
func (this *DefaultLogListener) logHandler(logChannel chan string, storageChannel chan interface{}) {
	// TODO 添加LogChannel的关闭状态的监听：logChannel关闭 && logChannel中没有元素时，关闭storageChannel并退出goroutine
	for logStr := range logChannel {
		if logStr != "" {
			LOG.Info("日志处理", logStr)
		}
	}
}

// 日志写入模块
func (this *DefaultLogListener) dataStorage(storageChannel chan interface{}, pool *interface{}) {
	// dosomething
}

// 启动监听
func (this *DefaultLogListener) Run() {
	go this.ReadFileLineByLine(this.FilePath, this.logChannel)
	for i := 0; i < 15; i++ {
		go this.logHandler(this.logChannel, this.storageChannel)
	}
}

// 停止监听
func (this *DefaultLogListener) Stop() {
	this.watchChannel <- 1
}

// 重新启动监听
func (this *DefaultLogListener) Restart() {
	this.Stop()
	this.logChannel = make(chan string, this.logChannelSize)
	this.Run()
}
