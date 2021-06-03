package service

import (
	"fmt"
	"testing"
)



func TestRegisterPerson(t *testing.T) {
	svv := Service{}
	classin_UID,err :=svv.RegisterPerson("18357036164","TIM")
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

