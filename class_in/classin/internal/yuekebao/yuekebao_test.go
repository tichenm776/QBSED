package yuekebao

import (
	"fmt"
	"testing"
)

func Test_GetMember(t *testing.T) {
	data,err := GetAllMember()
	if err != nil{
		fmt.Println("err is",err)
	}
	fmt.Println("data is",data)
}

func Test_GetTeacher(t *testing.T) {
	data,err := GetAllTeacher()
	if err != nil{
		fmt.Println("err is",err)
	}
	fmt.Println("data is",data)
}

func Test_GetAllCourse(t *testing.T) {
	data,err := GetAllCourse()
	if err != nil{
		fmt.Println("err is",err)
	}
	fmt.Println("data is",data)
}

