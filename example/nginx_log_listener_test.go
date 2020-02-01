package example

import (
	"bufio"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/taoey/go-log-listener/listener"
	"os"
	"strconv"
	"testing"
	"time"
)

// Nginx 日志监听器，实现网站PV UV统计，为便于Nginx日志处理，Nginx日志配置成JSON格式
// PV统计：Redis sorted-set  : ZINCRBY key increment member
// UV统计：Redis hyperLogLog : PFADD key element

/*
Nginx 日志配置
	log_format  main  '{"addr":"$remote_addr","time":"$time_local","req":"$request"}';
	access_log  logs/access.log  main;
*/
var (
	redisClient    *redis.Pool
	REDIS_ADDR     = "127.0.0.1:6379"
	REDIS_PASSWORD = "yourpassword"
	FILE_PATH      = "access.log"
)

func init() {
	maxIdle := 500
	maxActive := 2000
	maxIdleTimeout := 20

	// 建立连接池
	redisClient = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: time.Duration(maxIdleTimeout) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", REDIS_ADDR,
				redis.DialPassword(REDIS_PASSWORD))
			if err != nil {
				return nil, err
			}
			return con, nil
		},
	}
}

type NginxLog struct {
	time string
	ip   string
}

type StorageBlock struct {
	counterType  string
	storageModel string
	nginxLog     NginxLog
}

func logHandler(logStr string) interface{} {
	fmt.Print("日志处理：", logStr)
	return NginxLog{}
}

func logStatageHandler(log interface{}) {
	//redisPool.Cmd("set",log.url)
	//fmt.Println(log)
}

// 模拟Nginx日志产生
func nginxLogCreater() {

	if fileObj, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err == nil {
		defer fileObj.Close()
		// 使用WriteString方法,写入字符串并返回写入字节数和错误信息
		writeObj := bufio.NewWriterSize(fileObj, 4096)
		for i := 0; i < 100; i++ {
			content := strconv.Itoa(i) + "\n"
			fmt.Print("日志写入", content)
			if _, err := writeObj.WriteString(content); err != nil {
				fmt.Println("write string wrong:", i, err)
			}
			writeObj.Flush()
			time.Sleep(time.Second * 1)
		}
	}
	file, err := os.Open(FILE_PATH)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}

}

func Test00(t *testing.T) {
	client := redisClient.Get()
	defer client.Close()

	//reply, err := client.Do("set", "tao", "12")
	//fmt.Println(reply,err)
	go nginxLogCreater()
	time.Sleep(time.Second)

	logListener := listener.NewDefaultLogListener(FILE_PATH)
	logListener.SetHandler(logHandler, logStatageHandler)
	logListener.Run()

	time.Sleep(time.Hour)
}
