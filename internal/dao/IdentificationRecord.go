package dao

import (
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/pkg/errors"
	"strconv"
	"time"
	"zhiyuan/QBSED/internal/model"
	"zhiyuan/zyutil"
)

func(d *Dao) GetIdentificationRecords(Records model.EmployeeRecords_json)(result []model.IdentificationRecord,err error,count int,total int){

	DBdate := d.crmdb
	DBcount := d.crmdb

	if Records.Subject_id != 0 {
		DBdate = DBdate.Where("subject_id = ?", Records.Subject_id)
		DBcount = DBcount.Where("subject_id = ?", Records.Subject_id)
	}

	if Records.Screen_id > 0 {
		DBdate = DBdate.Where("screen_id = ?", Records.Screen_id)
		DBcount = DBcount.Where("screen_id = ?", Records.Screen_id)
	}

	if Records.Name != "" {
		DBdate = DBdate.Where("name like ?", "%" + Records.Name +"%")
		DBcount = DBcount.Where("name like ?","%" + Records.Name +"%")
	}
	if Records.Snap_position != "" {
		//DBdate = DBdate.Where("snap_position = ?",  Records.Snap_position )
		//DBcount = DBcount.Where("snap_position = ?", Records.Snap_position)
	}
	if Records.User_role >= 0 {
		DBdate = DBdate.Where("subject_type = ?", Records.User_role)
		DBcount = DBcount.Where("subject_type = ?", Records.User_role)
	}

	if Records.Snap_begin_time != "" {
		DBdate = DBdate.Where("snap_time >= ?", Records.Snap_begin_time)
		DBcount = DBcount.Where("snap_time >= ?", Records.Snap_begin_time)
	}
	if Records.Snap_end_time != "" {
		DBdate = DBdate.Where("snap_time <= ?", Records.Snap_end_time)
		DBcount = DBcount.Where("snap_time <= ?", Records.Snap_end_time)
	}

	DBdate = DBdate.Order("snap_time desc")
	if Records.Page > 0 {
		DBdate = DBdate.Limit(Records.Size).Offset((Records.Page - 1) * Records.Size)
	}
	if err := DBdate.Model(model.IdentificationRecord{}).Find(&result);err.Error!= nil {
		return result, err.Error, 0,0
	}

	if err := DBcount.Model(model.IdentificationRecord{}).Count(&count); err.Error != nil {
		return result, err.Error, 0,0
	}
	total = zyutil.GetTotal(count, Records.Size)
	return result, nil, count, total
}

func(d *Dao) GetStrangerRecords(Records model.StrangerRecords_json)(result []model.Stranger,err error,count int,total int){

	DBdate := d.crmdb
	DBcount := d.crmdb
	if Records.Screen_id > 0 {
		DBdate = DBdate.Where("screen_id = ?", Records.Screen_id)
		DBcount = DBcount.Where("screen_id = ?", Records.Screen_id)
	}

	if Records.Snap_begin_time != "" {
		DBdate = DBdate.Where("snap_time >= ?", Records.Snap_begin_time)
		DBcount = DBcount.Where("snap_time >= ?", Records.Snap_begin_time)
	}
	if Records.Snap_end_time != "" {
		DBdate = DBdate.Where("snap_time <= ?", Records.Snap_end_time)
		DBcount = DBcount.Where("snap_time <= ?", Records.Snap_end_time)
	}

	DBdate = DBdate.Order("snap_time desc")
	if Records.Page > 0 {
		DBdate = DBdate.Limit(Records.Size).Offset((Records.Page - 1) * Records.Size)
	}
	if err := DBdate.Model(model.Stranger{}).Find(&result);err.Error!= nil {
		return result, err.Error, 0,0
	}

	if err := DBcount.Model(model.Stranger{}).Count(&count); err.Error != nil {
		return result, err.Error, 0,0
	}
	total = zyutil.GetTotal(count, Records.Size)
	return result, nil, count, total
}

const (
	timeasc = "SELECT * FROM `identification_record`  ORDER BY snap_time ASC LIMIT 1 "
	timedesc = "SELECT * FROM `identification_record`  ORDER BY snap_time DESC LIMIT 1 "
)

//查询历史记录最早与最晚的的时间
func (d *Dao) GetTime() (map[string]interface{}, error) {
	var recordASC model.IdentificationRecord
	var recordDESC model.IdentificationRecord
	var record = make(map[string]interface{})
	//if err := Db.Debug().Order("checktime asc").First(&recordASC).Error; err != nil {
	if err := d.crmdb.Debug().Raw(timeasc).Scan(&recordASC).Error; err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("查询识别记录升序失败!")
	}
	//if err := Db.Debug().Order("checktime desc").First(&recordDESC).Error; err != nil {
	if err := d.crmdb.Debug().Raw(timedesc).Scan(&recordDESC).Error; err != nil {
		log4go.Error(err.Error())
		return nil, errors.New("查询识别记录升序失败!")
	}
	record["time_asc"] = recordASC.Snap_time
	record["time_desc"] = recordDESC.Snap_time
	return record, nil
}

func (d *Dao) GetRecords() (bool, error) {
	var recordASC model.IdentificationRecord
	var record = make(map[string]interface{})
	//if err := Db.Debug().Order("checktime asc").First(&recordASC).Error; err != nil {
	if err := d.crmdb.Debug().Raw(timeasc).Scan(&recordASC).Error; err != nil {
		log4go.Error(err.Error())
		return false, errors.New("record not found")
	}
	//if err := Db.Debug().Order("checktime desc").First(&recordDESC).Error; err != nil {
	record["time_asc"] = recordASC.Snap_time

	return true, nil
}

func(d *Dao) AddIdentificationRecord(record model.IdentificationRecord) (error) {

	if err := d.crmdb.Debug().Create(&record).Error; err != nil {
		log4go.Error(err.Error())
		return errors.New("新增历史记录失败!")
	}
	return nil
}

func(d *Dao) AddStrangerRecord(record model.Stranger) (error) {

	if err := d.crmdb.Debug().Create(&record).Error; err != nil {
		log4go.Error(err.Error())
		return errors.New("新增历史记录失败!")
	}
	return nil
}

func(d *Dao) GetIdentificationRecords_GroupBy(Records model.EmployeeRecords_json)(result []model.IdentificationRecord,err error,count int,total int,headcountmap map[int]int){

	DBdate := d.crmdb
	DBcount := d.crmdb

	Headcount := d.crmdb

	if Records.Subject_id != 0 {
		DBdate = DBdate.Where("subject_id = ?", Records.Subject_id)
		DBcount = DBcount.Where("subject_id = ?", Records.Subject_id)
		Headcount = Headcount.Where("subject_id = ?", Records.Subject_id)
	}

	if Records.Screen_id > 0 {
		DBdate = DBdate.Where("screen_id = ?", Records.Screen_id)
		DBcount = DBcount.Where("screen_id = ?", Records.Screen_id)
		Headcount = Headcount.Where("screen_id = ?", Records.Screen_id)
	}
	if Records.Snap_position != "" {
		DBdate = DBdate.Where("snap_position = ?", Records.Snap_position)
		DBcount = DBcount.Where("snap_position = ?", Records.Snap_position)
		Headcount = Headcount.Where("snap_position = ?", Records.Snap_position)
	}
	if Records.Name != "" {
		DBdate = DBdate.Where("name like ?", "%" + Records.Name +"%")
		DBcount = DBcount.Where("name like ?","%" + Records.Name +"%")
		Headcount = Headcount.Where("name like ?","%" + Records.Name +"%")
	}

	if Records.User_role >= 0 {
		DBdate = DBdate.Where("subject_type = ?", Records.User_role)
		DBcount = DBcount.Where("subject_type = ?", Records.User_role)
		Headcount = Headcount.Where("subject_type = ?", Records.User_role)
	}

	if Records.Snap_begin_time != "" {
		DBdate = DBdate.Where("snap_time >= ?", Records.Snap_begin_time)
		DBcount = DBcount.Where("snap_time >= ?", Records.Snap_begin_time)
		Headcount = Headcount.Where("snap_time >= ?", Records.Snap_begin_time)
	}
	if Records.Snap_end_time != "" {
		DBdate = DBdate.Where("snap_time <= ?", Records.Snap_end_time)
		DBcount = DBcount.Where("snap_time <= ?", Records.Snap_end_time)
		Headcount = Headcount.Where("snap_time <= ?", Records.Snap_end_time)
	}

	DBdate = DBdate.Order("snap_time desc")
	if Records.Page > 0 {
		DBdate = DBdate.Limit(Records.Size).Offset((Records.Page - 1) * Records.Size)
	}
	if err := DBdate.Model(model.IdentificationRecord{}).Group("subject_id").Find(&result);err.Error!= nil {
		return result, err.Error, 0,0,nil
	}

	if err := DBcount.Model(model.IdentificationRecord{}).Group("subject_id").Count(&count); err.Error != nil {
		return result, err.Error, 0,0,nil
	}
	headcountmaps := make(map[int]int)
	var headcount int
	//subjectids := make([]int,0)
	for _,v:= range result{
		//subjectids = append(subjectids,v.Subject_id)
		if err := Headcount.Model(model.IdentificationRecord{}).Where("subject_id = ?",v.Subject_id).Count(&headcount); err.Error != nil {
			//return result, err.Error, 0,0,nil
		}
		headcountmaps[v.Subject_id] = headcount
	}

	total = zyutil.GetTotal(count, Records.Size)
	return result, nil, count, total,headcountmaps
}

func (d *Dao) CreateStatistic(statistic model.Statistic)(error){

	sta:= model.Statistic{}
	if err := d.crmdb.Debug().Where("date = ? ",statistic.Date).Last(&sta).Error;err != nil{
		log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	}

	if sta.Date ==""{
		if err := d.crmdb.Debug().Create(&statistic).Error;err != nil{
			log4go.Error(err.Error())
			return errors.New("新增记录失败!")
		}
	}else {
		if err := d.crmdb.Debug().Table("statistic").Where("date = ?",statistic.Date).Updates(&statistic).Error;err != nil{
			log4go.Error(err.Error())
			return errors.New("更新记录失败!")
		}
	}
	return nil
}

func (d *Dao) GetStatistic(date string)(model.Statistic,error){

	var sta model.Statistic
	if err := d.crmdb.Debug().Where("date = ? ",date).Last(&sta).Error;err != nil{
		log4go.Error(err.Error())
		return sta,errors.New("查询记录失败!")
	}
	return sta,nil
}
func (d *Dao) GetStatistic2(date string,grandparent_id,parent_id int)(model.Statistics,error){

	var sta model.Statistics
	if err := d.crmdb.Debug().Where("grandparent_id = ? and  parent_id = ? and date = ? ",grandparent_id,parent_id,date).Last(&sta).Error;err != nil{

		log4go.Error(err.Error())
		return sta,errors.New("查询记录失败!")
	}
	return sta,nil
}


func (d *Dao) Statistic_Project_Create(name string)(statistic_obj model.Statistic_Main,err error){

	sta:= model.Statistic_Main{
		Name:name,
	}
	if err := d.crmdb.Debug().Create(&sta).Error;err != nil{
		log4go.Error(err.Error())
		return statistic_obj,errors.New("新增记录失败!")
	}
	if err := d.crmdb.Debug().Where("name = ? ",name).Last(&statistic_obj).Error;err != nil{
		log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	}
	return statistic_obj,nil
}

func (d *Dao) Statistic_Project_Item_Create(SPmain_id int,Statistic_Item_Json model.Statistic_Item_Json)(err error){



	sta_item:= model.Statistic_Item{
		Parent_id:SPmain_id,
		Name:Statistic_Item_Json.Name,
		Subjects_group:Statistic_Item_Json.Subjects_group,
		//Screens_group:Statistic_Item_Json.Screens_group,

		Start_time:Statistic_Item_Json.Start_time,
		End_time:Statistic_Item_Json.End_time,
	}
	groups := ""
	fmt.Println("Statistic_Item_Json.Screens_groups is",Statistic_Item_Json.Screens_groups)

	if len(Statistic_Item_Json.Screens_groups) > 1{
		for i := 0 ;i<len(Statistic_Item_Json.Screens_groups);i++{
			if i < len(Statistic_Item_Json.Screens_groups)-1{
				id := strconv.Itoa(Statistic_Item_Json.Screens_groups[i])
				temp := id + ","
				groups =  groups + temp
			}
			if i == len(Statistic_Item_Json.Screens_groups)-1{
				id := strconv.Itoa(Statistic_Item_Json.Screens_groups[i])
				temp := id
				groups =  groups + temp
			}
		}
	}
	if len(Statistic_Item_Json.Screens_groups) == 1{
		fmt.Println("ids is ----------------------------------------- 1")
		fmt.Println(Statistic_Item_Json.Screens_groups[0])

			id := strconv.Itoa(Statistic_Item_Json.Screens_groups[0])
			temp := id
			groups =  groups + temp
		fmt.Println(id)
		fmt.Println(groups)
	}
	if len(Statistic_Item_Json.Screens_groups) == 0{
		groups =  "-1"
	}
	sta_item.Screens_groups = groups

	if err := d.crmdb.Debug().Create(&sta_item).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("新增记录失败!")
	}

	return nil
}

func (d *Dao) Statistic_Project_Update(name string,SPmain_id int)(statistic_obj model.Statistic_Main,err error){

	sta:= model.Statistic_Main{
		Name: name,
	}
	if err := d.crmdb.Debug().Table("statistic_main").Where("id = ? ",SPmain_id).Update(&sta).Error;err != nil{
		log4go.Error(err.Error())
		return statistic_obj,errors.New("修改记录失败!")
	}
	if err := d.crmdb.Debug().Where("id = ? ",SPmain_id).Last(&statistic_obj).Error;err != nil{
		log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	}
	return statistic_obj,nil
}

func (d *Dao) Statistic_Project_Items_Get(SPmain_id int)(Statistic_Items []model.Statistic_Item , err error){

	if err := d.crmdb.Debug().Where("parent_id = ?",SPmain_id).Find(&Statistic_Items).Error;err != nil{
		log4go.Error(err.Error())
		return Statistic_Items,errors.New("查询子项失败!")
	}

	return Statistic_Items,nil
}

func (d *Dao) Statistic_Project_Items_Delete(SPmain_id int)(err error){

	if err := d.crmdb.Debug().Where("parent_id = ?",SPmain_id).Delete(&model.Statistic_Item{}).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("删除子项失败!")
	}

	return nil
}
func (d *Dao) Statistics_Project_Delete(id int)(err error){

	if err := d.crmdb.Debug().Where("parent_id = ?",id).Delete(&model.Statistics{}).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("删除历史统计项目失败!")
	}

	return nil
}
func (d *Dao) Statistic_Project_Delete(id int)(err error){

	if err := d.crmdb.Debug().Where("id = ?",id).Delete(&model.Statistic_Main{}).Error;err != nil{
		log4go.Error(err.Error())
		return errors.New("删除统计项目失败!")
	}

	return nil
}

func (d *Dao) Statistic_Project_Get()([]model.Statistic_Main,error){
	var (
		Statistic_Main []model.Statistic_Main
		err error
		)
	if err = d.crmdb.Debug().Find(&Statistic_Main).Error;err != nil{
		log4go.Error(err.Error())
		return Statistic_Main,errors.New("查询子项失败!")
	}

	return Statistic_Main,nil
}


//func (d *Dao) GetOneItemStatistic(screens,subjects []int,subject_type int,start_time,end_time string)([]model.IdentificationRecord,error){
//func (d *Dao) GetOneItemStatistic(screens,subjects []int,subject_type int,start_time,end_time string)(){
func (d *Dao) GetOneItemStatistic(screens []int,subject,subject_type int,start_time,end_time string)([]model.IdentificationRecord , error){
	//部分相机 部分人员
	var(
		//screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql string
		screen_sql,screenid_sql,time_sql,finalsql,sql string
		subject_objs []model.IdentificationRecord
	)

	finalsql = "select * from identification_record where subject_id = %d and "
	screenid_sql = "screen_id="
	if len(screens) == 0{
		return nil,nil
	}
	for k,v:=range screens{
		if k < len(screens)-1{
			screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
		}else{
			screen_sql += screenid_sql+strconv.Itoa(v)+" "
		}
	}
	screen_sql = "( " + screen_sql + " ) order by snap_time desc;"

	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
	if subject_type != 0{
		//type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		//sql = finalsql + type_sql + time_sql + screen_sql
		sql = finalsql  + time_sql + screen_sql
		}else {
		sql = finalsql + time_sql +screen_sql
	}

	sql = fmt.Sprintf(sql,subject)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,errors.New("查询数据库出错")
	}
	//fmt.Println(sql)
	//fmt.Println(subject_objs)
	//fmt.Println(screen_sql)
	return subject_objs,nil
}
func (d *Dao) GetOneItemStatistic_copy(screens []int,subject []int,subject_type int,start_time,end_time string)([]model.IdentificationRecord , error){
	//部分相机 部分人员
	var(
		//screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql string
		screen_sql,screenid_sql,subject_sql,subject_id_sql,type_sql,time_sql,finalsql,sql string
		subject_objs []model.IdentificationRecord
	)
	if len(subject) != 0{
		finalsql = "select * from identification_record where  ("
		subject_id_sql = " subject_id = "
		for i,c:=range subject{
			if i < len(subject)-1{
				subject_sql += subject_id_sql+strconv.Itoa(c)+" OR "
			}else{
				subject_sql += subject_id_sql+strconv.Itoa(c)+" ) and"
			}
		}
	}else{
		finalsql = "select * from identification_record where subject_id > 0 and "

		subject_id_sql = ""
	}
	screenid_sql = "screen_id="
	if len(screens) == 0{
		return nil,nil
	}
	for k,v:=range screens{
		if k < len(screens)-1{
			screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
		}else{
			screen_sql += screenid_sql+strconv.Itoa(v)+" "
		}
	}
	screen_sql = "( " + screen_sql + " ) group by subject_id order by snap_time desc;"

	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
	if subject_type != 0{
		type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		//sql = finalsql + type_sql + time_sql + screen_sql
		sql = finalsql  +subject_sql+ type_sql + time_sql + screen_sql
	}else {
		sql = finalsql +subject_sql+ time_sql +screen_sql
	}
	fmt.Println(sql)
	fmt.Println("--------------------------------------------------------------这里")
	//sql = fmt.Sprintf(sql,subject)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,errors.New("查询数据库出错")
	}
	//fmt.Println(sql)
	//fmt.Println(subject_objs)
	//fmt.Println(screen_sql)
	return subject_objs,nil
}
func (d *Dao) GetOneItemStatistic2(screens []int,subject,subject_type int,start_time,end_time string)([]model.IdentificationRecord , error){
	//部分相机 全人员
	var(
		//screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql string
		screen_sql,screenid_sql,time_sql,finalsql,sql string
		subject_objs []model.IdentificationRecord
	)

	finalsql = "select * from identification_record where   "
	screenid_sql = "screen_id="
	if len(screens) == 0{
		return nil,nil
	}
	for k,v:=range screens{
		if k < len(screens)-1{
			screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
		}else{
			screen_sql += screenid_sql+strconv.Itoa(v)+" "
		}
	}
	screen_sql = "( " + screen_sql + " ) order by snap_time desc;"

	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
	if subject_type != 0{
		//type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		//sql = finalsql + type_sql + time_sql + screen_sql
		sql = finalsql  + time_sql + screen_sql
	}else {
		sql = finalsql + time_sql +screen_sql
	}

	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,errors.New("查询数据库出错")
	}
	//fmt.Println(sql)
	//fmt.Println(subject_objs)
	//fmt.Println(screen_sql)
	return subject_objs,nil
}
func (d *Dao) GetOneItemStatistic3(subject,subject_type int,start_time,end_time string)([]model.IdentificationRecord , error){
	//全相机 部分人员
	var(
		//screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql string
		screen_sql,time_sql,finalsql,sql string
		subject_objs []model.IdentificationRecord
	)

	finalsql = "select * from identification_record where subject_id = %d and "

	screen_sql =  "  order by snap_time desc;"

	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s'  ", start_time,end_time)
	if subject_type != 0{
		//type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		//sql = finalsql + type_sql + time_sql + screen_sql
		sql = finalsql  + time_sql + screen_sql
	}else {
		sql = finalsql + time_sql +screen_sql
	}

	sql = fmt.Sprintf(sql,subject)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,errors.New("查询数据库出错")
	}
	//fmt.Println(sql)
	//fmt.Println(subject_objs)
	//fmt.Println(screen_sql)
	return subject_objs,nil
}
func (d *Dao) GetOneItemStatistic4(start_time,end_time string)([]model.IdentificationRecord , error){
	//全相机 全人员
	var(
		//screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql string
		screen_sql,time_sql,finalsql,sql string
		subject_objs []model.IdentificationRecord
	)
	finalsql = "select * from identification_record where "
	screen_sql = "  order by snap_time desc;"
	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' ", start_time,end_time)
	sql = finalsql  + time_sql + screen_sql
	fmt.Println(sql)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,errors.New("查询数据库出错")
	}
	//fmt.Println(sql)
	//fmt.Println(subject_objs)
	//fmt.Println(screen_sql)
	return subject_objs,nil
}

func (d *Dao) GetOneItemStatisticRecords(screens []int,subject,subject_type,page,size int,start_time,end_time string)([]model.IdentificationRecord ,int,int, error){

	var(
		screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql,countsql2 string
		subject_objs,count_objs []model.IdentificationRecord
		count int
	)
	if subject < 0{
		finalsql = "select * from identification_record where subject_id > 0   and "
		countsql2 = "select * from identification_record where  subject_id > 0  and "
		screenid_sql = "screen_id="
		if len(screens) == 0{
			return nil,0,0,nil
		}
		for k,v:=range screens{
			if k < len(screens)-1{
				screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
			}else{
				screen_sql += screenid_sql+strconv.Itoa(v)+" "
			}
		}
		screen_sql_copy := "( " +screen_sql+ " )"
		screen_sql = "( " + screen_sql + " ) group by subject_id order by snap_time desc LIMIT %d OFFSET %d;"
		count_sql :=  "group by subject_id order by snap_time desc;"
		screen_sql = fmt.Sprintf(screen_sql,size,(page-1)*size)
		time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
		countsql3 := countsql2 + time_sql + screen_sql_copy +count_sql
		if subject_type != 0{
			type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
			sql = finalsql + type_sql + time_sql + screen_sql
		}else {
			sql = finalsql + time_sql +screen_sql
		}

		//sql = fmt.Sprintf(sql,subject)
		//countsql3 = fmt.Sprintf(countsql3,subject)
		if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
			log4go.Error("查询数据库出错:",err)
			//continue
			return nil,0,0,errors.New("查询数据库出错")
		}
		log4go.Info("countsql3",countsql3)
		if err:=d.crmdb.Debug().Raw(countsql3).Scan(&count_objs).Error;err!=nil{
			log4go.Error("查询数据库出错:",err)
			//continue
			return nil,0,0,errors.New("查询数据库出错")
		}
		count = len(count_objs)
		total := zyutil.GetTotal(count, size)
		return subject_objs,count,total,nil
	}
	finalsql = "select * from identification_record where subject_id = %d and "
	countsql2 = "select * from identification_record where subject_id = %d and "
	screenid_sql = "screen_id="
	if len(screens) == 0{
		return nil,0,0,nil
	}
	for k,v:=range screens{
		if k < len(screens)-1{
			screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
		}else{
			screen_sql += screenid_sql+strconv.Itoa(v)+" "
		}
	}
	screen_sql_copy := "( " +screen_sql+ " )"
	screen_sql = "( " + screen_sql + " ) order by snap_time desc LIMIT %d OFFSET %d;"
	count_sql :=  "order by snap_time desc;"
	screen_sql = fmt.Sprintf(screen_sql,size,(page-1)*size)
	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
	countsql3 := countsql2 + time_sql + screen_sql_copy +count_sql
	if subject_type != 0{
		type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		sql = finalsql + type_sql + time_sql + screen_sql
	}else {
		sql = finalsql + time_sql +screen_sql
	}

	sql = fmt.Sprintf(sql,subject)
	countsql3 = fmt.Sprintf(countsql3,subject)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)

		//continue
		return nil,0,0,errors.New("查询数据库出错:"+err.Error())
	}
	//log4go.Info("countsql3",countsql3)
	if err:=d.crmdb.Debug().Raw(countsql3).Scan(&count_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,0,0,errors.New("查询数据库出错")
	}
	count = len(count_objs)
	total := zyutil.GetTotal(count, size)
	return subject_objs,count,total,nil
}


func (d *Dao) GetOneItemStatisticRecord(screens []int,subject,subject_type,page,size int,start_time,end_time string)([]model.IdentificationRecord ,int,int, error){

	var(
		screen_sql,screenid_sql,time_sql,finalsql,type_sql,sql,countsql2 string
		subject_objs,count_objs []model.IdentificationRecord
		count int
	)

	finalsql = "select * from identification_record where subject_id = %d and "
	countsql2 = "select * from identification_record where subject_id = %d and "
	screenid_sql = "screen_id="
	for k,v:=range screens{
		if k < len(screens)-1{
			screen_sql += screenid_sql+strconv.Itoa(v)+" OR "
		}else{
			screen_sql += screenid_sql+strconv.Itoa(v)+" "
		}
	}
	screen_sql_copy := "( " +screen_sql+ " )"
	screen_sql = "( " + screen_sql + " ) order by snap_time desc LIMIT %d OFFSET %d;"
	count_sql :=  "order by snap_time desc;"
	screen_sql = fmt.Sprintf(screen_sql,size,(page-1)*size)
	time_sql = fmt.Sprintf("  snap_time BETWEEN '%s' and '%s' and ", start_time,end_time)
	countsql3 := countsql2 + time_sql + screen_sql_copy +count_sql
	if subject_type != 0{
		type_sql = fmt.Sprintf("  subject_type = %d and ",subject_type)
		sql = finalsql + type_sql + time_sql + screen_sql
	}else {
		sql = finalsql + time_sql +screen_sql
	}

	sql = fmt.Sprintf(sql,subject)
	countsql3 = fmt.Sprintf(countsql3,subject)
	if err:=d.crmdb.Debug().Raw(sql).Scan(&subject_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,0,0,errors.New("查询数据库出错")
	}
	fmt.Println(countsql3)
	if err:=d.crmdb.Debug().Raw(countsql3).Scan(&count_objs).Error;err!=nil{
		log4go.Error("查询数据库出错:",err)
		//continue
		return nil,0,0,errors.New("查询数据库出错")
	}
	count = len(count_objs)
	total := zyutil.GetTotal(count, size)
	return subject_objs,count,total,nil
}



func(d *Dao)UpdateOneItemStatistic(s model.Statistics)(error){
	sta:= model.Statistics{}
	if err := d.crmdb.Debug().Where("grandparent_id = ? " +
		"and parent_id = ? and date = ? ",s.Grandparent_id,s.Parent_id,s.Date).Last(&sta).Error;err != nil{
		log4go.Error(err.Error())
		//return errors.New("查询记录失败!")
	}
	if sta.Date ==""{
		//if err := d.crmdb.Debug().Create(&s).Error;err != nil{
		//	log4go.Error(err.Error())
		//	return errors.New("新增记录失败!")
		//}
		go func() {
			//d.crmdb.Debug().Table("statistics").Where("grandparent_id = ? " +
			//	"and parent_id = ? and date = ? ",s.Grandparent_id,s.Parent_id,s.Date).Update(&s)
			if err := d.crmdb.Debug().Create(&s).Error;err != nil{
				log4go.Error(err.Error())
			}
		}()
	}else {
		if err := d.crmdb.Debug().Table("statistics").Where("grandparent_id = ? " +
			"and parent_id = ? and date = ? ",s.Grandparent_id,s.Parent_id,s.Date).Update(&s).Error;err != nil{
			log4go.Error(err.Error())
			return errors.New("更新记录失败!")
		}
		//go func() {
		//	//d.crmdb.Debug().Table("statistics").Where("grandparent_id = ? " +
		//	//	"and parent_id = ? and date = ? ",s.Grandparent_id,s.Parent_id,s.Date).Update(&s)
		//	if err := d.crmdb.Debug().Table("statistics").Where("grandparent_id = ? " +
		//		"and parent_id = ? and date = ? ",s.Grandparent_id,s.Parent_id,s.Date).Update(&s).Error;err != nil{
		//		log4go.Error(err.Error())
		//	}
		//}()
	}
	return nil
}

func(d *Dao) GetIdentificationRecordById(screen_id,subject_id int)(result model.IdentificationRecord,err error){

	DBdate := d.crmdb
	if err := DBdate.Model(model.IdentificationRecord{}).Where("subject_id = ? and screen_id = ?",subject_id,screen_id).Last(&result);err.Error!= nil {
		return result, err.Error
	}
	return result, nil
}

func(d *Dao) TruncateTable()(err error){

	if err := d.crmdb.Debug().Exec("truncate table identification_record").Error;err != nil{
		log4go.Error("truncate table records err",err)
		return err
	}
	if err := d.crmdb.Debug().Exec("truncate table stranger").Error;err != nil{
		log4go.Error("truncate stranger err",err)
		return err
	}
	if err := d.crmdb.Debug().Exec("truncate table statistics").Error;err != nil{
		log4go.Error("truncate stranger err",err)
		return err
	}
	return nil

}


func(d *Dao) DeleteRecord(day int) ( error) {
	timeLayout := "2006-01-02 15:04:05"
	nowTime := time.Now()
	days := day
	startTime := nowTime.AddDate(0, 0, -days).Unix()
	dataTimeStr := time.Unix(startTime, 0).Format(timeLayout)
	if err := d.crmdb.Debug().Where("snap_time < ?",dataTimeStr).Delete(model.IdentificationRecord{}).Error;err != nil{
		log4go.Error(err.Error())
		return  errors.New("查询识别记录条数失败!")
	}
	if err := d.crmdb.Debug().Where("snap_time < ?",dataTimeStr).Delete(model.Stranger{}).Error;err != nil{
		log4go.Error(err.Error())
		return  errors.New("查询识别记录条数失败!")
	}
	return  nil

}