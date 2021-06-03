package service

import (
	"github.com/alecthomas/log4go"
	"zhiyuan/QBSED/internal/model"
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