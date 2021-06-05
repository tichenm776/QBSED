package http

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/DeanThompson/ginpprof"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
	"zhiyuan/QBSED/configs"
	//"zhiyuan/koala_api_go/koala_api"
	"zhiyuan/QBSED/internal/model"
	"zhiyuan/QBSED/internal/service"
	//"zhiyuan/license"
)

var (
	svc *service.Service
)

// New new a gin server.
func New() {
	var (
		hc struct {
			Server *model.ServerConfig
		}
		cpath string
	)
	log4go_path :="./statistic.xml"
	log4go.LoadConfiguration(log4go_path)
	// 初始化
	if runtime.GOOS == "linux" {
		cpath = "./configs/http.toml"
		if runtime.GOARCH == "arm" {
			cpath = "./configs/http.toml"
			configs.Init("./conf.yaml")
		} else if runtime.GOARCH == "amd64" {
			cpath = "./configs/http.toml"
			configs.Init("./conf.yaml")
		}
	} else if runtime.GOOS == "windows" {
		cpath = "F:/program/code/go/go_project/src/zhiyuan/QBSED/configs/http.toml"
		configs.Init("F:/program/code/go/go_project/src/zhiyuan/QBSED/conf.yaml")
	}

	if _, err := toml.DecodeFile(cpath, &hc); err != nil {
		log4go.Error("read toml file error(%v)", err)
	}

	svc = service.New(configs.Conf)
	//svc.Koala_Action()
	//koala.Init(configs.Gconf.KoalaHost)
	//koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	//err := initAll()
	//if err != nil{
	//	log4go.Error("未激活")
	//	return
	//}
	//go SyncStatusCron()
	//go SyncStatisticCron()
	//go Statitic_nowday_()
	//go Statistic_List_Cron()
	//go GetPositionMap_Cron()
	//go DeleteRecord_Cron()
	//svc
	engine := gin.Default()
	initRouter(engine)
	gin.SetMode(gin.ReleaseMode)
	engine.Run(hc.Server.Addr)

}
var lock sync.Mutex
func  SyncStatusCron() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "*/10 * * * * ?"
	cronTarget.AddFunc(spec, func() {
		lock.Lock()
		svc.PrepareForRecord(1)
		lock.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}
var flag1 = 0
func  SyncStatisticCron() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "*/24 * * * * ?"
	cronTarget.AddFunc(spec, func() {
		//lock.Lock()
		if Position_map[-99] == nil{
			GetPositionMap()
		}
		Position_map_copy := Position_map
		//svc.PrepareForstatistic("",Position_map_copy)
		//fmt.Println("function begin ---------------------------------------------------")
		log4go.Info("function begin ---------------------------------------------------",time.Now().Unix())
		fmt.Println("function end ---------------------------------------------------",time.Now().Unix())
		//lock3.Lock()
		svc.PrepareForstatistic2("",Position_map_copy)
		//lock3.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}

func  Statitic_nowday_() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "*/15 * * * * ?"
	cronTarget.AddFunc(spec, func() {
		//lock2.Lock()
		Statitic_nowday_W()
		//lock2.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}
func  Statistic_List_Cron() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "55 59 23 * * ?"
	//spec := "00 37 14 * * ?"
	cronTarget.AddFunc(spec, func() {
		//lock.Lock()
		nowtime := time.Now()
		nowdate := nowtime.String()[0:10]
		SaveFile(nowdate)
		//lock.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}
func  GetPositionMap_Cron() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "00 */30 * * * ?"
	//spec := "00 37 14 * * ?"
	cronTarget.AddFunc(spec, func() {
		//lock.Lock()
		GetPositionMap()
		//lock.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}
func  DeleteRecord_Cron() {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "00 01 * * * ?"
	//spec := "25 10 * * * ?"
	//spec := "00 37 14 * * ?"
	cronTarget.AddFunc(spec, func() {
		//lock.Lock()
		svc.DeleteRecord(30)
		//svc.DeleteRecord(30)
		//lock.Unlock()
	})
	cronTarget.Start()
	log4go.Info("定时开始-----！")
}

func initRouter(e *gin.Engine) {

	e.Use(Cors())
	ginpprof.Wrap(e)
	e.POST("/test", Test23)
	authority := e.Group("/v1")
	{
		authority.GET("/start", howToStart)
		authority.POST("/login", LoginIn)
		authority.POST("/logout", Logout)
		authority.POST("/register", Register)
		authority.GET("/signurl", GetSignUrl)

	}

	bedroom := e.Group("/v1")
	{
		bedroom.GET("/test1", Test)
		bedroom.GET("/test", Test2)
		bedroom.GET("/syncstatistic", Test3)
		bedroom.GET("/test23", EmployeerecordsGroup)
		bedroom.GET("/test4", Test4)
		bedroom.GET("/test5", Test5)
		bedroom.GET("/G_MAP", Test6)
		bedroom.GET("/G_MAP2", Test6)
	}


	statistic_group := e.Group("/v1")
	{
		//statistic_group.POST("/statistics",Statistic2)
		statistic_group.GET("/statistics",Statistic2)
		statistic_group.GET("/statistic/today",Statistic_nowday)
		statistic_group.GET("/statistic/list",Statistic_list)
		statistic_group.GET("/statistic/records",GetOneStatisticRecords)
		statistic_group.GET("/statistic/record",GetStatisticRecord)
		statistic_group.GET("/time",Time)
	}

	subject := e.Group("/v1")
	{
		//subject.GET("/koala_departments", GetDepartment)
		//subject.GET("/subjects", GetSubjectByCategory)
		subject.GET("/event/constants", Eventconstants)
		//subject.POST("/employee/records", Employeerecords)
		//subject.POST("/employee/records_count", Employeerecords_days)
		subject.POST("/employee_group/records", EmployeerecordsGroup)
		subject.POST("/stranger/records", Strangerrecords)


		subject.POST("/statistic", Statistic_Project_Create)
		subject.DELETE("/statistic/:id", Statistic_Project_Delete)
		subject.PUT("/statistic/:id", Statistic_Project_Update)
		subject.GET("/statistic", Statistic_Project_Get)
		subject.GET("/accesslist", GetAccessList)
		subject.GET("/persongrouplist", GetPersonGroupList)


	}


	classin := e.Group("")
	{
		classin.GET("/sync/student", Test8)
	}

}

// example for http request handler.
func howToStart(c *gin.Context) {
	c.String(0, "Golang 大法好 !!!")
}

func initAll()error{
	//licensePath := ""
	//if runtime.GOOS == "linux" {
	//	if runtime.GOARCH == "arm" {
	//		licensePath = "/home/zybox/face_server/lic.txt"
	//	} else if runtime.GOARCH == "amd64" {
	//		licensePath = "/home/zhiyuan/ai_dormitory_apis/lic.txt"
	//	}
	//} else if runtime.GOOS == "windows" {
		log4go.LoadConfiguration("./log4go.xml")
		configs.Init("./conf.yaml")
	//}
	//if runtime.GOOS != "windows" {
	//	if !license.IsValidLicOfGeneral(licensePath) {
	//		return errors.New("-99")
	//	}
	//}
	return nil
}
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method      //请求方法
		origin := c.Request.Header.Get("Origin")        //请求头部
		var headerKeys []string                             // 声明请求头keys
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Origin", "*")        // 这是允许访问所有域
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")      //服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			//  header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//              允许跨域设置                                                                                                      可以返回其他子段
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")      // 跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Max-Age", "172800")        // 缓存请求信息 单位为秒
			c.Header("Access-Control-Allow-Credentials", "false")       //  跨域请求是否需要带cookie信息 默认设置为true
			c.Set("content-type", "application/json")       // 设置返回格式是json
		}

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 处理请求
		c.Next()        //  处理请求
	}
}





