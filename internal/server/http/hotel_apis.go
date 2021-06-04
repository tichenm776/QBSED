package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/common/log"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
	//"zhiyuan/koala_api_go/koala_api"

	//"zhiyuan/koala_api_go/koala_api"
	"zhiyuan/QBSED/configs"
	"zhiyuan/QBSED/internal/koala"
	"zhiyuan/QBSED/internal/model"
	"zhiyuan/zyutil"
)

func LoginIn(c *gin.Context) {
	code := 0
	err_msg := ""

	// Read the Body content
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}

	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	koala.LogRequest(c)

	username := c.PostForm("username")
	if username == "" {
		code = ERR_INPUT_NULL
		err_msg = "缺少用户名"
	}
	password := c.PostForm("password")
	if password == "" {
		code = ERR_INPUT_NULL
		err_msg = "缺少密码"
	}
	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
		})
		return
	}

	level := 0
	if strings.Contains(username, "@") {
		configs.Init("./conf.yaml")
		koala.Init(configs.Gconf.KoalaHost)
		_,err := koala.KoalaLogin(username, password)
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    ERR_KOALA_LOGIN,
				"err_msg": "账号密码错误!",
			})
			return
		}
		level = 1
		//configs.Init("./conf.yaml")
		_, err1 := configs.Gconf.EditYaml(username, password)

		if err1 != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    ERR_WRITE_YAML,
				"err_msg": err1.Error(),
			})
			return
		}
	}

	if level == 1 {
		data := make(map[string]interface{})
		data["level"] = 1
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "username",
			Value:    username,
			MaxAge:   0,
			Path:     "/",
			Domain:   "",
			Secure:   false,
			HttpOnly: true,
		})
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
			"data":    data,
		})
		return
	} else {
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		koala.LoginZybox(level, c)
	}

}

func Logout(c *gin.Context){

	//ip,_:=LocalIPv4s()

	c.SetCookie("username", "admin", -1, "/", "", false, true)
	//_logout := "http://"+ip[0]

	//c.Redirect(http.StatusFound, _logout)
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"err_msg": "",
	})
	return
}


func Register(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
		User model.User_json
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.Bind(&User)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
	}
	if User.Password != User.Password_check{
		resp4Device.Code = -100
		resp4Device.Err_msg = "两次密码不一致"
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	//注册

	svc.UpdateCourse()

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
	})
	return


}

func Employeerecords_days(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
		EmployeeRecords model.EmployeeRecords_json
		//data           interface{}
	)
	//subject_id_string := c.Query("subject_id")
	//if subject_id_string != ""{
	//	subject_id,_ := strconv.Atoi(subject_id_string)
	//}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&EmployeeRecords)

	if EmployeeRecords.Page == 0 {
		EmployeeRecords.Page = 1
	}
	if EmployeeRecords.Size == 0 {
		EmployeeRecords.Size = 10
	}
	if err != nil {
		//resp4Device.Code = -100
		//EmployeeRecords.Snap_position = ""
		EmployeeRecords.Subject_id = 0
		EmployeeRecords.Screen_id= 0
		EmployeeRecords.User_role = -1
		EmployeeRecords.Name= ""
		EmployeeRecords.Snap_begin_time= ""
		EmployeeRecords.Snap_end_time= ""
	}
	EmployeeRecords.Snap_begin_time = svc.UTCchange(EmployeeRecords.Snap_begin_time)
	EmployeeRecords.Snap_end_time = svc.UTCchange(EmployeeRecords.Snap_end_time)
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map

	if (EmployeeRecords.Snap_begin_time != "" && EmployeeRecords.Snap_end_time != "") && EmployeeRecords.Snap_begin_time[0:10] == EmployeeRecords.Snap_end_time[0:10]{
		res, err, count, total := svc.GetIdentificationRecords(EmployeeRecords,Position_map_copy)
		if err != nil {
			resp4Device.Code = -100
			resp4Device.Err_msg = "查询记录失败"
			resp4Device.Data = res
		}
		if resp4Device.Code != 0 {
			zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
			return
		}
		resp4Device.Data = res
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data":    resp4Device.Data,
			"count":   count,
			"current": EmployeeRecords.Page,
			"size":    EmployeeRecords.Size,
			"total":   total,
		})
	}else{
		datelist := svc.GetBetweenDates(EmployeeRecords.Snap_begin_time,EmployeeRecords.Snap_end_time)
		time1 := EmployeeRecords.Snap_begin_time[11:]
		time2 := EmployeeRecords.Snap_end_time[11:]
		temptotal:=0
		tempcount:=0
		tempdata := make([]model.IdentificationRecord,0)
		for _,v := range datelist{
			EmployeeRecords.Snap_begin_time = v+" "+time1
			EmployeeRecords.Snap_end_time = v+" "+time2
			res, err, count, total := svc.GetIdentificationRecords(EmployeeRecords,Position_map_copy)
			if err != nil{
				resp4Device.Code = -100
				resp4Device.Err_msg = "查询记录失败"
				resp4Device.Data = res
			}
			if resp4Device.Code != 0 {
				zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
				return
			}
			temptotal += total
			tempcount += count
			tempdata = append(tempdata,res... )
		}
		resp4Device.Data = tempdata
		total := zyutil.GetTotal(tempcount, EmployeeRecords.Size)
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data":    resp4Device.Data,
			"count":   tempcount,
			"current": EmployeeRecords.Page,
			"size":    EmployeeRecords.Size,
			"total":   total,
		})
	}




	//res, err, count, total := svc.GetIdentificationRecords(EmployeeRecords,Position_map_copy)
	//
	//if err != nil {
	//	resp4Device.Code = -100
	//	resp4Device.Err_msg = "查询记录失败"
	//	resp4Device.Data = res
	//}
	//if resp4Device.Code != 0 {
	//	zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
	//	return
	//}
	//
	//resp4Device.Data = res
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"code":    resp4Device.Code,
	//	"err_msg": resp4Device.Err_msg,
	//	"data":    resp4Device.Data,
	//	"count":   count,
	//	"current": EmployeeRecords.Page,
	//	"size":    EmployeeRecords.Size,
	//	"total":   total,
	//})
}

func EmployeerecordsGroup(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
		EmployeeRecords model.EmployeeRecords_json
		//data           interface{}
	)
	//subject_id_string := c.Query("subject_id")
	//if subject_id_string != ""{
	//	subject_id,_ := strconv.Atoi(subject_id_string)
	//}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&EmployeeRecords)

	if EmployeeRecords.Page == 0 {
		EmployeeRecords.Page = 1
	}
	if EmployeeRecords.Size == 0 {
		EmployeeRecords.Size = 10
	}
	if err != nil {
		//resp4Device.Code = -100
		//EmployeeRecords.Snap_position = ""
		EmployeeRecords.Subject_id = 0
		EmployeeRecords.Screen_id= 0
		EmployeeRecords.User_role = -1
		EmployeeRecords.Name= ""
		EmployeeRecords.Snap_begin_time= ""
		EmployeeRecords.Snap_end_time= ""
	}
	EmployeeRecords.Snap_begin_time = svc.UTCchange(EmployeeRecords.Snap_begin_time)
	EmployeeRecords.Snap_end_time = svc.UTCchange(EmployeeRecords.Snap_end_time)
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map

	//str := ""
	matched, err := regexp.MatchString(`^[1-9]\d{3}-(0[1-9]|1[0-2])-(0[1-9]|[1-2][0-9]|3[0-1])\s+(20|21|22|23|[0-1]\d):[0-5]\d:[0-5]\d$`, EmployeeRecords.Snap_begin_time)
	//fmt.Println(matched, err)
	if matched == false{
		temp1 := time.Now().String()[0:10] + EmployeeRecords.Snap_begin_time
		temp2 := time.Now().String()[0:10] + EmployeeRecords.Snap_end_time
		EmployeeRecords.Snap_begin_time = temp1
		EmployeeRecords.Snap_end_time = temp2
	}

	if (EmployeeRecords.Snap_begin_time != "" && EmployeeRecords.Snap_end_time != "") && EmployeeRecords.Snap_begin_time[0:10] == EmployeeRecords.Snap_end_time[0:10]{
		res, err, count, total := svc.GetIdentificationRecords_GroupBy(EmployeeRecords,Position_map_copy)
		if err != nil {
			resp4Device.Code = -100
			resp4Device.Err_msg = "查询记录失败"
			resp4Device.Data = res
		}
		if resp4Device.Code != 0 {
			zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
			return
		}
		resp4Device.Data = res
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data":    resp4Device.Data,
			"count":   count,
			"current": EmployeeRecords.Page,
			"size":    EmployeeRecords.Size,
			"total":   total,
		})
	}else{
		datelist := svc.GetBetweenDates(EmployeeRecords.Snap_begin_time,EmployeeRecords.Snap_end_time)
		time1 := EmployeeRecords.Snap_begin_time[11:]
		time2 := EmployeeRecords.Snap_end_time[11:]
		temptotal:=0
		tempcount:=0
		tempdata := make([]model.IdentificationRecord,0)
		for _,v := range datelist{
			EmployeeRecords.Snap_begin_time = v+" "+time1
			EmployeeRecords.Snap_end_time = v+" "+time2
			EmployeeRecords.Size = 10000
			EmployeeRecords.Page = 1
			res, err, count, total := svc.GetIdentificationRecords_GroupBy(EmployeeRecords,Position_map_copy)
			if err != nil{
				resp4Device.Code = -100
				resp4Device.Err_msg = "查询记录失败"
				resp4Device.Data = res
			}
			if resp4Device.Code != 0 {
				zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
				return
			}
			temptotal += total
			tempcount += count
			tempdata = append(tempdata,res... )
		}

		resp4Device.Data = tempdata
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data":    resp4Device.Data,
			"count":   len(tempdata),
			"current": EmployeeRecords.Page,
			"size":    EmployeeRecords.Size,
			"total":   1,
		})
	}

}

func Strangerrecords(c *gin.Context) {
	var (
	resp4Device    model.Resp4Device
	StrangerRecords model.StrangerRecords_json
	//data           interface{}
)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&StrangerRecords)

	if StrangerRecords.Page == 0 {
		StrangerRecords.Page = 1
	}
	if StrangerRecords.Size == 0 {
		StrangerRecords.Size = 10
	}
	if err != nil {
		//resp4Device.Code = -100
		//StrangerRecords.Snap_position = ""
		StrangerRecords.Screen_id= 0
		StrangerRecords.Snap_begin_time= ""
		StrangerRecords.Snap_end_time= ""
	}
	StrangerRecords.Snap_begin_time = svc.UTCchange(StrangerRecords.Snap_begin_time)
	StrangerRecords.Snap_end_time = svc.UTCchange(StrangerRecords.Snap_end_time)
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map
	res, err, count, total := svc.GetStrangerRecords(StrangerRecords,Position_map_copy)

	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询记录失败"
		resp4Device.Data = res
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	resp4Device.Data = res

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
		"count":   count,
		"current": StrangerRecords.Page,
		"size":    StrangerRecords.Size,
		"total":   total,
	})
	}

func Eventconstants(c *gin.Context) {	var (
	resp4Device    model.Resp4Device
	//data           interface{}
)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	koala.Init(configs.Gconf.KoalaHost)
	koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	datavalue , err := koala.Constants()
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "koala获取参数错误"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	screens := datavalue.Get("screens")
	subject_type_options := datavalue.Get("select_opt_user_role")
	resp4Device.Data = map[string]interface{}{
		"screens":screens,
		"subject_type_options":subject_type_options,
	}
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "koala获取参数错误"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})}

func GetDepartment(c *gin.Context) {
	code := 0
	err_msg := ""

	//LogRequest(c)

	//koala.KoalaLogin(configs.Gconf.KoalaUsername, configs.Gconf.KoalaPassword)
	data, err := koala.GetDepartment()

	if err != nil {
		zyutil.KoalaReturn(c, ERR_KOALA_DATA, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"err_msg": err_msg,
		"data":    data,
	})
}

func GetSubjectByCategory(c *gin.Context) {
	code := 0
	err_msg := ""

	//LogRequest(c)

	category := c.Query("category")
	if category == "" {
		category = "employee"
	}

	name := c.Query("name")
	name = strings.Replace(name, " ", "+", -1)
	department := c.Query("department")
	tempPage := c.Query("page")
	tempSize := c.Query("size")
	page, err := strconv.Atoi(tempPage)
	if err != nil {
		page = 1
	}
	tempSize = "50000"
	size, err := strconv.Atoi(tempSize)
	if err != nil {
		size = 5000
	}

	if code != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    code,
			"err_msg": err_msg,
		})
		return
	}
	koala.Init(configs.Gconf.KoalaHost)
	koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	datas, data, err := koala.GetSubjectsByCondition(category, name, department, page, size, configs.Gconf.KoalaHost)
	if err != nil {
		zyutil.KoalaReturn(c, ERR_KOALA_DATA, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    code,
		"err_msg": err_msg,
		"data":    data,
		"count":   datas["count"],
		"current": datas["current"],
		"size":    datas["size"],
		"total":   datas["total"],
	})
}

var Position_map = map[int]interface{}{}
var Persongrouplist_g = map[int]interface{}{}
var Accesslist_g = map[int]interface{}{}
func GetPositionMap()(map[int]interface{},error){
	//koala.Init(configs.Gconf.KoalaHost)
	//koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	koala.Init(configs.Gconf.KoalaHost)
	koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	datavalue , err := koala.Constants()
	if err != nil{
		log4go.Error(err)
		return nil, err
	}
	screens := datavalue.Get("screens")
	screenslist,err := screens.Array()
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	for _,v := range screenslist{
		key ,_:= v.(map[string]interface{})["value"].(json.Number).Int()
		Position_map[key]=v.(map[string]interface{})["label"].(string)
	}
	Position_map[-99]="activate"
	PersonGroupList_subject ,err := koala.GetPersonGroupList("0")
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	PersonGroupList_visitor ,err := koala.GetPersonGroupList("1")
	if err != nil {
		log4go.Error(err)
		return nil, err
	}
	if len(PersonGroupList_visitor) != 0{
		PersonGroupList_subject = append(PersonGroupList_subject,PersonGroupList_visitor...)
	}
	for _,v := range PersonGroupList_subject{
		key ,_:= v["id"].(json.Number).Int()
		Persongrouplist_g[key]=v["name"]
	}
	AccessList,_ := koala.GetAccessList()
	for _,v := range AccessList{
		key ,_:= v["id"].(json.Number).Int()
		Accesslist_g[key]=v["name"]
	}
	Accesslist_g[-99]="activate"
	Persongrouplist_g[-99]="activate"
	return Position_map,nil
}



func Statistic(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		Statistic_json model.Statistic_json
		datelist []string
		datetime string
	)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	datetime = nowdate
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&Statistic_json)
	if err != nil {
		Statistic_json.Snap_begin_time= datetime
		Statistic_json.Snap_end_time= datetime
	}
	if Statistic_json.Snap_begin_time == "" || Statistic_json.Snap_end_time == ""{
		Statistic_json.Snap_begin_time= datetime
		Statistic_json.Snap_end_time= datetime
	}
	Statistic_json.Snap_begin_time = svc.UTCchange(Statistic_json.Snap_begin_time)
	Statistic_json.Snap_end_time = svc.UTCchange(Statistic_json.Snap_end_time)
	if Statistic_json.Snap_begin_time[0:10] == Statistic_json.Snap_end_time[0:10]{
		res , err := svc.GetStatistic(Statistic_json.Snap_begin_time[0:10])
		if err != nil{
			//resp4Device.Code = -404
			//resp4Device.Err_msg = "无历史记录"
			//zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
			resp4Device.Code = 0
			resp4Device.Err_msg = "无历史记录"
			resp4Device.Data = model.Statistic{
				Desayuno:0,
				Almuerzo:0,
				Jantar:0,
			}
			c.JSON(http.StatusOK, gin.H{
				"code":    resp4Device.Code,
				"err_msg": resp4Device.Err_msg,
				"data": resp4Device.Data,
			})
			return
		}
		resp4Device.Data = res
	}else{
		datelist = svc.GetBetweenDates(Statistic_json.Snap_begin_time,Statistic_json.Snap_end_time)
		temp := model.Statistic{}
		for _,v := range datelist{
			res , err := svc.GetStatistic(v)
			if err != nil{
				//resp4Device.Code = -100
				//resp4Device.Err_msg = "无此日期数据"
				log4go.Error("无此日期数据")
				//zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
				continue
			}
			temp.Desayuno += res.Desayuno
			temp.Almuerzo += res.Almuerzo
			temp.Jantar += res.Jantar
		}
		resp4Device.Data = temp
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})
}
func Statistic_nowday(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device

	)
	resp4Device.Data = model.Resp
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})
}
func Statistic_list(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		date_time		string
		//data_list = make([]map[string]interface{},0)
	)
	date_time = c.Query("date_time")
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	if date_time == ""{
		date_time = nowdate
	}

	data_list,err := Statistic_List(date_time)
	if err != nil{
		zyutil.DeviceErrorReturn(c,-100,"接口错误:"+err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		//"data3": data_list,
		"data": data_list,
		"total":len(data_list),
		//"data2": model.G_map_time.Getmap(),
	})
}

func SaveFile(date_time string)(){
	data_list,err := Statistic_List(date_time)
	if err != nil{
		log4go.Error("获取数据出错:",err)
		return
	}
	err = model.SaveFile(date_time,data_list)
	if err != nil{
		log4go.Error("定时数据存储失败:",err)
	}
}

func ReadFile(date_time string)([]map[string]interface{},error){
	datalist,err := model.ReadFile2(date_time)
	if err != nil{
		log4go.Error("定时数据存储失败:",err)
		return []map[string]interface{}{},err
	}
	return datalist, nil
}


func Statistic_List(date_time string)([]map[string]interface{},error){
	data_list := make([]map[string]interface{},0)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	loc, _ := time.LoadLocation("Local")
	//old_datetime, _ := time.ParseInLocation("2006-01-02", date_time, loc)
	//now_datetime, _ := time.ParseInLocation("2006-01-02", nowdate, loc)
	//fmt.Println("--------------------------------------")
	//fmt.Println(nowdate)
	//fmt.Println(now_datetime)
	//fmt.Println(date_time)
	//fmt.Println(old_datetime)
	//fmt.Println("--------------------------------------")
	if date_time != nowdate{
		fmt.Println("------------------------------------------1")
		data_read,err := ReadFile(date_time)
		if err != nil{
			return data_read,err
		}
		return data_read, nil
		}
	if date_time == nowdate{
		fmt.Println("------------------------------------------2")
		data := model.G_map.Getmap()
		//timemap := model.G_map_time.Getmap()
		for k,_ := range data{
			times := model.G_map_time.Get(k)
			if value,ok := times.(map[string]interface{});ok{
				snap_time_in := value["snap_time_in"].(string)
				snap_time_out := value["snap_time_out"].(string)
				data[k].(map[string]interface{})["enter_time"] = snap_time_in
				data[k].(map[string]interface{})["leave_time"] = snap_time_out
				data[k].(map[string]interface{})["remain_12"] = false
				data[k].(map[string]interface{})["remain_24"] = false
				if snap_time_in == "" || snap_time_out == ""{
					data_list = append(data_list, data[k].(map[string]interface{}))
					continue
				}
				snap_time_in_time, _ := time.ParseInLocation("2006-01-02 15:04:05", snap_time_in, loc)
				snap_time_out_time, _ := time.ParseInLocation("2006-01-02 15:04:05", snap_time_out, loc)
				sub := snap_time_out_time.Sub(snap_time_in_time)
				fmt.Println("------------------------------------subtime is")
				fmt.Println(sub.Hours())
				if sub.Hours() >= 12.0{
					data[k].(map[string]interface{})["remain_12"] = true
				}
				if sub.Hours() >= 24.0{
					data[k].(map[string]interface{})["remain_24"] = true
				}
			}
			data_list = append(data_list, data[k].(map[string]interface{}))
		}

	}
	return data_list,nil
}


func Statitic_nowday_W()(){
	if flag1 == 1{
		return
	}
	flag1 = 1
	defer func() {
		flag1 = 0
	}()
	var (
		resp4Device    model.Resp4Device
		datetime string
		data = make(map[string]interface{},0)
		resp = make(map[string]interface{},0)
		snap_begin_time,snap_end_time string
	)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	datetime = nowdate
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	snap_begin_time = datetime
	snap_end_time = datetime
	if Position_map[-99] == nil{
		GetPositionMap()//获取Koala相机map{screen_id:name} 防止改名字与koala内名字不匹配
	}
	Position_map_copy := Position_map
	snap_begin_time = svc.UTCchange(snap_begin_time)
	snap_end_time = svc.UTCchange(snap_end_time)
	//if snap_begin_time[0:10] == snap_end_time[0:10]{
	if snap_begin_time[0:10] != nowdate{
		svc.PrepareForstatistic2(snap_begin_time[0:10],nil)
	}
	//这里
	item_map,gdata,err:=svc.Statistic_main_item_map()
	if err!=nil{
		return
	}
	res , err := svc.GetStatistic2(snap_begin_time[0:10],item_map,gdata)
	if err != nil{
		return
	}
	//到这 都是数据库内找出统计项和子项
	if len(res) == 0{
		return
	}
	data = res[0]//默认第一项
	main_id := data["id"].(int)
	items := data["items"].([]map[string]interface{})
	for k := range items{
		name := items[k]["name"].(string)
		item_id := items[k]["id"].(int)
		flag := 0
		//进出比对 1是进标签 2是出标签
		if strings.Contains(name,"进") {
			number := items[k]["number"].(int)
			resp["num_of_in"] = number
			flag = 1
		}
		if strings.Contains(name,"出"){
			number := items[k]["number"].(int)
			resp["num_of_out"] = number
			flag = 2
		}
		//获取属于该统计项的识别记录
		res ,count,total, err := svc.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,1,100000)
		if err != nil{
			resp4Device.Code = -100
			resp4Device.Err_msg = "获取识别列表失败"
			return
		}
		fmt.Println(count)
		fmt.Println(total)
		log4go.Info("res",res)
		log4go.Info("count",count)
		log4go.Info("total",total)
		loc, _ := time.LoadLocation("Local")
		//开始归类
		if flag == 1{
			for k := range res{
				//screen_id := res[k].Screen_id
				subject_id := res[k].Subject_id
				name := res[k].Name
				snaptime := res[k].Snap_time
				come_from := res[k].Come_from
				remark := res[k].Remark
				start_time := res[k].Start_time
				end_time := res[k].End_time
				snap_position := res[k].Snap_position
				data := map[string]interface{}{
					"subject_id":subject_id,
					"snap_position":snap_position,
					"name":name,
					"snap_time":snaptime,
					"come_from":come_from,
					"remark":remark,
					"start_time":start_time,
					"end_time":end_time,
					"mark":1,
				}
				snap_timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", snaptime, loc)
				str_sub_id := strconv.Itoa(subject_id)
				haskey :=model.G_map.Get(str_sub_id)
				//如果map里有了只更新时间
				if haskey != nil{
					old_snap_time := haskey.(map[string]interface{})["snap_time"].(string)
					old_snap_timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", old_snap_time, loc)
					if old_snap_timestamp.Unix() < snap_timestamp.Unix(){
						model.G_map.Set(str_sub_id,data)
						time_diary := model.G_map_time.Get(str_sub_id)
						if time_diary != nil{
							data2 := map[string]interface{}{
								"snap_time_in":snaptime,
								"snap_time_out":time_diary.(map[string]interface{})["snap_time_out"],
							}
							model.G_map_time.Set(str_sub_id,data2)
						}
					}
				}else{
					//没有就记录
					model.G_map.Set(str_sub_id,data)
					data2 := map[string]interface{}{
						"snap_time_in":snaptime,
						"snap_time_out":"",
					}
					//
					//res,err := svc.GetRecordById(screen_id,subject_id)
					//if err == nil{
					//	data2 = map[string]interface{}{
					//		"snap_time_out":res.Snap_time,
					//	}
					//}
					model.G_map_time.Set(str_sub_id,data2)
				}
			}
		}
		if flag == 2{
			for k := range res{
				subject_id := res[k].Subject_id
				name := res[k].Name
				snaptime := res[k].Snap_time
				come_from := res[k].Come_from
				remark := res[k].Remark
				start_time := res[k].Start_time
				end_time := res[k].End_time
				snap_position := res[k].Snap_position
				data := map[string]interface{}{
					"subject_id":subject_id,
					"snap_position":snap_position,
					"name":name,
					"snap_time":snaptime,
					"come_from":come_from,
					"remark":remark,
					"start_time":start_time,
					"end_time":end_time,
					"mark":2,
				}
				str_sub_id := strconv.Itoa(subject_id)
				snap_timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", snaptime, loc)
				haskey :=model.G_map.Get(str_sub_id)
				if haskey != nil{
					old_snap_time := haskey.(map[string]interface{})["snap_time"].(string)
					old_snap_timestamp, _ := time.ParseInLocation("2006-01-02 15:04:05", old_snap_time, loc)
					if old_snap_timestamp.Unix() < snap_timestamp.Unix(){
						model.G_map.Set(str_sub_id,data)
						time_diary := model.G_map_time.Get(str_sub_id)
						if time_diary != nil{
							data2 := map[string]interface{}{
								"snap_time_out":snaptime,
								"snap_time_in":time_diary.(map[string]interface{})["snap_time_in"],
							}
							model.G_map_time.Set(str_sub_id,data2)
						}
					}
				}else{
					model.G_map.Set(str_sub_id,data)
					data2 := map[string]interface{}{
						"snap_time_out":snaptime,
						"snap_time_in":"",
					}
					model.G_map_time.Set(str_sub_id,data2)
				}
			}
		}
	}
	counter := 0
	for _,v := range model.G_map.Data{
		if v.(map[string]interface{})["mark"].(int) == 1{
			counter += 1
		}
	}
	log4go.Info("G_map_data is",model.G_map.Getmap())
	resp["num_of_inside"] = counter
	model.Resp = resp
}



func Statistic2(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
		datetime string
	)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	datetime = nowdate
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	//err := c.BindJSON(&Statistic_json)

	snap_begin_time := c.Query("snap_begin_time")
	snap_end_time := c.Query("snap_end_time")
	//if err != nil {
	//	Statistic_json.Snap_begin_time= datetime
	//	Statistic_json.Snap_end_time= datetime
	//}
	if snap_begin_time == ""{
		snap_begin_time = datetime
	}
	if snap_end_time == ""{
		snap_end_time = datetime
	}
	//if Statistic_json.Snap_begin_time == "" || Statistic_json.Snap_end_time == ""{
	//	Statistic_json.Snap_begin_time= datetime
	//	Statistic_json.Snap_end_time= datetime
	//}
	snap_begin_time = svc.UTCchange(snap_begin_time)
	snap_end_time = svc.UTCchange(snap_end_time)
	if snap_begin_time[0:10] == snap_end_time[0:10]{
		if snap_begin_time[0:10] != nowdate{
			svc.PrepareForstatistic2(snap_begin_time[0:10],nil)
		}
		item_map,gdata,err:=svc.Statistic_main_item_map()
		if err!=nil{

		}
		res , err := svc.GetStatistic2(snap_begin_time[0:10],item_map,gdata)
		if err != nil{
			resp4Device.Code = 0
			resp4Device.Err_msg = "无历史记录"
			c.JSON(http.StatusOK, gin.H{
				"code":    resp4Device.Code,
				"err_msg": resp4Device.Err_msg,
				"data": resp4Device.Data,
			})
			return
		}
		resp4Device.Data = res
	}
	//else{
	//	datelist = svc.GetBetweenDates(snap_begin_time,snap_end_time)
	//	temp := model.Statistic{}
	//	for _,v := range datelist{
	//		res , err := svc.GetStatistic(v)
	//		if err != nil{
	//			//resp4Device.Code = -100
	//			//resp4Device.Err_msg = "无此日期数据"
	//			log4go.Error("无此日期数据")
	//			//zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
	//			continue
	//		}
	//		temp.Desayuno += res.Desayuno
	//		temp.Almuerzo += res.Almuerzo
	//		temp.Jantar += res.Jantar
	//	}
	//	resp4Device.Data = temp
	//}

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})
}

func GetOneStatisticRecords(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string

		datetime string
	)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	datetime = nowdate
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	//err := c.BindJSON(&Statistic_json)

	snap_begin_time := c.Query("snap_begin_time")
	snap_end_time := c.Query("snap_end_time")

	main_id,err := strconv.Atoi(c.Query("main_id"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数main_id错误"
	}
	item_id,err := strconv.Atoi(c.Query("item_id"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数item_id错误"
	}
	subject_id,err := strconv.Atoi(c.Query("subject_id"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数subject_id错误"
	}
	page,err := strconv.Atoi(c.Query("page"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数page错误"
	}
	size,err := strconv.Atoi(c.Query("size"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数size错误"
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if snap_begin_time == ""{
		snap_begin_time = datetime
	}
	if snap_end_time == ""{
		snap_end_time = datetime
	}

	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map

	snap_begin_time = svc.UTCchange(snap_begin_time)
	snap_end_time = svc.UTCchange(snap_end_time)
	if snap_begin_time[0:10] == snap_end_time[0:10]{
		item_map,gdata,err:=svc.Statistic_main_item_map()
		if err!=nil{

		}
		res ,count,total, err := svc.GetOneStatisticRecords(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,subject_id,page,size)
		if err != nil{
			resp4Device.Code = 0
			resp4Device.Err_msg = "无历史记录"
			c.JSON(http.StatusOK, gin.H{
				"code":    resp4Device.Code,
				"err_msg": resp4Device.Err_msg,
				"data": resp4Device.Data,
			})
			return
		}
		resp4Device.Data = res
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data": resp4Device.Data,
			"count":   count,
			"current": page,
			"size":    size,
			"total":   total,
		})
	}
	return
	//else{
	//	datelist = svc.GetBetweenDates(snap_begin_time,snap_end_time)
	//	temp := model.Statistic{}
	//	for _,v := range datelist{
	//		res , err := svc.GetStatistic(v)
	//		if err != nil{
	//			//resp4Device.Code = -100
	//			//resp4Device.Err_msg = "无此日期数据"
	//			log4go.Error("无此日期数据")
	//			//zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
	//			continue
	//		}
	//		temp.Desayuno += res.Desayuno
	//		temp.Almuerzo += res.Almuerzo
	//		temp.Jantar += res.Jantar
	//	}
	//	resp4Device.Data = temp
	//}

}

func GetStatisticRecord(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		datetime string
	)
	nowtime := time.Now()
	nowdate := nowtime.String()[0:10]
	datetime = nowdate
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	snap_begin_time := c.Query("snap_begin_time")
	snap_end_time := c.Query("snap_end_time")

	main_id,err := strconv.Atoi(c.Query("main_id"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数main_id错误"
	}
	item_id,err := strconv.Atoi(c.Query("item_id"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数item_id错误"
	}
	//subject_id,err := strconv.Atoi(c.Query("subject_id"))
	//if err != nil{
	//	resp4Device.Code = -100
	//	resp4Device.Err_msg = "参数subject_id错误"
	//}
	page,err := strconv.Atoi(c.Query("page"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数page错误"
	}
	size,err := strconv.Atoi(c.Query("size"))
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "参数size错误"
	}
	//page := 1
	//size := 5
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if snap_begin_time == "" || len(snap_begin_time)<10{
		snap_begin_time = datetime
	}
	if snap_end_time == ""|| len(snap_end_time)<10{
		snap_end_time = datetime
	}

	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map

	snap_begin_time = svc.UTCchange(snap_begin_time)
	snap_end_time = svc.UTCchange(snap_end_time)
	if snap_begin_time[0:10] == snap_end_time[0:10]{
		item_map,gdata,err:=svc.Statistic_main_item_map()
		if err!=nil{

		}
		res ,count,total, err := svc.GetOneStatisticOneRecord2(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
		if err != nil{
			fmt.Println(err)
			resp4Device.Code = 0
			resp4Device.Err_msg = "无历史记录"
			c.JSON(http.StatusOK, gin.H{
				"code":    resp4Device.Code,
				"err_msg": resp4Device.Err_msg,
				"data": resp4Device.Data,
				"count":   count,
				"current": page,
				"size":    size,
				"total":   total,
			})
			return
		}
		resp4Device.Data = res
		c.JSON(http.StatusOK, gin.H{
			"code":    resp4Device.Code,
			"err_msg": resp4Device.Err_msg,
			"data": resp4Device.Data,
			"count":   count,
			"current": page,
			"size":    size,
			"total":   total,
		})
	}
	return
}


func Time(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
	)

	res , err := svc.GetTime()
	if err != nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = "koala获取参数错误"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	resp4Device.Data = res

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})
}

func Test3(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	datetime := c.Query("datetime")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")


	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	if Position_map[-99] == nil{
		GetPositionMap()
	}
	Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	data,_:= svc.PrepareForstatistic2(datetime,Position_map_copy)
	resp4Device.Data = data
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
		"length": len(data),
	})
}
func Test4(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//datetime := c.Query("datetime")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")


	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	//if Position_map[-99] == nil{
	//	GetPositionMap()
	//}
	//Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	err:= svc.TruncateTables()
	if err != nil {
		zyutil.DeviceErrorReturn(c, -100,	err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
	})
}

func Test5(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	days := c.Query("days")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")

	day,_ := strconv.Atoi(days)
	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	//if Position_map[-99] == nil{
	//	GetPositionMap()
	//}
	//Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	err:= svc.DeleteRecord(day)
	if err != nil {
		zyutil.DeviceErrorReturn(c, -100,	err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
	})
}

func Test6(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//days := c.Query("days")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")

	//day,_ := strconv.Atoi(days)
	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	//if Position_map[-99] == nil{
	//	GetPositionMap()
	//}
	//Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	//err:= svc.DeleteRecord(day)
	//if err != nil {
	//	zyutil.DeviceErrorReturn(c, -100,	err.Error())
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":model.G_map.Getmap(),
		"data2":model.G_map_time.Getmap(),
	})
}



func Test8(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//days := c.Query("days")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")

	//day,_ := strconv.Atoi(days)
	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	//if Position_map[-99] == nil{
	//	GetPositionMap()
	//}
	//Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	svc.SyncMember()
	svc.SyncTeacher()
	//err:= svc.DeleteRecord(day)
	//if err != nil {
	//	zyutil.DeviceErrorReturn(c, -100,	err.Error())
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		//"data":model.G_map.Getmap(),
		//"data2":model.G_map_time.Getmap(),
	})
}

func Test7(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//Statistic_json model.Statistic_json
		//datelist []string
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//days := c.Query("days")
	//Snap_begin_time := c.Query("begin_time")
	//Snap_end_time := c.Query("end_time")

	//day,_ := strconv.Atoi(days)
	//err := c.BindJSON(&Statistic_json)
	//if err != nil {
	//	Statistic_json.Snap_begin_time= ""
	//	Statistic_json.Snap_end_time= ""
	//}
	//if Position_map[-99] == nil{
	//	GetPositionMap()
	//}
	//Position_map_copy := Position_map
	//datelist = svc.GetBetweenDates(Snap_begin_time,Snap_end_time)
	//for _,v := range datelist{
	//	svc.PrepareForstatistic(v,Position_map_copy)
	//}
	//err:= svc.DeleteRecord(day)
	//if err != nil {
	//	zyutil.DeviceErrorReturn(c, -100,	err.Error())
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":model.G_map.Getmap(),
		"data2":model.G_map_time.Getmap(),
	})
}



func Test(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//data           interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	tmp1 := c.Query("id1")
	tmp2 := c.Query("id2")
	id1, err := strconv.Atoi(tmp1)
	id2, err := strconv.Atoi(tmp2)

	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "无请求参数或请求参数类型错误!"
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	//_,err,normalcount,_ := svc.Getsubjectsall_normal(id,1,10000)
	//_,err,abnormalcount,_ := svc.Getsubjectsall_abnormal(id,1,10000)
	//if err != nil {
	//	resp4Device.Code = -100
	//	resp4Device.Err_msg = "查询失败"
	//}
	//svc.InsertPerson4Koala(id)
	_,err = koala.InsertPerson2PersonGroup(id1,[]int{id2})

	if err != nil{
		log.Error(err.Error())
		resp4Device.Err_msg = err.Error()
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	//results := make([]map[string]interface{},0)
	//result := map[string]interface{}{
	//	"id":id,
	//	"normalcount":normalcount,
	//	"abnormalcount":   abnormalcount,
	//}
	//results= append(results,result)
	//resp4Device.Data = results
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": resp4Device.Data,
	})
}

func Test2(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//data           interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	//var(
	//	message= make(map[string]interface{})
	//)
	////ip,err := Get_Ip()
	ips,err := svc.LocalIPv4s()
	if err != nil{
		log.Info("获取本机IP失败！", err)
		resp4Device.Code = -100
		resp4Device.Err_msg = "获取IP失败!"
	}
	//message["ip"] = ips[0]
	//c.JSON(message, err)

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	GetPositionMap()
	//svc.PrepareForRecord(1)
	//svc.GetIdentificationRecords_GroupBy()
	svc.PrepareForstatistic2("",Position_map)
	//statistic , _ := svc.GetStatistic("2020-07-19")
	time := svc.UTCchange("")

	//GetScreenList,_:= koala.GetScreenList()
	//employee , _ := koala.GetSubjects("employee")
	//visitor,_:=koala.GetSubjects("visitor")
	//_,_,err=koala.GetPersonData(employee)
	//visitor_data,_,err=koala.GetPersonData(visitor)
	//svc.PrepareForstatistic2("",Position_map_copy)
	if err != nil {
		fmt.Println(err)
	}
	//visitor,_,_:=koala.GetPersonData(visitor)
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data": ips[0],
		"data_map": time,
		//"GetScreenList": GetScreenList,
		//"employee": employee,
		//"visitor": visitor,
	})
}


func Test23(c *gin.Context) () {
	var (
		resp4Device    model.Resp4Device
		//data           interface{}
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	log4go.Info("met header is ",c.Request.Header)
	log4go.Info("met header is ",c)


	body, _ := ioutil.ReadAll(c.Request.Body)
	log4go.Info("req_body---------------------------------------------------")
	log4go.Info(body)
	log4go.Info(string(body))

	data := strings.ReplaceAll(string(body),"&quot;",`"`)

	m := make(map[string]interface{})
	json.Unmarshal([]byte(data), &m)
	log4go.Info("met body is ",m)
	//var(
	//	message= make(map[string]interface{})
	//)
	////ip,err := Get_Ip()
	//ips,err := svc.LocalIPv4s()
	//if err != nil{
	//	log.Info("获取本机IP失败！", err)
	//	resp4Device.Code = -100
	//	resp4Device.Err_msg = "获取IP失败!"
	//}
	//message["ip"] = ips[0]
	//c.JSON(message, err)

	//if resp4Device.Code != 0 {
	//	zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
	//	return
	//}
	//GetPositionMap()
	//svc.PrepareForRecord(1)
	//svc.GetIdentificationRecords_GroupBy()
	//svc.PrepareForstatistic2("",Position_map)
	//statistic , _ := svc.GetStatistic("2020-07-19")
	//time := svc.UTCchange("")

	//GetScreenList,_:= koala.GetScreenList()
	//employee , _ := koala.GetSubjects("employee")
	//visitor,_:=koala.GetSubjects("visitor")
	//_,_,err=koala.GetPersonData(employee)
	//visitor_data,_,err=koala.GetPersonData(visitor)
	//svc.PrepareForstatistic2("",Position_map_copy)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//visitor,_,_:=koala.GetPersonData(visitor)
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		//"data": ips[0],
		//"data_map": time,
		//"GetScreenList": GetScreenList,
		//"employee": employee,
		//"visitor": visitor,
	})
}


func Statistic_Project_Create(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
		Statistic_Main model.Statistic_Main_Json
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&Statistic_Main)
	if Statistic_Main.Name == "" {
		resp4Device.Code = -100
		resp4Device.Err_msg = "统计项目不能为空"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if len(Statistic_Main.Items) == 0{
		resp4Device.Code = -100
		resp4Device.Err_msg = "缺少子项"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if err != nil {

	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	Statistic_Project_Main,err := svc.Statistic_Project_Create(Statistic_Main.Name)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "新增记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	var parent_id = Statistic_Project_Main.Id
	for _,v := range Statistic_Main.Items{
		err = svc.Statistic_Project_Item_Create(parent_id,v)
		if err != nil {
			resp4Device.Code = -100
			resp4Device.Err_msg = "新增子项记录失败"
			svc.Statistic_Project_Items_Delete(parent_id)
			zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
			return
		}
	}
	Items ,err :=svc.Statistic_Project_Items_Get(parent_id)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	var res = map[string]interface{}{
		"id": Statistic_Project_Main.Id,
		"name": Statistic_Project_Main.Name,
		"items":Items,
	}
	resp4Device.Data = res
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
	})
}

func Statistic_Project_Get(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	if Persongrouplist_g[-99] == nil|| Accesslist_g[-99] == nil{
		GetPositionMap()
	}
	Statistic_Project_Main,err := svc.Statistic_Project_Get(Persongrouplist_g,Accesslist_g)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "获取记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	resp4Device.Data = Statistic_Project_Main
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
	})
}


func GetAccessList(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	koala.Init(configs.Gconf.KoalaHost)
	koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	//koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	AccessList,err := koala.GetAccessList()
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "获取记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if Accesslist_g[-99] == nil{
		for _,v := range AccessList{
			key ,_:= v["id"].(json.Number).Int()
			Accesslist_g[key]=v["name"]
		}
		Persongrouplist_g[-99]="activate"
	}
	AccessList = append(AccessList, map[string]interface{}{
		"id":-1,
		"name":"全相机",
	} )
	resp4Device.Data = AccessList
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
		//"Accesslist_g":    Accesslist_g,
	})
}

func GetPersonGroupList(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	//koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	PersonGroupList_subject,err := koala.GetPersonGroupList("0")
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "获取记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	PersonGroupList_visitor,err := koala.GetPersonGroupList("1")
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "获取记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if len(PersonGroupList_visitor) !=0{
		PersonGroupList_subject = append(PersonGroupList_subject,PersonGroupList_visitor...)
	}

	if Persongrouplist_g[-99] == nil{
		for _,v := range PersonGroupList_subject{
			key ,_:= v["id"].(json.Number).Int()
			Persongrouplist_g[key]=v["name"]
		}
		Persongrouplist_g[-99]="activate"
	}
	PersonGroupList_subject = append(PersonGroupList_subject, map[string]interface{}{
		"id":-1,
		"name":"全人员",
	} )
	resp4Device.Data = PersonGroupList_subject
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
		//"Persongrouplist_g":    Persongrouplist_g,
	})
}

func Statistic_Project_Delete(c *gin.Context) {
	var (
		resp4Device    model.Resp4Device
		data        interface{}
	)
	resp4Device.Code = 0
	resp4Device.Err_msg = ""

	tmp := c.Param("id")
	id, err := strconv.Atoi(tmp)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		log4go.Error(err.Error())
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	err = svc.Statistic_Project_Items_Delete(id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = ""
	}

	err = svc.Statistic_Project_Delete(id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = ""
	}

	err = svc.Statistics_Project_Delete(id)

	if err !=nil{
		resp4Device.Code = -100
		resp4Device.Err_msg = err.Error()
		data = ""
	}

	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	data = ""

	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    data,
	})
}

func Statistic_Project_Update(c *gin.Context){
	var (
		resp4Device    model.Resp4Device
		Statistic_Main model.Statistic_Main_Json
	)

	resp4Device.Code = 0
	resp4Device.Err_msg = ""
	err := c.BindJSON(&Statistic_Main)

	tmp := c.Param("id")
	parent_id, err := strconv.Atoi(tmp)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "传入的参数有误!"
		log4go.Error(err.Error())
	}

	if Statistic_Main.Name == "" {
		resp4Device.Code = -100
		resp4Device.Err_msg = "统计项目不能为空"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	if err != nil {
		return
	}
	if resp4Device.Code != 0 {
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	Statistic_Project_Main,err := svc.Statistic_Project_Update(Statistic_Main.Name,parent_id)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "修改记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}

	err = svc.Statistic_Project_Items_Delete(parent_id)

	for _,v := range Statistic_Main.Items{
		err = svc.Statistic_Project_Item_Create(parent_id,v)
		if err != nil {
			resp4Device.Code = -100
			resp4Device.Err_msg = "新增子项记录失败"
			svc.Statistic_Project_Items_Delete(parent_id)
			zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
			return
		}
	}

	Items ,err :=svc.Statistic_Project_Items_Get(parent_id)
	if err != nil {
		resp4Device.Code = -100
		resp4Device.Err_msg = "查询记录失败"
		zyutil.DeviceErrorReturn(c, resp4Device.Code, resp4Device.Err_msg)
		return
	}
	var res = map[string]interface{}{
		"id": Statistic_Project_Main.Id,
		"name": Statistic_Project_Main.Name,
		"items":Items,
	}
	resp4Device.Data = res
	c.JSON(http.StatusOK, gin.H{
		"code":    resp4Device.Code,
		"err_msg": resp4Device.Err_msg,
		"data":    resp4Device.Data,
	})
}





