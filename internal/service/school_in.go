package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"strconv"
	"strings"
	"time"
	"zhiyuan/QBSED/internal/model"
	"zhiyuan/QBSED/internal/yuekebao"
)

func (s *Service)SyncMember()(error){

	datafromyuekebao,err := yuekebao.GetAllMember()
	log4go.Info("datafromyuekebao length is",len(datafromyuekebao))
	if err != nil{
		log4go.Error("get data from yuekebao err:",err)
		return err
	}
	log4go.Info("datafromyuekebao len is",datafromyuekebao)
	if len(datafromyuekebao) == 0{
		log4go.Error("get data need to sync")
		return err
	}
	datafromlocaldb,err := s.GetStudentsFromDb()
	if err != nil{
		log4go.Error("get data from localdb err:",err)
		return err
	}
	log4go.Info("datafromlocaldb len is",datafromlocaldb)
	number := 0
	if len(datafromlocaldb) <= 0{
		number = 1
	}else{
		number = 2
	}
	//compare
	log4go.Info("number is",number)
	switch number {
	case 1:
		fmt.Println("a")
		name := ""
		phone := ""
		id := ""
		for _,v := range datafromyuekebao{
			name = v["name"].(string)
			id = v["id"].(string)
			phone = v["phone"].(string)
		if phone == ""{
			log4go.Error("phone number is null",id)
			continue
		}
		if strings.Index(phone,"86-") != -1{
			phone = PhoneParse(phone)
		}
		SID_classin,err := s.RegisterPerson(phone,name,"1")
		if err != nil{
			log4go.Error("class in register student err",err)
		}
		result,err := s.AddSchoolStudent(phone,name)
		if err != nil{
			log4go.Error("class in register student err",err)
		}
		if result == 1 || result == 135 || result == 133{
			log4go.Info("success")
		}else{
			log4go.Info("result code is:",result)
			continue
		}
		DB_id ,_ := strconv.Atoi(id)
		DB_SID := strconv.Itoa(SID_classin)
		addmember := model.Member{
			Id: DB_id,
			Phone: phone,
			Name: name,
			SID: DB_SID,
			Member_Type: "1",
			//IsChange: 1,
		}
		err = s.AddMember(addmember)
		if err != nil{
			log4go.Error("add student to local db err",err)
			return err
		}
	}
		return nil
	case 2:
		fmt.Println("b")
		name := ""
		phone := ""
		id := ""
		tempmap := make(map[string]interface{},0)
		for _,v := range datafromlocaldb{
			tempmap[strconv.Itoa(v.Id)]=v
		}
		for _,v := range datafromyuekebao{
			name = v["name"].(string)
			id = v["id"].(string)
			phone = v["phone"].(string)
			if phone == ""{
				log4go.Error("phone number is null",id)
				continue
			}
			if strings.Index(phone,"86-") != -1{
				phone = PhoneParse(phone)
			}
			if tempmap[id] != nil{
				memberobj := tempmap[id].(model.Member)
				if memberobj.Name != name{
					//is exist but value changed we update
					_,err := s.EditStudent(memberobj.SID,name)
					if err != nil{
						log4go.Error("update student name err:",err)
						log4go.Info("student phone is:",memberobj.Phone)
						continue
					}
					memberobj.Name = name
					err = s.UpdateMember(memberobj)
					if err != nil{
						log4go.Error("local db - update student name err:",err)
						log4go.Info("student phone is:",memberobj.Phone)
						continue
					}
				}
			}else{
				// not exist so we create
				SID_classin,err := s.RegisterPerson(phone,name,"1")
				if err != nil{
					log4go.Error("class in register student err",err)
				}
				DB_id ,_ := strconv.Atoi(id)
				DB_SID := strconv.Itoa(SID_classin)
				result,err := s.AddSchoolStudent(phone,name)
				if err != nil{
					log4go.Error("class in register student err",err)
				}
				if result == 1 || result == 135 || result == 133{
					log4go.Info("success")
				}else{
					log4go.Info("result code is:",result)
					continue
				}
				addmember := model.Member{
					Id: DB_id,
					Phone: phone,
					Name: name,
					SID: DB_SID,
					Member_Type: "1",
					//IsChange: 1,
				}
				err = s.AddMember(addmember)
				if err != nil{
					log4go.Error("add student to local db err",err)
					return err
				}
			}
		}
	default:
		log4go.Info("SyncMember default")
	}
	return nil
}

func (s *Service)SyncTeacher()(error){

	datafromyuekebao,err := yuekebao.GetAllTeacher()
	log4go.Info("datafromyuekebao length is",len(datafromyuekebao))
	if err != nil{
		log4go.Error("get data from yuekebao err:",err)
		return err
	}
	log4go.Info("datafromyuekebao len is",datafromyuekebao)
	if len(datafromyuekebao) == 0{
		log4go.Error("get data need to sync")
		return err
	}
	datafromlocaldb,err := s.GetTeachersFromDb()
	if err != nil{
		log4go.Error("get data from localdb err:",err)
		return err
	}
	log4go.Info("datafromlocaldb len is",datafromlocaldb)
	number := 0
	if len(datafromlocaldb) <= 0{
		number = 1
	}else{
		number = 2
	}
	//compare
	log4go.Info("number is",number)
	switch number {
	case 1:
		fmt.Println("a")
		name := ""
		phone := ""
		id := ""
		for _,v := range datafromyuekebao{
			name = v["name"].(string)
			id = v["id"].(string)
			phone = v["phone"].(string)
			if phone == ""{
				log4go.Error("phone number is null",v)
				continue
			}
			if strings.Index(phone,"86-") != -1{
				phone = PhoneParse(phone)
			}
			SID_classin,err := s.RegisterPerson(phone,name,"2")
			if err != nil{
				log4go.Error("class in register teacher err",err)
				log4go.Info("err person is:",v)
				continue
			}
			result,err := s.AddTeacher(phone,name)
			if err != nil{
				log4go.Error("class in register teacher err",err)
			}
			if result == 1 {
			//if result == 1 || result == 135 || result == 133{
				log4go.Info("success")
			}else{
				log4go.Info("result code is:",result)
				//continue
			}
			DB_id ,_ := strconv.Atoi(id)
			DB_SID := strconv.Itoa(SID_classin)
			addmember := model.Member{
				Id: DB_id,
				Phone: phone,
				Name: name,
				SID: DB_SID,
				Member_Type: "2",
				//IsChange: 1,
			}
			err = s.AddMember(addmember)
			if err != nil{
				log4go.Error("add student to local db err",err)
				return err
			}
		}
		return nil
	case 2:
		fmt.Println("b")
		name := ""
		phone := ""
		id := ""
		tempmap := make(map[string]interface{},0)
		for _,v := range datafromlocaldb{
			tempmap[strconv.Itoa(v.Id)]=v
		}
		for _,v := range datafromyuekebao{
			name = v["name"].(string)
			id = v["id"].(string)
			phone = v["phone"].(string)
			if phone == ""{
				log4go.Error("phone number is null",v)
				continue
			}
			if strings.Index(phone,"86-") != -1{
				phone = PhoneParse(phone)
			}
			if tempmap[id] != nil{
				memberobj := tempmap[id].(model.Member)
				if memberobj.Name != name{
					//is exist but value changed we update
					_,err := s.EditTeacher(memberobj.SID,name)
					if err != nil{
						log4go.Error("update student name err:",err)
						log4go.Info("student phone is:",memberobj.Phone)
						continue
					}
					memberobj.Name = name
					err = s.UpdateMember(memberobj)
					if err != nil{
						log4go.Error("local db - update student name err:",err)
						log4go.Info("student phone is:",memberobj.Phone)
						continue
					}
				}
			}else{
				// not exist so we create

				SID_classin,err := s.RegisterPerson(phone,name,"2")
				if err != nil{
					log4go.Error("class in register student err",err)
					log4go.Info("err person is:",v)
					continue
				}
				DB_id ,_ := strconv.Atoi(id)
				DB_SID := strconv.Itoa(SID_classin)
				result,err := s.AddTeacher(phone,name)
				if err != nil{
					log4go.Error("class in register student err",err)
				}
				//if result == 1 || result == 135 || result == 133{
				if result == 1{
					log4go.Info("success")
				}else{
					log4go.Info("result code is:",result)
					//continue
				}
				addmember := model.Member{
					Id: DB_id,
					Phone: phone,
					Name: name,
					SID: DB_SID,
					Member_Type: "2",
					//IsChange: 1,
				}
				err = s.AddMember(addmember)
				if err != nil{
					log4go.Error("add student to local db err",err)
					return err
				}
			}
		}
	default:
		log4go.Info("SyncMember default")
	}
	return nil
}

func PhoneParse(phonenumber string)(string){

	phone_arr := strings.Split(phonenumber,"-")

	//if len(phone_arr) > 0{
	log4go.Info(phone_arr)
	return phone_arr[1]
	//}
	//return phonenumber
}

//func (s *Service)SyncCourse()(error){
//
//	datafromyuekebao,err := yuekebao.GetAllCourse()
//	log4go.Info("datafromyuekebao length is",len(datafromyuekebao))
//	if err != nil{
//		log4go.Error("get data from yuekebao err:",err)
//		return err
//	}
//	log4go.Info("datafromyuekebao len is",datafromyuekebao)
//	if len(datafromyuekebao) == 0{
//		log4go.Error("get data need to sync")
//		return err
//	}
//	teacherfromlocaldb,err := s.GetTeachersFromDb()
//	if err != nil{
//		log4go.Error("get data from localdb err:",err)
//		return err
//	}
//	log4go.Info("datafromlocaldb len is",teacherfromlocaldb)
//
//	tmpteacher := make(map[string]interface{})
//
//	for _,v := range teacherfromlocaldb{
//		id := strconv.Itoa(v.Id)
//		tmpteacher[id] = v
//	}
//	studentfromlocaldb,err := s.GetStudentsFromDb()
//	if err != nil{
//		log4go.Error("get data from localdb err:",err)
//		return err
//	}
//	log4go.Info("datafromlocaldb len is",studentfromlocaldb)
//
//	tmpstudent := make(map[string]interface{})
//
//	for _,v := range studentfromlocaldb{
//		id := strconv.Itoa(v.Id)
//		tmpstudent[id] = v
//	}
//	datafromlocaldb,err := s.GetCourses()
//	if err != nil{
//		log4go.Error("get data from localdb err:",err)
//		return err
//	}
//	log4go.Info("datafromlocaldb len is",datafromlocaldb)
//
//	number := 0
//	if len(datafromlocaldb) <= 0{
//		number = 1
//	}else{
//		number = 2
//	}
//	//compare
//	log4go.Info("number is",number)
//	switch number {
//	case 1:
//		fmt.Println("a")
//		name,teacher_id_yuekebao,teacher_id_classin,start_day,start_end_time,folderId,start_time,end_time := "","","","","","","",""
//		teachMode:="1"
//		isAutoOnstage:="0"
//		seatNum:="0"
//		isHd := "0"
//		id := ""
//		for _,v := range datafromyuekebao{
//			name = v["category_name"].(string)
//			id = v["id"].(string)
//			teacher_id_yuekebao = v["teacher"].(string)
//			start_day = v["start_day"].(string)
//			start_end_time = v["start_end_time"].(string)
//			seatNum = v["num"].(string)
//			orderArr := make([]interface{},0)
//			if seatNum == "0"{
//				log4go.Info("orderArr is empty")
//			}else{
//				orderArr = v["orderArr"].([]interface{})
//			}
//			//处理时间
//			timearr := strings.Split(start_end_time,"-")
//			if len(timearr) == 2{
//				temp_start_time := start_day + " " +timearr[0]
//				temp_end_time := start_day + " " +timearr[1]
//				start_time , err = GetUnixTime(temp_start_time)
//				if err != nil{
//					log4go.Error("时间转换失败:",err)
//				}
//				end_time , err = GetUnixTime(temp_end_time)
//				if err != nil{
//					log4go.Error("时间转换失败:",err)
//				}
//			}
//			result,err := s.CreateCourse(name,"")
//			if err != nil{
//				log4go.Error("class in create course err",err)
//				log4go.Info("err course is:",v)
//				//return err
//				continue
//			}
//			courseId := strconv.Itoa(result)
//			if tmpteacher[teacher_id_yuekebao] != nil{
//				teacher_id_classin = tmpteacher[teacher_id_yuekebao].(model.Member).SID
//			}
//			foldeid,err := s.GetTopFolderId()
//			if err != nil{
//				log4go.Error("class in get FolderList err",err)
//				//log4go.Info("err course is:",v)
//				//return err
//			}
//			folderId = strconv.Itoa(foldeid)
//			class_result,err := s.CreateCourseClass(courseId,name,start_time,end_time,teacher_id_classin,folderId,
//				teachMode,isAutoOnstage,seatNum,isHd)
//			//result,err := s.AddTeacher(phone,name)
//			if err != nil{
//				log4go.Error("class in CreateCourseClass  err",err)
//				continue
//			}
//			classId := strconv.Itoa(class_result)
//			//if result == 1 {
//				//if result == 1 || result == 135 || result == 133{
//				//log4go.Info("success")
//			//}else{
//			//	log4go.Info("result code is:",result)
//				//continue
//			//}
//			//课节添加学生
//			studentJson := make([]map[string]interface{},0)
//			studentids:=make([]string,0)
//			studentsids:=make([]string,0)
//			if len(orderArr) > 0{
//
//				for _,v:=range orderArr{
//					if obj,ok := v.(map[string]interface{});ok{
//						studentid := obj["member"].(string)
//						if tmpstudent[studentid]!=nil{
//							uid := tmpstudent[studentid].(model.Member).SID
//							studentJson = append(studentJson, map[string]interface{}{"uid":uid} )
//							studentsids = append(studentsids, uid)
//							studentids = append(studentids, studentid)
//						}
//					}
//				}
//			}
//			STUDENT_JSON := ""
//			studentJson_bytes,err := json.Marshal(studentJson)
//			if err != nil{
//				log4go.Error("studentJson marshal err",err)
//			}else{
//				STUDENT_JSON = string(studentJson_bytes)
//			}
//			addsuccess,err := s.AddClassStudentMultiple(courseId,classId,STUDENT_JSON)
//			if err != nil{
//				log4go.Error("add class student multiple err",err)
//			}
//			if addsuccess == 1{
//				log4go.Info("class add student success")
//			}
//			//id_yuekebao courseid classid name starttime endtime teacherid teachersid studentsid studentid
//			student_id :=Getids(studentids)
//			student_sid :=Getids(studentsids)
//			DB_id ,_ := strconv.Atoi(id)
//			DB_SID := courseId
//			addCourse := model.Course{
//				Id: DB_id,
//				Name: name,
//				Course_id: DB_SID,
//				Class_id:classId,
//				Starttime:start_time,
//				Endtime:end_time,
//				Teacher_id:teacher_id_yuekebao,
//				Teacher_sid:teacher_id_classin,
//				Student_id:student_id,
//				Student_sid:student_sid,
//				SeatNum:seatNum,
//			}
//			err = s.AddCourse(addCourse)
//			if err != nil{
//				log4go.Error("add student to local db err",err)
//				return err
//			}
//		}
//		return nil
//	case 2:
//		fmt.Println("b")
//
//		//取交并差
//		yuekebao_ids := make([]string,0)
//		classin_ids := make([]string,0)
//		for _,v := range datafromyuekebao{
//			id := v["id"].(string)
//			yuekebao_ids = append(yuekebao_ids, id)
//		}
//
//		for _,v := range datafromyuekebao{
//			id := v["id"].(string)
//			yuekebao_ids = append(yuekebao_ids, id)
//		}
//		for _,v := range datafromlocaldb{
//			classin_ids = append(classin_ids, strconv.Itoa(v.Id))
//		}
//
//		allids := union(yuekebao_ids,classin_ids)
//
//		addids := difference(allids,classin_ids)
//		updateids := intersect(yuekebao_ids,classin_ids)
//		deleteids := difference(allids,yuekebao_ids)
//
//		addmap :=make(map[string]interface{},0)
//		updatemap :=make(map[string]interface{},0)
//		//deletemap :=make(map[string]interface{},0)
//
//		if len(addids) > 0{
//			for _,v:=range addids{
//				addmap[v]=1
//			}
//		}
//		if len(updateids) > 0{
//			for _,v:=range addids{
//				updatemap[v]=1
//			}
//		}
//		if len(deleteids) > 0{
//			for _,v:=range deleteids{
//				deleteobj,err := s.GetCoursesbyid(v)
//				if err != nil{
//					log4go.Error("get course from db err:",err)
//					continue
//				}
//				_ , err = s.EndCourse(deleteobj.Course_id)
//				if err!=nil{
//					log4go.Error("end course err:",err)
//					continue
//				}
//			}
//		}
//		name,teacher_id_yuekebao,teacher_id_classin,start_day,start_end_time,folderId,start_time,end_time := "","","","","","","",""
//		teachMode:="1"
//		isAutoOnstage:="0"
//		seatNum:="0"
//		isHd := "0"
//		id := ""
//		for _,v := range datafromyuekebao {
//			name = v["category_name"].(string)
//			id = v["id"].(string)
//			teacher_id_yuekebao = v["teacher"].(string)
//			start_day = v["start_day"].(string)
//			start_end_time = v["start_end_time"].(string)
//			seatNum = v["num"].(string)
//			orderArr := make([]interface{}, 0)
//			if seatNum == "0" {
//				log4go.Info("orderArr is empty")
//			} else {
//				orderArr = v["orderArr"].([]interface{})
//			}
//			//处理时间
//			timearr := strings.Split(start_end_time, "-")
//			if len(timearr) == 2 {
//				temp_start_time := start_day + " " + timearr[0]
//				temp_end_time := start_day + " " + timearr[1]
//				start_time, err = GetUnixTime(temp_start_time)
//				if err != nil {
//					log4go.Error("时间转换失败:", err)
//				}
//				end_time, err = GetUnixTime(temp_end_time)
//				if err != nil {
//					log4go.Error("时间转换失败:", err)
//				}
//			}
//			//add or update or delete
//			if addmap[id] != nil{
//				err := s.AddCourseToClassin(id,name,start_time,end_time,teachMode,isAutoOnstage,seatNum,isHd,teacher_id_yuekebao ,teacher_id_classin,folderId,tmpteacher,tmpstudent ,orderArr)
//				if err != nil{
//					log4go.Error(err)
//				}
//				continue
//				}
//			if updatemap[id] != nil{
//				s.UpdateCourseToClassin(id,name,start_time,end_time,teachMode,isAutoOnstage,seatNum,isHd,teacher_id_yuekebao ,teacher_id_classin,folderId,tmpteacher,tmpstudent ,orderArr)
//				if err != nil{
//					log4go.Error(err)
//				}
//				continue
//			}
//	}
//
//
//		return nil
//	default:
//		log4go.Info("SyncMember default")
//	}
//	return nil
//}

func(s *Service) AddCourseToClassin(id,name,start_time,end_time,teachMode,isAutoOnstage,seatNum,isHd,teacher_id_yuekebao,teacher_id_classin,folderId string ,tmpteacher,tmpstudent map[string]interface{},orderArr []interface{})(error){

	result,err := s.CreateCourse(name,"")
	if err != nil{
		log4go.Error("class in create course err",err)
		log4go.Info("err course is:",name)
		//return err
		return err
	}
	courseId := strconv.Itoa(result)
	if tmpteacher[teacher_id_yuekebao] != nil{
		teacher_id_classin = tmpteacher[teacher_id_yuekebao].(model.Member).SID
	}
	foldeid,err := s.GetTopFolderId()
	if err != nil{
		log4go.Error("class in get FolderList err",err)
		//log4go.Info("err course is:",v)
		//return err
	}
	folderId = strconv.Itoa(foldeid)
	class_result,err := s.CreateCourseClass(courseId,name,start_time,end_time,teacher_id_classin,folderId,
		teachMode,isAutoOnstage,seatNum,isHd)
	//result,err := s.AddTeacher(phone,name)
	if err != nil{
		log4go.Error("class in CreateCourseClass  err",err)
		return err
	}
	classId := strconv.Itoa(class_result)

	//课节添加学生
	studentJson := make([]map[string]interface{},0)
	studentids:=make([]string,0)
	studentsids:=make([]string,0)
	if len(orderArr) > 0{
		for _,v:=range orderArr{
			if obj,ok := v.(map[string]interface{});ok{
				studentid := obj["member"].(string)
				if tmpstudent[studentid]!=nil{
					uid := tmpstudent[studentid].(model.Member).SID
					studentJson = append(studentJson, map[string]interface{}{"uid":uid} )
					studentsids = append(studentsids, uid)
					studentids = append(studentids, studentid)
				}
			}
		}
	}
	STUDENT_JSON := ""
	studentJson_bytes,err := json.Marshal(studentJson)
	if err != nil{
		log4go.Error("studentJson marshal err",err)
	}else{
		STUDENT_JSON = string(studentJson_bytes)
	}
	addsuccess,err := s.AddClassStudentMultiple(courseId,classId,STUDENT_JSON)
	if err != nil{
		log4go.Error("add class student multiple err",err)
	}
	if addsuccess == 1{
		log4go.Info("class add student success")
	}
	//id_yuekebao courseid classid name starttime endtime teacherid teachersid studentsid studentid
	student_id :=Getids(studentids)
	student_sid :=Getids(studentsids)
	DB_id ,_ := strconv.Atoi(id)
	DB_SID := courseId
	addCourse := model.Course{
		Id: DB_id,
		Name: name,
		Course_id: DB_SID,
		Class_id:classId,
		Starttime:start_time,
		Endtime:end_time,
		Teacher_id:teacher_id_yuekebao,
		Teacher_sid:teacher_id_classin,
		Student_id:student_id,
		Student_sid:student_sid,
	}
	err = s.AddCourse(addCourse)
	if err != nil{
		log4go.Error("add student to local db err",err)
		return err
	}
	return nil
}


//func(s *Service) UpdateCourseToClassin(id,name,start_time,end_time,teachMode,isAutoOnstage,seatNum,isHd,teacher_id_yuekebao,teacher_id_classin,folderId string ,tmpteacher,tmpstudent map[string]interface{},orderArr []interface{})(error){
//
//	obj,err := s.GetCoursesbyid(id)
//	if obj.Teacher_id != teacher_id_yuekebao || obj.Name != name{
//		_,err := s.EditCourse(name,"",id)
//		if err != nil{
//			log4go.Error("class in create course err",err)
//			log4go.Info("err course is:",name)
//			//return err
//			return err
//		}
//	}
//
//	if tmpteacher[teacher_id_yuekebao] != nil{
//		teacher_id_classin = tmpteacher[teacher_id_yuekebao].(model.Member).SID
//	}
//	foldeid,err := s.GetTopFolderId()
//	if err != nil{
//		log4go.Error("class in get FolderList err",err)
//		//log4go.Info("err course is:",v)
//		return err
//	}
//	folderId = strconv.Itoa(foldeid)
//	courseId := obj.Course_id
//	classId := obj.Class_id
//
//	if name != obj.Name || start_time != obj.Starttime || end_time != obj.Endtime || seatNum != obj.SeatNum,teacher_id_classin != obj.Teacher_sid {
//		_,err := s.EditCourseClass(courseId,classId,name,start_time,end_time,teacher_id_classin,folderId,
//			teachMode,isAutoOnstage,seatNum,isHd)
//		//result,err := s.AddTeacher(phone,name)
//		if err != nil{
//			log4go.Error("class in CreateCourseClass  err",err)
//			return err
//		}
//	}
//	if seatNum != obj.SeatNum{
//		_,err := s.ModifyClassSeatNum(courseId,classId,seatNum,isHd)
//		//result,err := s.AddTeacher(phone,name)
//		if err != nil{
//			log4go.Error("class in ModifyClassSeatNum  err",err)
//			return err
//		}
//	}
//
//	//课节添加学生
//	studentJson := make([]map[string]interface{},0)
//	studentids:=make([]string,0)
//	studentsids:=make([]string,0)
//	if len(orderArr) > 0{
//		for _,v:=range orderArr{
//			if obj,ok := v.(map[string]interface{});ok{
//				studentid := obj["member"].(string)
//				if tmpstudent[studentid]!=nil{
//					uid := tmpstudent[studentid].(model.Member).SID
//					studentJson = append(studentJson, map[string]interface{}{"uid":uid} )
//					studentsids = append(studentsids, uid)
//					studentids = append(studentids, studentid)
//				}
//			}
//		}
//	}
//	localstudentsids := strings.Split(obj.Student_sid,",")
//
//	difference_sids := difference(studentsids,localstudentsids)
//
//	if len(difference_sids) > 0{
//		STUDENT_JSON_DELETE := ""
//		studentJson_bytes_delete,err := json.Marshal(localstudentsids)
//		if err != nil{
//			log4go.Error("studentJson marshal err",err)
//			return err
//		}else{
//			STUDENT_JSON_DELETE = string(studentJson_bytes_delete)
//		}
//		_,err = s.DelClassStudentMultiple(courseId,classId,STUDENT_JSON_DELETE)
//		if err != nil{
//			log4go.Error("add class student multiple err",err)
//			return err
//		}
//		STUDENT_JSON := ""
//		studentJson_bytes,err := json.Marshal(studentJson)
//		if err != nil{
//			log4go.Error("studentJson marshal err",err)
//		}else{
//			STUDENT_JSON = string(studentJson_bytes)
//		}
//		addsuccess,err := s.AddClassStudentMultiple(courseId,classId,STUDENT_JSON)
//		if err != nil{
//			log4go.Error("add class student multiple err",err)
//		}
//		if addsuccess == 1{
//			log4go.Info("class add student success")
//		}
//	}
//
//	student_id :=Getids(studentids)
//	student_sid :=Getids(studentsids)
//	DB_id ,_ := strconv.Atoi(id)
//	DB_SID := courseId
//	updateCourse := model.Course{
//		Id: DB_id,
//		Name: name,
//		Course_id: DB_SID,
//		Class_id:classId,
//		Starttime:start_time,
//		Endtime:end_time,
//		Teacher_id:teacher_id_yuekebao,
//		Teacher_sid:teacher_id_classin,
//		Student_id:student_id,
//		Student_sid:student_sid,
//		SeatNum:seatNum,
//	}
//	err = s.UpdateCourse(updateCourse)
//	if err != nil{
//		log4go.Error("add student to local db err",err)
//		return err
//	}
//	return nil
//}

func Getids(idlist []string)(string){
	if len(idlist) == 0{
		return ""
	}
	return strings.Join(idlist,",")
}

func GetUnixTime(time_ string)(string,error){
	loc, _ := time.LoadLocation("Local") //获取时区
	//ASCtime, err := time.ParseInLocation("2006-01-02 15:04:05", recordTime["time_asc"].(string), loc)
	DESCtime, err := time.ParseInLocation("2006-01-02 15:04", time_, loc)
	if err != nil {
		log4go.Error(err)
		return "-100", errors.New("时间转换错误!")
	}
	//加1s以确保不会把相同的识别记录拉至本地，避免造成相同记录
	DESCt := DESCtime.Unix()

	s := strconv.FormatInt(DESCt, 10)
	return s,nil
}


//求并集
func union(slice1, slice2 []string) []string {
	m := make(map[string]int)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 0 {
			slice1 = append(slice1, v)
		}
	}
	return slice1
}

//求交集
func intersect(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	for _, v := range slice1 {
		m[v]++
	}

	for _, v := range slice2 {
		times, _ := m[v]
		if times == 1 {
			nn = append(nn, v)
		}
	}
	return nn
}

//求差集 slice1-并集
func difference(slice1, slice2 []string) []string {
	m := make(map[string]int)
	nn := make([]string, 0)
	inter := intersect(slice1, slice2)
	for _, v := range inter {
		m[v]++
	}

	for _, value := range slice1 {
		times, _ := m[value]
		if times == 0 {
			nn = append(nn, value)
		}
	}
	return nn
}