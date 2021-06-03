package service

import (
	"fmt"
	"testing"
)



func TestRegisterPerson(t *testing.T) {
	svv := Service{}
	classin_UID,err :=svv.RegisterPerson("18357036164","TIM","1")
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println(classin_UID)

	}
func TestEditStudent(t *testing.T) {
	svv := Service{}
	svv.EditStudent("34296792","TIM")
}
func TestAddSchoolStudent(t *testing.T) {
	svv := Service{}
	svv.AddSchoolStudent("18357036164","TIM")
}
func TestAddTeacher(t *testing.T) {
	svv := Service{}
	svv.AddTeacher("18357036164","TIM")
}
func TestRestartUsingTeacher(t *testing.T) {
	svv := Service{}
	svv.RestartUsingTeacher("3464235")
}
func TestStopUsingTeacher(t *testing.T) {
	svv := Service{}
	svv.StopUsingTeacher("3464235")
}


func TestTest2(t *testing.T) {
	svv := Service{}
	svv.Test2(5)
}

func TestCreateCourse(t *testing.T) {
	svv := Service{}
	svv.CreateCourse("test2","")
}

func TestEditCourse(t *testing.T) {
	svv := Service{}
	svv.EditCourse("test3","","154564453")
}

func TestCreateCourseClass(t *testing.T) {
	svv := Service{}


	courseId:="154564453"
	className:="test2"
	beginTime:="1622437500"
	endTime:="1622440500"
	teacherUid:="34825556"
	folderId:=""
	teachMode:="1"
	isAutoOnstage:="0"
	seatNum:="0"
	isHd := "0"





	svv.CreateCourseClass(courseId,className,beginTime,endTime,teacherUid,folderId,teachMode,isAutoOnstage,seatNum,isHd)
}


func TestGetTopFolderId(t *testing.T) {
	svv := Service{}
	result,_:=svv.GetTopFolderId()
	fmt.Println(result)
}
