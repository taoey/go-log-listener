package listener

type ILogListener interface {
	// 将日志文件读入到logChannel管道
	ReadFileLineByLine(filePath string, logChannel chan string) error
	// 日志处理
	logHandler(logChannel chan string, storageChannel chan interface{})
	// 数据存储
	dataStorage(storageChannel chan interface{}, pool *interface{})
	// 启动
	Run()
	//停止
	Stop()
}
