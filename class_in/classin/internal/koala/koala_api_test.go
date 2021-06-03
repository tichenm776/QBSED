package koala

import (
	"fmt"
	"testing"
)

func TestGetTimer(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}
	arr,err :=GetTimer()
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestGetAccessList(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}
	arr,err :=GetAccessList()
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}


func TestGetAccess(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}
	arr,err :=GetAccess(1)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestGetPerson(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}
	arr,subjecttype,err :=GetPerson(3)
	if err != nil{
		fmt.Println(err)
		fmt.Println(subjecttype)
	}else{
		fmt.Println(arr)
		fmt.Println(subjecttype)
	}

}

func TestCreateAccess(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}
	name := "test"
	comment := "vivlavida"
	arr,err :=CreateAccess(name,comment)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestDeleteAccess(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 41
	arr,err :=DeleteAccess(id)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestInsertScreen2Access(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 4
	screenids := []int{9,12,13,14,15,16}
	arr,err :=InsertScreen2Access(id,screenids)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestDeleteScreen2Access(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 4
	screenids := []int{9,12,13,14,15,16}
	arr,err :=DeleteScreen4Access(id,screenids)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}
func TestCreatePersonGroup(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	//id := strconv.Itoa(4)
	//screenids := []int{9,12,13,14,15,16}
	name :="testc"
	comment := "vivalavida"
	arr,err :=CreatePersonGroup(name,comment)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}


func TestDeletePersonGroup(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	id := 9
	//id := strconv.Itoa(9)
	//screenids := []int{9,12,13,14,15,16}
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=DeletePersonGroup(id)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestInsertPerson2PersonGroup(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 42
	//subjectids := []int{8,7,6,4,3,1}
	subjectids := []int{3}
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=InsertPerson2PersonGroup(id,subjectids)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}
	//{"subject_ids":[6]}
	//{"subject_ids":[6]}
}


func TestDeletePerson2PersonGroup(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 42
	subjectids := []int{3}
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=DeletePerson2PersonGroup(id,subjectids)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestDeleteAccessControlRules(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 12
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=DeleteAccessControlRules(id)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}
func TestGetPersonGroupList(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	//id := 12
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=GetPersonGroupList("0")
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}

func TestGetPersonGroupListByGroupId(t *testing.T) {
	//for _, ut := range ucTests {
	//	uc := UpperCase(ut.in)
	//	if uc != ut.out {
	//		t.Errorf("uppercase(%s) = %s,must be %s", ut.in, uc, ut.out)
	//	}
	//}

	//id := "3"
	id := 11
	//name :="testc"
	//comment := "vivalavida"
	arr,err :=GetPersonGroupListByGroupId(id)
	if err != nil{
		fmt.Println(err)
	}else{
		fmt.Println(arr)
	}

}