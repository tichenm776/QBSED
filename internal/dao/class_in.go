package dao

import (
	"errors"
	"github.com/alecthomas/log4go"
	"zhiyuan/QBSED/internal/model"
)

func(d *Dao) GetStudents()(result []model.Member,err error){

	DBdate := d.crmdb

	DBdate = DBdate.Where("member_type = ?", "1")

	if err := DBdate.Model(model.Member{}).Find(&result);err.Error!= nil {
		return result, err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return result, nil
}

func(d *Dao) GetTeachers()(result []model.Member,err error){

	DBdate := d.crmdb

	DBdate = DBdate.Where("member_type = ?", "2")

	if err := DBdate.Model(model.Member{}).Find(&result);err.Error!= nil {
		return result, err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return result, nil
}

func(d *Dao) AddMembers(record model.Member)(err error){

	if err := d.crmdb.Debug().Create(&record).Error; err != nil {
		log4go.Error(err.Error())
		return errors.New("新增历史记录失败!")
	}
	return nil
}

func (d *Dao) UpdateMembers(statistic_obj model.Member)(err error){


	if err := d.crmdb.Debug().Table("member").Where("id = ? ",statistic_obj.Id).Update(&statistic_obj).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("修改记录失败!")
	}
	if err := d.crmdb.Debug().Where("id = ? ",statistic_obj.Id).Last(&statistic_obj).Error;err != nil{
		log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	}
	return nil
}

func(d *Dao) GetCourses()(result []model.Member,err error){

	DBdate := d.crmdb

	DBdate = DBdate.Where("member_type = ?", "2")

	if err := DBdate.Model(model.Member{}).Find(&result);err.Error!= nil {
		return result, err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return result, nil
}



func(d *Dao) AddCourse(record model.Course)(err error){

	if err := d.crmdb.Debug().Create(&record).Error; err != nil {
		log4go.Error(err.Error())
		return errors.New("新增课程失败!")
	}
	return nil
}

func (d *Dao) UpdateCourse(statistic_obj model.Course)(err error){


	if err := d.crmdb.Debug().Table("course").Where("id = ? ",statistic_obj.Id).Update(&statistic_obj).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("修改记录失败!")
	}
	//if err := d.crmdb.Debug().Where("id = ? ",statistic_obj.Id).Last(&statistic_obj).Error;err != nil{
	//	log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	//}
	return nil
}

func(d *Dao) GetCourse()(result []model.Course,err error){

	DBdate := d.crmdb

	//DBdate = DBdate.Where("member_type = ?", "1")

	if err := DBdate.Model(model.Course{}).Find(&result);err.Error!= nil {
		return result, err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return result, nil
}
func(d *Dao) GetCoursebyid()(result model.Course,err error){

	DBdate := d.crmdb

	//DBdate = DBdate.Where("member_type = ?", "1")

	if err := DBdate.Model(model.Course{}).Last(&result);err.Error!= nil {
		return result, err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return result, nil
}
func(d *Dao) DeleteCourse(statistic_obj model.Course)(err error){

	DBdate := d.crmdb

	//DBdate = DBdate.Where("member_type = ?", "1")

	if err := DBdate.Model(model.Course{}).Where("id = ?",statistic_obj.Id).Delete(&model.Course{});err.Error!= nil {
		return  err.Error
	}

	//total = zyutil.GetTotal(count, Records.Size)
	return  nil
}


func(d *Dao) AddUser(record model.User)(err error){

	if err := d.crmdb.Debug().Create(&record).Error; err != nil {
		log4go.Error(err.Error())
		return errors.New("新增课程失败!")
	}
	return nil
}
func(d *Dao) CheckUser(record model.User)(err error){

	if err := d.crmdb.Debug().Where("phone = ?",record.Phone).Last(&record).Error; err != nil {
		log4go.Error(err.Error())
		return err
	}
	return nil
}
