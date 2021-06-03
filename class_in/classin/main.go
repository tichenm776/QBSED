package main

import (
	"zhiyuan/classin/internal/server/http"
)
//var log = logrus.New()
//var isConnected bool
func main() {


	// ==== SET LOGGING
	//Formatter := new(log.JSONFormatter)
	//log.SetFormatter(Formatter)
	//InitLog()
	http.New()
}

//func InitLog() {
//	path := "/home/zybox/logs/hotel.log"
//	//path := "./hotel.log"
//	/* 日志轮转相关函数
//	 `WithLinkName` 为最新的日志建立软连接
//	 `WithRotationTime` 设置日志分割的时间，隔多久分割一次
//	 WithMaxAge 和 WithRotationCount二者只能设置一个
//	  `WithMaxAge` 设置文件清理前的最长保存时间
//	  `WithRotationCount` 设置文件清理前最多保存的个数
//	 */
//	// 下面配置日志每隔 1 分钟轮转一个新文件，保留最近 3 分钟的日志文件，多余的自动清理掉。
//	writer, _ := rotatelogs.New(
//		path+".%Y%m%d%H%M",
//		rotatelogs.WithLinkName(path),
//		rotatelogs.WithMaxAge(time.Duration(6000*72)*time.Second),
//		rotatelogs.WithRotationTime(time.Duration(6000*12)*time.Second),
//	)
//	log.SetOutput(writer)
//	//log.SetFormatter(&log.JSONFormatter{})
//}