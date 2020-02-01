# go-log-listener
> 本项目始于2019年春节期间，因新型冠状病毒爆发，在家里闲着无聊，遂进行Golang的并发编程学习，在看完慕课网的《 基于Golang协程实现流量统计系统》后有感而发，并结合《Go语言并发之道》进行Go语言的深入学习了解，并最终抽象成为一个简单易用的日志监听库。本人编程能力有限，如有设计考虑不周，还请大佬进行指点。

## 一、功能介绍

// TODO



## 二、快速开始

// TODO



## 三、实例：Nginx日志监听

// TODO 


## 四、实现原理
```
log_format  main  '{"addr":"$remote_addr","time":"$time_local","req":"$request"}';
access_log  logs/access.log  main;
```