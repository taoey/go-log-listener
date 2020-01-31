package listener

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

// 默认日志监听处理器
type DefaultLogListener struct {
	once           sync.Once
	logChannel     chan string
	logChannelSize int
	storageChannel chan interface{}
	watchChannel   chan int
	startLineNum   int //从特定的行开始监听
	//currentLineNum int //当前遍历的行，goroutine-ReadFileLineByLine进行维护

	logHandlerNum        int
	logStarageHandlerNum int
	varLogHandler        func(string) interface{} // 日志处理回调
	varStorageHandler    func(interface{})        // 存储回调

	RefreshTime int64
	FilePath    string
}

// 构造函数
func NewDefaultLogListener(filePath string, refreshTime int64) *DefaultLogListener {
	return NewDefaultLogListenerWithParams(filePath, refreshTime, 15, 15, 15, 1)
}

// 带参数的构造函数
func NewDefaultLogListenerWithParams(filePath string, refreshTime int64,
	logChannelSize int, logHandlerNum int, storageChannelSize int, logStarageHandlerNum int) *DefaultLogListener {
	return &DefaultLogListener{
		logChannelSize:       logChannelSize,
		logChannel:           make(chan string, logChannelSize),
		storageChannel:       make(chan interface{}, storageChannelSize),
		watchChannel:         make(chan int),
		logHandlerNum:        logHandlerNum,
		logStarageHandlerNum: logStarageHandlerNum,

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
func (this *DefaultLogListener) logHandler(
	logChannel chan string, storageChannel chan interface{}, f func(string) interface{}) {
	// 使用Once 防止 多个goroutine关闭同一个channel
	defer func() {
		this.once.Do(func() {
			LOG.Info("handle data finished!")
			close(storageChannel)
		})
	}()
	// logChannel 关闭后待logChannel缓冲区为空时，自动跳出该方法
	for logStr := range logChannel {
		if logStr != "" {
			logObject := f(logStr)
			storageChannel <- logObject
		}
	}
}

// 日志写入模块
func (this *DefaultLogListener) dataStorage(storageChannel chan interface{}, f func(interface{})) {
	// logChannel 关闭后待logChannel缓冲区为空时，自动跳出该方法
	for logObj := range storageChannel {
		if logObj != nil {
			f(logObj)
		}
	}
	LOG.Info("storage data finished!")
}

// 启动监听
func (this *DefaultLogListener) Run() {
	go this.ReadFileLineByLine(this.FilePath, this.logChannel)
	for i := 0; i < this.logHandlerNum; i++ {
		go this.logHandler(this.logChannel, this.storageChannel, this.varLogHandler)
	}
	for i := 0; i < this.logStarageHandlerNum; i++ {
		go this.dataStorage(this.storageChannel, this.varStorageHandler)
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

// 设置处理方法
func (this *DefaultLogListener) SetHandler(logHandler func(string) interface{}, storageHandler func(interface{})) {
	this.varLogHandler = logHandler
	this.varStorageHandler = storageHandler
}
