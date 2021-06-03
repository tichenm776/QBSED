package service

import (
	"github.com/alecthomas/log4go"
	"strconv"
	"time"
	"zhiyuan/classin/internal/classin"
)

func (s *Service) CreateStudent()(){



}


func (s *Service) RegisterPerson(telephone,nickname string)(int ,error){
//func  RegisterPerson(telephone,nickname string)(){
//是否加入机构

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"telephone":telephone,
		"nickname":nickname,
		"password":telephone,
		//"md5pass":"",
		//"addToSchoolMember":telephone,
	}
	data,err := classin.G_Class.RPC("register",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("data").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}

func (s *Service) EditStudent(studentUid,studentName string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"studentUid":studentUid,
		"studentName":studentName,
	}
	data,err := classin.G_Class.RPC("editSchoolStudent",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return	0 ,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}
func (s *Service) AddSchoolStudent(studentAccount,studentName string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"studentAccount":studentAccount,
		"studentName":studentName,
	}
	data,err := classin.G_Class.RPC("addSchoolStudent",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0, err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}

func (s *Service) AddTeacher(teacherAccount,teacherName string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"teacherAccount":teacherAccount,
		"teacherName":teacherName,
	}
	data,err := classin.G_Class.RPC("addTeacher",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0, err
	}
	if value,err :=data.Get("data").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}


func (s *Service) EditTeacher(teacherUid,teacherName string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"teacherUid":teacherUid,
		"teacherName":teacherName,
	}
	data,err := classin.G_Class.RPC("editTeacher",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return	0 ,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}

//3464235
func (s *Service) RestartUsingTeacher(teacherUid string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"teacherUid":teacherUid,
		//"teacherName":teacherName,
	}
	data,err := classin.G_Class.RPC("restartUsingTeacher",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return	0 ,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}

func (s *Service) StopUsingTeacher(teacherUid string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"teacherUid":teacherUid,
		//"teacherName":teacherName,
	}
	data,err := classin.G_Class.RPC("stopUsingTeacher",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return	0 ,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		return 0, err
	}
}








//editTeacher
//https://api.eeo.cn/partner/api/course.api.php?action=editSchoolStudent
//https://root_url/partner/api/course.api.php?action=addSchoolStudent