package example

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/taoey/go-log-listener/listener"
	"math/rand"
	"os"
	"strings"
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
	Addr string
	Time string
	Req  string
}

func logHandler(logStr string) interface{} {
	log := NginxLog{}
	json.Unmarshal([]byte(logStr), &log)
	log.Time = log.Time[0:11]
	log.Req = strings.Split(log.Req, " ")[1]
	fmt.Println("处理：", log)
	return log
}

func logStatageHandler(logObj interface{}) {
	log := (logObj).(NginxLog)
	client := redisClient.Get()
	defer client.Close()
	// PV 统计
	client.Send("ZINCRBY", fmt.Sprintf("PV:%s", log.Time), 1, log.Req)
	// UV 统计
	client.Send("PFADD", fmt.Sprintf("UV:%s", log.Time), log.Addr)
	// 使用redigo的pipeline方式进行存储
	client.Flush()
	fmt.Println("存储：", log)

}

func randInt(min int, max int) int {
	//rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + rand.Intn(max-min)
}

// 模拟Nginx日志产生
// {"addr":"127.0.0.1","time":"01/Feb/2020:16:57:03 +0800","req":"GET /favicon.ico HTTP/1.1"}
func nginxLogCreater() {

	if fileObj, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err == nil {
		defer fileObj.Close()
		// 使用WriteString方法,写入字符串并返回写入字节数和错误信息
		writeObj := bufio.NewWriterSize(fileObj, 4096)
		for i := 0; i < 100; i++ {
			addr := fmt.Sprintf("192.168.%d.%d", randInt(2, 66), randInt(2, 88))
			timeLocal := fmt.Sprintf("%02d/Feb/2020:16:57:03 +0800", randInt(1, 10))
			req := fmt.Sprintf("GET /%03d.png HTTP/1.1", randInt(1, 100))
			content := fmt.Sprintf("{\"addr\":\"%s\",\"time\":\"%s\",\"req\":\"%s\"}\n", addr, timeLocal, req)
			fmt.Print("写入:", content)

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
