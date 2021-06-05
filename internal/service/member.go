package service

import (
	"github.com/alecthomas/log4go"
	"strings"
	"time"
	"zhiyuan/QBSED/internal/model"
	"zhiyuan/QBSED/internal/oss_server"
)

func(s *Service)GetStudentsFromDb()(members []model.Member,err error){

	members,err = s.dao.GetStudents()
	if err != nil{
		log4go.Error("get students from local db err:",err)
		return nil,err
	}
	return members,nil
}


func(s *Service)GetTeachersFromDb()(members []model.Member,err error){

	members,err = s.dao.GetTeachers()
	if err != nil{
		log4go.Error("get teachers from local db err:",err)
		return nil,err
	}
	return members,nil
}


func(s *Service)AddMember(record model.Member)(err error){

	err = s.dao.AddMembers(record)
	if err != nil{
		log4go.Error("add member to local db err:",err)
		return err
	}
	return nil
}
func(s *Service)AddCourse(record model.Course)(err error){

	err = s.dao.AddCourse(record)
	if err != nil{
		log4go.Error("add Course to local db err:",err)
		return err
	}
	return nil
}

func(s *Service)UpdateMember(record model.Member)(err error){

	err = s.dao.UpdateMembers(record)
	if err != nil{
		log4go.Error("update member to local db err:",err)
		return err
	}
	return nil
}
func(s *Service)DeleteCourse(record model.Course)(err error){

	err = s.dao.DeleteCourse(record)
	if err != nil{
		log4go.Error("Delete Course local db err:",err)
		return err
	}
	return nil
}

func(s *Service)UpdateCourse(record model.Course)(err error){

	err = s.dao.UpdateCourse(record)
	if err != nil{
		log4go.Error("update Course to local db err:",err)
		return err
	}
	return nil
}
func(s *Service)GetSignUrl(filepath string)(result map[string]string,err error){

	result,err = oss_server.GetSignUrl(filepath)
	if err != nil{
		log4go.Error("oss get sign url err:",err)
		return result,err
	}
	return result,nil
}

func(s *Service)GetCourses()(result []model.Course,err error){

	result,err = s.dao.GetCourse()
	if err != nil{
		log4go.Error("update Course to local db err:",err)
		return result,err
	}
	return result,nil
}
func(s *Service)GetCoursesbyid(id string)(result model.Course,err error){

	result,err = s.dao.GetCoursebyid()
	if err != nil{
		log4go.Error("update Course to local db err:",err)
		return result,err
	}
	return result,nil
}

func(s *Service)Register(userinfo model.User_json)(result model.User,err error){

	user := model.User{
		Name:userinfo.Name,
		Password:userinfo.Password,
		Phone:userinfo.Phone,
		Created_at:time.Now().Unix(),
		Role:1,
	}
	err = s.dao.CheckUser(user)
	if err != nil{
		if strings.Contains(string(err.Error()), "record not found") {
			log4go.Info("无重复记录")
		}else{
			log4go.Error("检查出错:",err.Error())
			return result, err
		}
	}
	err = s.dao.AddUser(user)
	if err != nil{
		log4go.Error("注册用户失败:",err)
		return result,err
	}
	return result,nil
}
