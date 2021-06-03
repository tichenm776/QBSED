package service

import (
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"strconv"
	"time"
	"zhiyuan/QBSED/internal/classin"
)



func (s *Service) RegisterPerson(telephone,nickname,membertype string)(int ,error){
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
		"addToSchoolMember":membertype,
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
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
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
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
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

func (s *Service) Test2(number int)(int,error){
	if number%5 == 0{
		fmt.Println(number)
		n1 := number/5
		fmt.Println(n1)
		n2 := (n1-1)*5
		fmt.Println(n2)
		n3 := number - 1
		fmt.Println(n3)
	}
	return 0, nil
}




func(s *Service)CreateCourse(courseName,mainTeacherUid string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseName":courseName,
		"mainTeacherUid":mainTeacherUid,
	}
	data,err := classin.G_Class.RPC("addCourse",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("data").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}

func(s *Service)EditCourse(courseName,mainTeacherUid,courseId string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseName":courseName,
		"mainTeacherUid":mainTeacherUid,
		"courseId":courseId,
	}
	data,err := classin.G_Class.RPC("editCourse",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}

func(s *Service)EndCourse(courseId string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
	}
	data,err := classin.G_Class.RPC("endCourse",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}


func(s *Service)CreateCourseClass(courseId,className,beginTime,endTime,teacherUid,folderId,teachMode,isAutoOnstage,seatNum,isHd string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
		"className":className,
		"beginTime":beginTime,
		"endTime":endTime,
		"teacherUid":teacherUid,
		"folderId":folderId,
		"teachMode":teachMode,
		"isAutoOnstage":isAutoOnstage,
		"seatNum":seatNum,
		"isHd":isHd,
	}
	data,err := classin.G_Class.RPC("addCourseClass",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("data").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}
func(s *Service)EditCourseClass(courseId,classId,className,beginTime,endTime,teacherUid,folderId,teachMode,isAutoOnstage,seatNum,isHd string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
		"classId":classId,
		"className":className,
		"beginTime":beginTime,
		"endTime":endTime,
		"teacherUid":teacherUid,
		"folderId":folderId,
		"teachMode":teachMode,
		"isAutoOnstage":isAutoOnstage,
		"seatNum":seatNum,
		"isHd":isHd,
	}
	data,err := classin.G_Class.RPC("editCourseClass",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}
func(s *Service)ModifyClassSeatNum(courseId,classId,seatNum,isHd string)(int ,error){

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
		"classId":classId,
		"seatNum":seatNum,
		"isHd":isHd,
		//"isDc":isHd,
	}
	data,err := classin.G_Class.RPC("modifyClassSeatNum",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("error_info").Get("errno").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}

}


func (s *Service) GetgetFolderList()(int ,error){
	//func  RegisterPerson(telephone,nickname string)(){
	//是否加入机构

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
	}
	data,err := classin.G_Class.RPC_cloud("getFolderList",params)
	if err != nil{
		log4go.Error("getFolderList err:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("data").Int();err == nil{
		return value,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}
}
func (s *Service) GetTopFolderId()(int ,error){
	//func  RegisterPerson(telephone,nickname string)(){
	//是否加入机构

	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
	}
	data,err := classin.G_Class.RPC_cloud("getTopFolderId",params)
	if err != nil{
		log4go.Error("getFolderList err:",err)
		return 0,err
	}
	//fmt.Println(data)
	if value,err :=data.Get("data").String();err == nil{
		foldid,_:=strconv.Atoi(value)
		return foldid,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}
}

func (s *Service) AddClassStudentMultiple(courseId,classId,studentJson string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
		"classId":classId,
		"studentJson":studentJson,
	}
	data,err := classin.G_Class.RPC("addClassStudentMultiple",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0, err
	}
	if value,err :=data.Get("error_info").Get("errno").String();err == nil{

		result ,_:= strconv.Atoi(value)

		return result,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}
}

func (s *Service) DelClassStudentMultiple(courseId,classId,studentJson string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
		"classId":classId,
		"studentJson":studentJson,
	}
	data,err := classin.G_Class.RPC("delClassStudentMultiple",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0, err
	}
	if value,err :=data.Get("error_info").Get("errno").String();err == nil{

		result ,_:= strconv.Atoi(value)

		return result,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}
}

func (s *Service) DelCourseClass(courseId string)(int,error){
	Safekey := classin.G_Class.Md5V()
	timeStamp:=strconv.FormatInt(time.Now().Unix(),10)
	params := map[string]interface{}{
		"SID":classin.SID,
		"safeKey":Safekey,
		"timeStamp":timeStamp,
		"courseId":courseId,
	}
	data,err := classin.G_Class.RPC("endCourse",params)
	if err != nil{
		log4go.Error("classin注册失败:",err)
		return 0, err
	}
	if value,err :=data.Get("error_info").Get("errno").String();err == nil{

		result ,_:= strconv.Atoi(value)

		return result,nil
	}else{
		err_msg,_ := data.Get("error_info").Get("error").String()
		return 0, errors.New(err_msg)
	}
}
//https://root_url/partner/api/course.api.php?action=


//editTeacher
//https://api.eeo.cn/partner/api/course.api.php?action=editSchoolStudent
//https://root_url/partner/api/course.api.php?action=addSchoolStudent