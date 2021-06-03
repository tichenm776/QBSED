package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
	"net"
	"strconv"
	"strings"
	"time"
	//"zhiyuan/koala_api_go/koala_api"
	"zhiyuan/classin/configs"
	"zhiyuan/classin/internal/koala"
	"zhiyuan/classin/internal/model"
	"zhiyuan/zyutil"
)

func (s *Service) LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}

func (s *Service) GetIdentificationRecords(Records model.EmployeeRecords_json,position_map map[int]interface{})(results []model.IdentificationRecord,err error,count int,total int){

	results,err,count,total = s.dao.GetIdentificationRecords(Records)
	if err!=nil{
		return	nil,err,0,0
	}
	for k,v := range results{
		if v.Photo != ""{
			tmp_photo := "http://" + configs.Gconf.KoalaHost +v.Photo
			results[k].Photo = tmp_photo
		}
		tmp_Snap_photo:= "http://" + configs.Gconf.KoalaHost +v.Snap_photo
		results[k].Snap_photo = tmp_Snap_photo
		if v.Snap_position == ""{
			if position_map[v.Screen_id] != nil{
				tmp_Snap_position := position_map[v.Screen_id].(string)
				results[k].Snap_position = tmp_Snap_position
			}
		}
	}
	return results,nil,count,total
}
func (s *Service) GetIdentificationRecords_GroupBy(Records model.EmployeeRecords_json,position_map map[int]interface{})(results []model.IdentificationRecord,err error,count int,total int){

	results,err,count,total,headcountmap := s.dao.GetIdentificationRecords_GroupBy(Records)
	if err!=nil{
		return	nil,err,0,0
	}
	for k,v := range results{
		tmp_photo := "http://" + configs.Gconf.KoalaHost +v.Photo
		results[k].Photo = tmp_photo
		tmp_Snap_photo:= "http://" + configs.Gconf.KoalaHost +v.Snap_photo
		results[k].Snap_photo = tmp_Snap_photo
		if v.Snap_position == ""{
			tmp_Snap_position := position_map[v.Screen_id].(string)
			results[k].Snap_position = tmp_Snap_position
		}
		results[k].Recognition_times = headcountmap[v.Subject_id]
	}
	return results,nil,count,total
}
func (s *Service) GetStrangerRecords(Records model.StrangerRecords_json,position_map map[int]interface{})(results []model.Stranger,err error,count int,total int){

	results,err,count,total = s.dao.GetStrangerRecords(Records)
	if err!=nil{
		return	nil,err,0,0
	}
	for k,v := range results{
		tmp_Snap_photo:= "http://" + configs.Gconf.KoalaHost +v.Snap_photo
		results[k].Snap_photo = tmp_Snap_photo
		if v.Snap_position == ""{
			if position_map[v.Screen_id] != nil {
				tmp_Snap_position := position_map[v.Screen_id].(string)
				results[k].Snap_position = tmp_Snap_position
			}else{

			}
		}
	}
	return results,nil,count,total
}
func (s *Service)PrepareForRecord2() () {

	_, err := s.dao.GetRecords()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func (s *Service)PrepareForRecord(days int) (int, error) {


	//flag,err := koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	//if !flag || err != nil{
	//	return 0, nil
	//}
	koala.Init(configs.Gconf.KoalaHost)
	flag,err := koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	if !flag || err != nil{
		return 0, nil
	}

	//拉取koala中的历史记录
	count, err := s.dao.GetRecords()
	//fmt.Println(count)
	if err != nil {
		if strings.Contains(string(err.Error()), "record not found") {
			//return false, nil
		} else {
			log4go.Error(err.Error())
			return -100, err
		}
	}
	//要拉取的最早时间
	nowTime := time.Now()
	endTime := nowTime.Unix()
	//startTime := nowTime.AddDate(0, 0, -days).Unix()
	//fmt.Println(count)
	if count  {
		recordTime, err := s.dao.GetTime()
		if err != nil {
			return -100, err
		}
		loc, _ := time.LoadLocation("Local") //获取时区
		//ASCtime, err := time.ParseInLocation("2006-01-02 15:04:05", recordTime["time_asc"].(string), loc)
		DESCtime, err := time.ParseInLocation("2006-01-02 15:04:05", recordTime["time_desc"].(string), loc)
		if err != nil {
			log4go.Error(err)
			return -100, errors.New("时间转换错误!")
		}
		//加1s以确保不会把相同的识别记录拉至本地，避免造成相同记录
		DESCt := DESCtime.Unix()
		//fmt.Println(DESCt)
		//coed, err := s.GetRecord(DESCt, endTime)
		//if err != nil {
		//	return coed, err
		//}
		//ASCt := ASCtime.Unix()
		startTime := DESCt + 1
		//startTime := ASCt - 1
		//if startTime < ASCt {
			//减1s以确保不会把相同的识别记录拉至本地，避免造成相同记录
			//endTime = ASCt - 1
		code, err := s.GetRecord(startTime, endTime)
		if err != nil {
			return code, err
		}
		//}else{
		//	endTime = ASCt - 1
		//	code, err := s.GetRecord(endTime,startTime)
		//	if err != nil {
		//		return code, err
		//	}
		//}
	} else {
		days := 10

		beginTime := nowTime.AddDate(0, 0, -days).Unix()
		fmt.Println(days)
		code, err := s.GetRecord(beginTime, endTime)
		if err != nil {
			return code, err
		}
	}
	return 0, nil
}

func(s *Service) GetRecord(startTime int64, endTime int64) (int,error) {
	//fmt.Println("gooooooooooooooooooooooooooooooo")
	mapdata := make([]map[string]interface{},0)
	//koala.Init(configs.Gconf.KoalaHost)
	//koala.KoalaLogin(configs.Gconf.KoalaUsername,configs.Gconf.KoalaPassword)
	first_req_data, total, err := koala.GetEventsUser(startTime, endTime, 1)
	//log4go.Info("writeSug file.Write error(%v)",datas)
	log4go.Info("writeSug file.Write error-------------------------------",total)
	log4go.Info("first_req_data-------------------------------",first_req_data)
	if err != nil {
		fmt.Println(err)
		return -100, err
	}
	if total >1 {
		for i:= int64(1); i<=total;i++{
			temp := make([]map[string]interface{},0)
			data, _, err := koala.GetEventsUser(startTime, endTime, i)
			if err != nil{
				return 0,err
			}
			temp = append(temp,data...)
			s.GoAddIdentificationRecord(temp)
		}
		return 1,nil
	}
	if total == 1 {
		mapdata = append(mapdata,first_req_data...)
	}
	go s.GoAddIdentificationRecord(mapdata)
	return 0, nil
}

func (s *Service)GoAddIdentificationRecord(datas []map[string]interface{})(){
	var(
		stranger model.Stranger
		identificationRecord model.IdentificationRecord
		//tmpscreen_id int
	)
	//timer,_ := s.GetTime()
	//if timer.Screen_id != 0{
	//	tmpscreen_id = timer.Screen_id
	//}
	for _,v := range datas{
		log4go.Info("datas",datas)
		if _,ok := v["subject"];ok && v["subject_id"] != nil{
			subject := v["subject"].(map[string]interface{})
			screen := v["screen"].(map[string]interface{})
			subject_type, _ := subject["subject_type"].(json.Number).Int()
			log4go.Info("该条历史的subject_type为：" + strconv.Itoa(subject_type))
			photoURL := subject["avatar"].(string)
			var subject_id int
			if subject["id"] != nil{
				subject_id,_ = subject["id"].(json.Number).Int()
			}else {
				subject_id = 0
				fmt.Println(subject)
			}
			var screen_id int
			if strings.Contains(screen["camera_position"].(string),"已删除"){
				screen_id = -99
				identificationRecord.Snap_position = screen["camera_position"].(string)
			}else{
				screen_id,_ = screen["id"].(json.Number).Int()
				identificationRecord.Snap_position = screen["camera_position"].(string)
			}
			name := subject["name"].(string)
			identificationRecord.Photo = photoURL
			identificationRecord.Subject_type = subject_type
			identificationRecord.Event_type = 0
			identificationRecord.Subject_id = subject_id
			identificationRecord.Screen_id = screen_id
			identificationRecord.Snap_photo = v["photo"].(string)
			identificationRecord.Name = name
			if value,ok :=subject["start_time"].(json.Number);ok{
				start_time, _ := value.Int64()
				start_time_str := time.Unix(start_time, 0).Format("2006-01-02 15:04:05")
				identificationRecord.Start_time = start_time_str
			}
			if value,ok :=subject["end_time"].(json.Number);ok{
				end_time, _ := value.Int64()
				end_time_str := time.Unix(end_time, 0).Format("2006-01-02 15:04:05")
				identificationRecord.End_time = end_time_str
			}
			if value,ok :=subject["come_from"].(string);ok{
				identificationRecord.Come_from = value
			}
			if value,ok :=subject["remark"].(string);ok{
				identificationRecord.Remark = value
			}
			snapTime, _ := v["timestamp"].(json.Number).Int64()
			date_time := time.Unix(snapTime, 0)
			str_time := date_time.Format("2006-01-02 15:04:05")
			var timestamp time.Time
			identificationRecord.Snap_time = str_time
			timestamp, _ = time.Parse("2006-01-02", str_time[0:10])
			identificationRecord.Date_time = timestamp
			err := s.dao.AddIdentificationRecord(identificationRecord)
			if err != nil {
				log4go.Error("该数据插入DB出错")
				log4go.Error("该记录为：",identificationRecord)
				continue
			}
		}else{
			screen := v["screen"].(map[string]interface{})
			var screen_id int
			var snap_position string
			var snap_photo string
			var str_time string
			var snapTime int64
			var timestamp time.Time
			if strings.Contains(screen["camera_position"].(string),"已删除"){
				screen_id = -99
				snap_position=screen["camera_position"].(string)
			}else{
				screen_id,_ = screen["id"].(json.Number).Int()
				snap_position=screen["camera_position"].(string)
			}
			//if screen_id != tmpscreen_id{
			//	continue
			//}
			snap_photo=v["photo"].(string)
			snapTime, _ = v["timestamp"].(json.Number).Int64()
			date_time := time.Unix(snapTime, 0)
			str_time = date_time.Format("2006-01-02 15:04:05")
			timestamp, _ = time.Parse("2006-01-02", str_time[0:10])

			event_type, _ := v["event_type"].(json.Number).Int64()
			if event_type == 0 {
				stranger.Event_type = 0
				stranger.Snap_position = snap_position
				stranger.Screen_id = screen_id
				stranger.Snap_photo=snap_photo
				stranger.Snap_time = str_time
				stranger.Date_time = timestamp
				err := s.dao.AddStrangerRecord(stranger)
				if err != nil {
					log4go.Error("该数据插入DB出错")
					log4go.Error("该记录为：",stranger)
					continue
				}
			} else if event_type == 1 {
				identificationRecord_ := model.IdentificationRecord{}
				identificationRecord_.Event_type = 1
				identificationRecord_.Snap_position = snap_position
				identificationRecord_.Screen_id = screen_id
				identificationRecord_.Snap_photo=snap_photo
				identificationRecord_.Snap_time = str_time
				identificationRecord_.Date_time = timestamp
				identificationRecord_.Subject_type = 4
				if value,ok :=v["start_time"].(json.Number);ok{
					start_time, _ := value.Int64()
					start_time_str := time.Unix(start_time, 0).Format("2006-01-02 15:04:05")
					identificationRecord_.Start_time = start_time_str
				}
				if value,ok :=v["end_time"].(json.Number);ok{
					end_time, _ := value.Int64()
					end_time_str := time.Unix(end_time, 0).Format("2006-01-02 15:04:05")
					identificationRecord_.End_time = end_time_str
				}
				if value,ok :=v["come_from"].(string);ok{
					identificationRecord_.Come_from = value
				}
				if value,ok :=v["remark"].(string);ok{
					identificationRecord_.Remark = value
				}
				err := s.dao.AddIdentificationRecord(identificationRecord_)
				if err != nil {
					log4go.Error("该数据插入DB出错")
					log4go.Error("该记录为：",identificationRecord_)
					continue
				}
			}
		}
	}
	return
}

func (s *Service)PrepareForstatistic(datetime string,Position_map map[int]interface{}) (error) {

	if datetime != ""{

	}else{
		nowtime := time.Now()
		nowdate := nowtime.String()[0:10]
		datetime = nowdate
		}
	res , err := s.GetTime()
	Desayuno_begin:= " "+res.Desayuno_begin
	//Desayuno_begin:= " 07:00:00"
	Desayuno_end:= " "+res.Desayuno_end
	//Desayuno_end:= " 10:00:00"
	Almuerzo_begin:= " "+res.Almuerzo_begin
	//Almuerzo_begin:= " 11:00:00"
	Almuerzo_end:= " "+res.Almuerzo_end
	//Almuerzo_end:= " 13:00:00"
	Jantar_begin:= " "+res.Jantar_begin
	//Jantar_begin:= " 16:30:00"
	Jantar_end:= " "+res.Jantar_end
	//Jantar_end:= " 18:30:00"


	DesayunoCount,err := s.DesayunoCount(datetime+Desayuno_begin,datetime+Desayuno_end,Position_map)
	if err != nil{
		log4go.Error(err)
		return errors.New("查询日间统计数据失败")
	}
	AlmuerzoCount,err := s.AlmuerzoCount(datetime+Almuerzo_begin,datetime+Almuerzo_end,Position_map)
	if err != nil{
		log4go.Error(err)
		return errors.New("查询午间食堂统计数据失败")
	}
	JantarCount,err := s.JantarCount(datetime+Jantar_begin,datetime+Jantar_end,Position_map)
	if err != nil{
		log4go.Error(err)
		return errors.New("查询夜间食堂统计数据失败")
	}


	statistic := model.Statistic{
		Jantar:JantarCount,
		Almuerzo:AlmuerzoCount,
		Desayuno:DesayunoCount,
		Date:datetime,
	}
	err = s.dao.CreateStatistic(statistic)
	if err!= nil{
		return errors.New("插入统计数据失败")
	}
	return nil
}


func (s *Service)Statistic_main_item_map() (map[int][]model.Statistic_Item,[]model.Statistic_Main,error) {
	main_item := make(map[int][]model.Statistic_Item)
	stastics ,err := s.Statistic_Project_Parent_Get()
	if err != nil{
		log4go.Error(err.Error())
		return map[int][]model.Statistic_Item{},[]model.Statistic_Main{},err
	}
	for _,v := range stastics {
		//获取该主项的子项
		items, err := s.Statistic_Project_Items_Get(v.Id)
		if err != nil {
			log4go.Error(err.Error())
			continue
		}
		items_ids := make([]model.Statistic_Item,0)
		if len(items) == 0{
			main_item[v.Id] = []model.Statistic_Item{}
			continue
		}
		for k,_ := range items{
			items_ids = append(items_ids,items[k])
		}
		main_item[v.Id] = items_ids
	}
	log4go.Info("main_item",main_item)
	log4go.Info("stastics",stastics)
	return main_item,stastics,nil
}
var flag2 = 0

func (s *Service)PrepareForstatistic2(datetime string,Position_map map[int]interface{}) (map[int]interface{},error) {
//func (s *Service)PrepareForstatistic2() (error) {
	if flag2 == 1{
		return	nil,nil
	}
	flag2 = 1
	defer func() {
		flag2 = 0
	}()
	if datetime != ""{

	}else{
		nowtime := time.Now()
		nowdate := nowtime.String()[0:10]
		datetime = nowdate
	}
	main_item := make(map[int]interface{})
	stastics ,err := s.Statistic_Project_Parent_Get()
	if err != nil{
		log4go.Error(err.Error())
		return nil,err
	}
	for _,v := range stastics{
		//获取该主项的子项
		items , err := s.Statistic_Project_Items_Get(v.Id)
		if err != nil{
			log4go.Error(err.Error())
			return nil,err
		}
		if len(items) == 0{
			continue
		}
		items_maplist := make([]map[string]interface{},0)
		for _,i:= range items{
			log4go.Info("---------------------------------------------get GetChildren")
			items_map :=GetChildren2(i)
			if items_map == nil{
				continue
			}
			items_maplist = append(items_maplist,items_map)
		}
		main_item[v.Id] = items_maplist
	}

	s.Stastic_Main_Item2(datetime,main_item,stastics)

	log4go.Info("-------------------------------------------------------------main_item")
	log4go.Info(main_item)
	//fmt.Println(datetime)
	log4go.Info(Position_map)

	//fmt.Println(time.Now().Unix())
	//fmt.Println("function end --------------------------------------------------")
	log4go.Info("function end ---------------------------------------------------",time.Now().Unix())
	fmt.Println("function end ---------------------------------------------------",time.Now().Unix())
	return main_item,nil

	//DesayunoCount,err := s.DesayunoCount(datetime+Desayuno_begin,datetime+Desayuno_end,Position_map)
	//if err != nil{
	//	log4go.Error(err)
	//	return errors.New("查询日间统计数据失败")
	//}
	//AlmuerzoCount,err := s.AlmuerzoCount(datetime+Almuerzo_begin,datetime+Almuerzo_end,Position_map)
	//if err != nil{
	//	log4go.Error(err)
	//	return errors.New("查询午间食堂统计数据失败")
	//}
	//JantarCount,err := s.JantarCount(datetime+Jantar_begin,datetime+Jantar_end,Position_map)
	//if err != nil{
	//	log4go.Error(err)
	//	return errors.New("查询夜间食堂统计数据失败")
	//}
	//
	//
	//statistic := model.Statistic{
	//	Jantar:JantarCount,
	//	Almuerzo:AlmuerzoCount,
	//	Desayuno:DesayunoCount,
	//	Date:datetime,
	//}
	//err = s.dao.CreateStatistic(statistic)
	//if err!= nil{
	//	return errors.New("插入统计数据失败")
	//}
	//return nil
}

func (s *Service)Stastic_Main_Item(datetime string,main_item_map map[int]interface{},main []model.Statistic_Main)(){

	for _,v := range main{
		if main_item_map[v.Id] == nil{
			continue
		}
		main_id := v.Id
			//"screens":screenids,
			//"subjects":subjectids,
			//"subject_type":subject_type,
			//"start_time":item.Start_time,
			//"end_time":item.End_time,
			//"name":item.Name,
		items := main_item_map[v.Id].([]map[string]interface {})
		log4go.Info(items)
		log4go.Info("-------------------------------------------------------------items")
		for _,i :=range items{
			//知道了main_id 和 item_id
			item_id := i["id"].(int)
			screens := i["screens"].([]int)
			screens_group := i["screens_group"].(int)
			subjects_group := i["subjects_group"].(int)
			subjectids := i["subjects"]
			subject_type := i["subject_type"].(int)
			start_time := i["start_time"].(string)
			if start_time == ""{
				start_time = "00:00:00"
			}
			end_time := i["end_time"].(string)
			if end_time == ""{
				end_time = "23:59:59"
			}
			name := i["name"].(string)
			s.counter(name,screens,subjectids.([]int),subjects_group,screens_group,main_id,item_id,subject_type,datetime+" "+start_time,datetime+" "+end_time,datetime)
//(screens []int,subject []int,subjects_group ,screens_group,main_id,item_id,subject_type int,start_time,end_time,datetime string)

				//times := 0
			//var lock sync.Mutex
			//if len(subjectids.([]int)) == 0 && len(screens) == 0{
			//	//全相机 全人员
			//	subject_objs,err := s.GetoneItemStatistic4(datetime+" "+start_time,datetime+" "+end_time)
			//	if err != nil{
			//		log4go.Error("查询统计子项出错,该项名称为",name)
			//		return
			//	}
			//	times = len(subject_objs)
			//
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else if len(subjectids.([]int)) == 0 && len(screens) > 0{
			//	//部分相机 全人员
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else if len(subjectids.([]int)) > 0 && len(screens) == 0{
			//	//全相机 部分人员
			//	//GetOneItemStatistic3
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else
			// if  len(screens) > 0{
			//	//for _,o :=range subjectids.([]int){
			//
			//		//subject_objs,err := s.GetoneItemStatistic(screens,o,subject_type,datetime+" "+start_time,datetime+" "+end_time)
			//		subject_objs,err := s.GetoneItemStatistic_copy(screens,subjectids.([]int),subject_type,datetime+" "+start_time,datetime+" "+end_time)
			//		if err != nil{
			//			log4go.Error("查询统计子项出错,该项名称为",name)
			//			continue
			//		}
			//		if len(subject_objs) != 0{
			//			times = len(subject_objs)
			//		}
			//	//}
			//}

			//更新子项统计数值
			//update_data := model.Statistics{
			//	Number:times,
			//	End_time:end_time,
			//	Start_time:start_time,
			//	Date:datetime,
			//	Subjects_group:subjects_group,
			//	Screens_group:screens_group,
			//	Grandparent_id:main_id,
			//	Parent_id:item_id,
			//}
			//s.UpdateoneItemStatistic(update_data)
			}
	}
}
func (s *Service)Stastic_Main_Item2(datetime string,main_item_map map[int]interface{},main []model.Statistic_Main)(){

	for _,v := range main{
		if main_item_map[v.Id] == nil{
			continue
		}
		main_id := v.Id
		//"screens":screenids,
		//"subjects":subjectids,
		//"subject_type":subject_type,
		//"start_time":item.Start_time,
		//"end_time":item.End_time,
		//"name":item.Name,
		items := main_item_map[v.Id].([]map[string]interface {})
		log4go.Info(items)
		log4go.Info("-------------------------------------------------------------items")
		for _,i :=range items{
			//知道了main_id 和 item_id
			item_id := i["id"].(int)
			screens := i["screens"].([]int)
			screens_group := i["screens_group"].(string)
			subjects_group := i["subjects_group"].(int)
			subjectids := i["subjects"]
			subject_type := i["subject_type"].(int)
			start_time := i["start_time"].(string)
			if start_time == ""{
				start_time = "00:00:00"
			}
			end_time := i["end_time"].(string)
			if end_time == ""{
				end_time = "23:59:59"
			}
			name := i["name"].(string)
			s.counter2(name,screens_group,screens,subjectids.([]int),subjects_group,main_id,item_id,subject_type,datetime+" "+start_time,datetime+" "+end_time,datetime)
			//(screens []int,subject []int,subjects_group ,screens_group,main_id,item_id,subject_type int,start_time,end_time,datetime string)

			//times := 0
			//var lock sync.Mutex
			//if len(subjectids.([]int)) == 0 && len(screens) == 0{
			//	//全相机 全人员
			//	subject_objs,err := s.GetoneItemStatistic4(datetime+" "+start_time,datetime+" "+end_time)
			//	if err != nil{
			//		log4go.Error("查询统计子项出错,该项名称为",name)
			//		return
			//	}
			//	times = len(subject_objs)
			//
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else if len(subjectids.([]int)) == 0 && len(screens) > 0{
			//	//部分相机 全人员
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else if len(subjectids.([]int)) > 0 && len(screens) == 0{
			//	//全相机 部分人员
			//	//GetOneItemStatistic3
			//	//res ,count,total, err := s.GetOneStatisticOneRecord(Position_map_copy,snap_begin_time[0:10],item_map,gdata,main_id,item_id,page,size)
			//}else
			// if  len(screens) > 0{
			//	//for _,o :=range subjectids.([]int){
			//
			//		//subject_objs,err := s.GetoneItemStatistic(screens,o,subject_type,datetime+" "+start_time,datetime+" "+end_time)
			//		subject_objs,err := s.GetoneItemStatistic_copy(screens,subjectids.([]int),subject_type,datetime+" "+start_time,datetime+" "+end_time)
			//		if err != nil{
			//			log4go.Error("查询统计子项出错,该项名称为",name)
			//			continue
			//		}
			//		if len(subject_objs) != 0{
			//			times = len(subject_objs)
			//		}
			//	//}
			//}

			//更新子项统计数值
			//update_data := model.Statistics{
			//	Number:times,
			//	End_time:end_time,
			//	Start_time:start_time,
			//	Date:datetime,
			//	Subjects_group:subjects_group,
			//	Screens_group:screens_group,
			//	Grandparent_id:main_id,
			//	Parent_id:item_id,
			//}
			//s.UpdateoneItemStatistic(update_data)
		}
	}
}
//var C = make(chan int , 5000)
//var counter = 0
func (s *Service)counter(name string,screens []int,subject []int,subjects_group ,screens_group,main_id,item_id,subject_type int,start_time,end_time,datetime string){
	//time.Sleep(50*time.Millisecond)
	var times = 0
	subject_objs,err  := s.GetoneItemStatistic_copy(screens,subject,subject_type,start_time,end_time)
	if err != nil{
		log4go.Error("查询统计子项出错,该项名称为",name)
		return
	}
	if len(subject_objs) != 0{
		times = len(subject_objs)
	}
	update_data := model.Statistics{
		Number:times,
		End_time:end_time,
		Start_time:start_time,
		Date:datetime,
		Subjects_group:subjects_group,
		Screens_group:screens_group,
		Grandparent_id:main_id,
		Parent_id:item_id,
	}
	log4go.Info("update_data is ----------------",update_data)
	go func(){
		s.UpdateoneItemStatistic(update_data)
	}()
}
func (s *Service)counter2(name,screens_group string,screens []int,subject []int,subjects_group ,main_id,item_id,subject_type int,start_time,end_time,datetime string){
	//time.Sleep(50*time.Millisecond)
	var times = 0
	subject_objs,err  := s.GetoneItemStatistic_copy(screens,subject,subject_type,start_time,end_time)
	if err != nil{
		log4go.Error("查询统计子项出错,该项名称为",name)
		return
	}
	if len(subject_objs) != 0{
		times = len(subject_objs)
	}
	update_data := model.Statistics{
		Number:times,
		End_time:end_time,
		Start_time:start_time,
		Date:datetime,
		Subjects_group:subjects_group,
		Screens_groups:screens_group,
		Grandparent_id:main_id,
		Parent_id:item_id,
	}
	log4go.Info("update_data is ----------------",update_data)
	go func(){
		s.UpdateoneItemStatistic(update_data)
	}()
}

func (s *Service)GetOneStatisticRecords(position_map map[int]interface{},datetime string,item map[int][]model.Statistic_Item,main []model.Statistic_Main,main_id,item_id,subject_id,page,size int)([]model.IdentificationRecord,int,int,error){

	var(
		subject_objs []model.IdentificationRecord
		err error
		count,total int
	)
	if item[main_id] == nil{
		return nil,0,0, nil
	}
	items := item[main_id]
	//main_item := make(map[int]interface{})

	items_maplist := make([]map[string]interface{},0)
	for _,i:= range items{

		items_map :=GetChildren2(i)
		if items_map == nil{
			continue
		}
		items_maplist = append(items_maplist,items_map)
	}

	//if main_item[main_id] == nil{
	//	return nil, nil
	//}
		//main_id := v.Id
		//items := main_item[main_id].([]map[string]interface {})
		for _,i :=range items_maplist{
			if item_id != i["id"].(int){
				continue
			}
			//知道了main_id 和 item_id
			//item_id := i["id"].(int)
			screens := i["screens"].([]int)
			//screens_group := i["screens_group"].(int)
			//subjects_group := i["subjects_group"].(int)
			//subjectids := i["subjects"]
			subject_type := i["subject_type"].(int)
			start_time := i["start_time"].(string)
			if start_time == ""{
				start_time = "00:00:00"
			}
			end_time := i["end_time"].(string)
			if end_time == ""{
				end_time = "23:59:59"
			}
			name := i["name"].(string)
			subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,subject_id,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
			if err != nil{
				log4go.Error("查询统计子项出错,该项名称为",name)
				return nil,0,0, err
			}
			for k,v := range subject_objs{
				if v.Photo != ""{
					tmp_photo := "http://" + configs.Gconf.KoalaHost +v.Photo
					subject_objs[k].Photo = tmp_photo
				}
				if v.Snap_photo != ""{
					tmp_Snap_photo:= "http://" + configs.Gconf.KoalaHost +v.Snap_photo
					subject_objs[k].Snap_photo = tmp_Snap_photo
				}
				if v.Snap_position == ""{
					tmp_Snap_position := position_map[v.Screen_id].(string)
					subject_objs[k].Snap_position = tmp_Snap_position
				}
				//results[k].Recognition_times = headcountmap[v.Subject_id]
			}
		}
	return subject_objs,count,total, nil
}
func (s *Service)GetOneStatisticOneRecord2(position_map map[int]interface{},datetime string,item map[int][]model.Statistic_Item,main []model.Statistic_Main,main_id,item_id,page,size int)([]model.IdentificationRecord,int,int,error){

	var(
		subject_objs []model.IdentificationRecord
		resp_identify []model.IdentificationRecord
		err error
		count,total int
	)
	if item[main_id] == nil{
		fmt.Println(item)
		fmt.Println("item not find")
		return nil,0,0, nil
	}
	items := item[main_id]
	//main_item := make(map[int]interface{})
	fmt.Println("items is -----------------------------------------",items)
	items_maplist := make([]map[string]interface{},0)
	for _,i:= range items{

		items_map :=GetChildren2(i)
		if items_map == nil{
			continue
		}
		items_maplist = append(items_maplist,items_map)
	}

	//if main_item[main_id] == nil{
	//	return nil, nil
	//}
	//main_id := v.Id
	//items := main_item[main_id].([]map[string]interface {})
	//total_subjects := make([]int,0)
	//for _,i :=range items_maplist{
	//	subjectids := i["subjects"].([]int)
	//	total_subjects = append(total_subjects, subjectids...)
	//}
	//fmt.Println("total_subjects",len(total_subjects))
	//fmt.Println("items_maplist",items_maplist)


	for _,i :=range items_maplist{
		if item_id != i["id"].(int){
			continue
		}
		//知道了main_id 和 item_id
		//item_id := i["id"].(int)
		screens := i["screens"].([]int)
		//screens_group := i["screens_group"].(int)
		//subjects_group := i["subjects_group"].(int)
		subjectids := i["subjects"].([]int)
		subject_type := i["subject_type"].(int)
		start_time := i["start_time"].(string)
		if start_time == ""{
			start_time = "00:00:00"
		}
		end_time := i["end_time"].(string)
		if end_time == ""{
			end_time = "23:59:59"
		}
		name := i["name"].(string)

		//resp_identify := make([]model.IdentificationRecord,0)
		//length = len(subjectids)
		//if length==0{
		//	fmt.Println("length is 0")
		//	return nil, 0, 0, nil
		//}
		////数组切片
		//max_index := page*size//
		//var(
		//	begin_index,end_index int
		//)
		//if max_index <= 10{
		//	begin_index = 0//初始位置
		//}else{
		//	begin_index = max_index - 10//初始位置
		//}
		//
		//end_index = 0
		//
		//if length > max_index{
		//	end_index = max_index -1
		//}
		//if length < max_index {
		//	if max_index - length <10{
		//		end_index = max_index-(max_index - length)-1
		//	}else{
		//		return nil, 0, 0, nil
		//	}
		//}
		//fmt.Println("begin_index is ",begin_index)
		//fmt.Println("end_index is ",end_index)
		//fmt.Println("subjectids is  ",subjectids)
		//fmt.Println("subjectids range is  ",subjectids[begin_index:end_index])
		fmt.Println("subjectids",len(subjectids))
		if len(subjectids) != 0 {
			count = len(subjectids)
			m_max := len(subjectids)
			fmt.Println("m_max",m_max)
			total := zyutil.GetTotal(m_max,size)
			fmt.Println("total",total)
			supplement := total*size - m_max
			fmt.Println("supplement",supplement)
			for i:= 1 ; i<supplement;i++{
				subjectids = append(subjectids, 0)
			}

			var max = 0
			var min  = 0
			if page >1{

				max = size * page
				min = (page - 1)*size
				fmt.Println("page > 1",max)
				fmt.Println("page > 1",min)
			}else if page == 1{
				max = size * page
				min = 0
				fmt.Println("page < 1",max)
				fmt.Println("page < 1",min)
			}
			//fmt.Println(subjectids)
			for _,vi := range subjectids[min:]{
			//for _,vi := range subjectids{
				//page =1
				//size = 10000
				subject_objs,_,_,err = s.GetOneItemStatisticRecords(screens,vi,subject_type,1,10000,datetime+" "+start_time,datetime+" "+end_time)
				if err != nil{
					log4go.Error("查询统计子项出错,该项名称为",name)
					return nil,0,0, err
				}
				if len(subject_objs) != 0 {
					if subject_objs[0].Photo != ""{
						tmp_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[0].Photo
						subject_objs[0].Photo = tmp_photo
					}
					if subject_objs[0].Snap_photo != ""{
						tmp_Snap_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[0].Snap_photo
						subject_objs[0].Snap_photo = tmp_Snap_photo
					}
					if subject_objs[0].Snap_position == "" {
						tmp_Snap_position := position_map[subject_objs[0].Screen_id].(string)
						subject_objs[0].Snap_position = tmp_Snap_position
					}
					subject_objs[0].Recognition_times = len(subject_objs)
				}else{
					continue
				}
				resp_identify= append(resp_identify,subject_objs[0])
				if len(resp_identify) == 10{
					break
				}
				//if len(resp_identify) == max_index {
				//	break
				//}
			}
		}else{
			//for _,vi := range subjectids{
			page =1
			size = 10000
			subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,-99,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
			if err != nil{
				log4go.Error("查询统计子项出错,该项名称为",name)
				return nil,0,0, err
			}
			if len(subject_objs) != 0 {
				for k,_ :=range subject_objs{
					tmp_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[k].Photo
					subject_objs[k].Photo = tmp_photo
					tmp_Snap_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[k].Snap_photo
					subject_objs[k].Snap_photo = tmp_Snap_photo
					if subject_objs[k].Snap_position == "" {
						if position_map[subject_objs[k].Screen_id] != nil{
							tmp_Snap_position := position_map[subject_objs[k].Screen_id].(string)
							subject_objs[k].Snap_position = tmp_Snap_position
						}
					}
					_,count,total,err = s.GetOneItemStatisticRecords(screens,subject_objs[k].Subject_id,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
					if err != nil{
						log4go.Error("查询统计子项出错,该项名称为",name)
						continue
					}
					subject_objs[k].Recognition_times = count
				}
			}
			resp_identify= append(resp_identify,subject_objs...)
		}

		//subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,subject_id,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
		//if err != nil{
		//	log4go.Error("查询统计子项出错,该项名称为",name)
		//	return nil,0,0, err
		//}
	}
	if count != 0 {

	}else{
		count = len(resp_identify)
	}
	total = zyutil.GetTotal(count, size)
	return resp_identify,count,total, nil
}


func (s *Service)GetOneStatisticOneRecord(position_map map[int]interface{},datetime string,item map[int][]model.Statistic_Item,main []model.Statistic_Main,main_id,item_id,page,size int)([]model.IdentificationRecord,int,int,error){

	var(
		subject_objs []model.IdentificationRecord
		resp_identify []model.IdentificationRecord
		err error
		count,total int
	)
	if item[main_id] == nil{
		fmt.Println("item not find")
		return nil,0,0, nil
	}
	items := item[main_id]
	//main_item := make(map[int]interface{})

	items_maplist := make([]map[string]interface{},0)
	for _,i:= range items{

		items_map :=GetChildren2(i)
		if items_map == nil{
			continue
		}
		items_maplist = append(items_maplist,items_map)
	}

	//if main_item[main_id] == nil{
	//	return nil, nil
	//}
	//main_id := v.Id
	//items := main_item[main_id].([]map[string]interface {})
	for _,i :=range items_maplist{
		if item_id != i["id"].(int){
			continue
		}
		//知道了main_id 和 item_id
		//item_id := i["id"].(int)
		screens := i["screens"].([]int)
		//screens_group := i["screens_group"].(int)
		//subjects_group := i["subjects_group"].(int)
		subjectids := i["subjects"].([]int)
		subject_type := i["subject_type"].(int)
		start_time := i["start_time"].(string)
		if start_time == ""{
			start_time = "00:00:00"
		}
		end_time := i["end_time"].(string)
		if end_time == ""{
			end_time = "23:59:59"
		}
		name := i["name"].(string)

		//resp_identify := make([]model.IdentificationRecord,0)
		//length = len(subjectids)
		//if length==0{
		//	fmt.Println("length is 0")
		//	return nil, 0, 0, nil
		//}
		////数组切片
		//max_index := page*size//
		//var(
		//	begin_index,end_index int
		//)
		//if max_index <= 10{
		//	begin_index = 0//初始位置
		//}else{
		//	begin_index = max_index - 10//初始位置
		//}
		//
		//end_index = 0
		//
		//if length > max_index{
		//	end_index = max_index -1
		//}
		//if length < max_index {
		//	if max_index - length <10{
		//		end_index = max_index-(max_index - length)-1
		//	}else{
		//		return nil, 0, 0, nil
		//	}
		//}
		//fmt.Println("begin_index is ",begin_index)
		//fmt.Println("end_index is ",end_index)
		//fmt.Println("subjectids is  ",subjectids)
		//fmt.Println("subjectids range is  ",subjectids[begin_index:end_index])

		if len(subjectids) != 0 {
			for _,vi := range subjectids{
				page =1
				size = 10000
				subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,vi,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
				if err != nil{
					if strings.Contains(err.Error(),"Server shutdown in progress"){

					}
					log4go.Error("查询统计子项出错,该项名称为",name)
					return nil,0,0, err
				}
				if len(subject_objs) != 0 {
					if subject_objs[0].Photo != ""{
						tmp_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[0].Photo
						subject_objs[0].Photo = tmp_photo
					}
					if subject_objs[0].Snap_photo != ""{
						tmp_Snap_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[0].Snap_photo
						subject_objs[0].Snap_photo = tmp_Snap_photo
					}
					if subject_objs[0].Snap_position == "" {
						tmp_Snap_position := position_map[subject_objs[0].Screen_id].(string)
						subject_objs[0].Snap_position = tmp_Snap_position
					}
					subject_objs[0].Recognition_times = len(subject_objs)
				}else{
					continue
				}
				resp_identify= append(resp_identify,subject_objs[0])
				//if len(resp_identify) == max_index {
				//	break
				//}
			}
		}else{
			//for _,vi := range subjectids{
				page =1
				size = 10000
				subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,-99,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
				if err != nil{
					log4go.Error("查询统计子项出错,该项名称为",name)
					return nil,0,0, err
				}
				if len(subject_objs) != 0 {
					for k,_ :=range subject_objs{
						tmp_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[k].Photo
						subject_objs[k].Photo = tmp_photo
						tmp_Snap_photo := "http://" + configs.Gconf.KoalaHost + subject_objs[k].Snap_photo
						subject_objs[k].Snap_photo = tmp_Snap_photo
						if subject_objs[k].Snap_position == "" {
							if position_map[subject_objs[k].Screen_id] != nil{
								tmp_Snap_position := position_map[subject_objs[k].Screen_id].(string)
								subject_objs[k].Snap_position = tmp_Snap_position
							}
						}
						_,count,total,err = s.GetOneItemStatisticRecords(screens,subject_objs[k].Subject_id,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
						if err != nil{
							log4go.Error("查询统计子项出错,该项名称为",name)
							continue
						}
						subject_objs[k].Recognition_times = count
					}
				}
				resp_identify= append(resp_identify,subject_objs...)
			}

		//subject_objs,count,total,err = s.GetOneItemStatisticRecords(screens,subject_id,subject_type,page,size,datetime+" "+start_time,datetime+" "+end_time)
		//if err != nil{
		//	log4go.Error("查询统计子项出错,该项名称为",name)
		//	return nil,0,0, err
		//}
	}



	count = len(resp_identify)
	total = zyutil.GetTotal(count, size)
	return resp_identify,count,total, nil
}


//func(s *Service)GetoneItemStatistic(screens,subjects []int,subject_type int,start_time,end_time string)(){
func(s *Service)GetoneItemStatistic(screens []int,subject,subject_type int,start_time,end_time string)([]model.IdentificationRecord,error){

	subject_objs,err := s.dao.GetOneItemStatistic(screens,subject,subject_type,start_time,end_time)
	return subject_objs,err

}
func(s *Service)GetoneItemStatistic_copy(screens []int,subject []int ,subject_type int,start_time,end_time string)([]model.IdentificationRecord,error){

	subject_objs,err := s.dao.GetOneItemStatistic_copy(screens,subject,subject_type,start_time,end_time)
	return subject_objs,err

}
func(s *Service)GetOneItemStatistic3(subject,subject_type int,start_time,end_time string)([]model.IdentificationRecord,error){

	subject_objs,err := s.dao.GetOneItemStatistic3(subject,subject_type,start_time,end_time)
	return subject_objs,err

}

func(s *Service)GetoneItemStatistic4(start_time,end_time string)([]model.IdentificationRecord,error){

	subject_objs,err := s.dao.GetOneItemStatistic4(start_time,end_time)
	return subject_objs,err

}

func(s *Service)GetOneItemStatisticRecords(screens []int,subject,subject_type,page,size int,start_time,end_time string)([]model.IdentificationRecord,int,int,error){

	subject_objs,count,total,err := s.dao.GetOneItemStatisticRecords(screens,subject,subject_type,page,size,start_time,end_time)

	return subject_objs,count,total,err

}


func(s *Service)UpdateoneItemStatistic(statistic model.Statistics)(){

	//subject_objs,err := s.dao.UpdateOneItemStatistic(statistic)
	s.dao.UpdateOneItemStatistic(statistic)
	return
	//return subject_objs,err

}

var temp_map_screengroupid = make(map[int][]int)
var  temp_map_subjectgroupid = make(map[int][]int)
var temp_map_subjectgrouptype = make(map[int]int)

func GetChildren(item model.Statistic_Item)(map[string]interface{}){
	var(
		err error
		subject_type int
	)
	screenids := make([]int,0)
	subjectids := make([]int,0)
	item_id := item.Id
	screengroupid := item.Screens_group
	subjectgroupid := item.Subjects_group
	fmt.Println("---------------------------------------------get screengroupid",screengroupid)
	fmt.Println("---------------------------------------------get subjectgroupid",subjectgroupid)
	if screengroupid == -1{
		fmt.Println("---------------------------------------------get screengroupid is -1")
		screen_list ,err := koala.GetScreenList()
		if err!=nil{
			log4go.Error("get children failed, reason is :",err)
			return nil
		}
		fmt.Println("---------------------------------------------get screen_list")
		data_list,err := screen_list.Get("data").Array()
		if err!=nil{
			log4go.Error("screen list is not array, reason is :",err)
			return nil
		}
		if len(data_list) == 0 {
			log4go.Error("no screen")
			return nil
		}
		for _,v := range  data_list{
			id ,err := v.(interface{}).(map[string]interface{})["id"].(json.Number).Int()
			if err!=nil{
				log4go.Error("get screen id failed, reason is :",err)
				return nil
			}
			screenids = append(screenids,id)
		}
		temp_map_screengroupid[screengroupid] = screenids
		fmt.Println("---------------------------------------------get data_list")
	}else{
	//if screengroupid == -1{
	//	temp_map_screengroupid[screengroupid] = screenids
	//}else{
		screenids,err = koala.GetAccess(screengroupid)
		if err != nil{
			log4go.Error(err.Error())
			return nil
		}
		temp_map_screengroupid[screengroupid] = screenids
	}


	//}
	fmt.Println("---------------------------------------------get subjectgroupid")
	//if subjectgroupid == -1{
	//	employee,err := koala.GetSubjects("employee")
	//
	//	if err != nil{
	//		fmt.Println("---------------------------------------------get employee err")
	//		log4go.Error(err.Error())
	//		return nil
	//	}
	//	visitor,err :=koala.GetSubjects("visitor")
	//
	//	if err != nil{
	//		fmt.Println("---------------------------------------------get visitor err")
	//		log4go.Error(err.Error())
	//		return nil
	//	}
	//	employeeids,employeetype,err := koala.GetPersonData(employee)
	//	if err != nil{
	//		fmt.Println("---------------------------------------------get employee err")
	//		log4go.Error(err)
	//	}
	//	subjectids = append(subjectids, employeeids...)
	//	visitorids,_,err := koala.GetPersonData(visitor)
	//	if err != nil{
	//		fmt.Println("---------------------------------------------get visitor err")
	//		log4go.Error(err)
	//	}
	//	subjectids = append(subjectids, visitorids...)
	//	temp_map_subjectgroupid[subjectgroupid] = subjectids
	//	temp_map_subjectgrouptype[subjectgroupid] = employeetype
	//}else{
	if subjectgroupid == -1{
		temp_map_subjectgroupid[subjectgroupid] = subjectids
		temp_map_subjectgrouptype[subjectgroupid] = 0
	}else{
		subjectids,subject_type,err = koala.GetPerson(subjectgroupid)
		if err != nil{
			log4go.Error(err.Error())
			return nil
		}
		temp_map_subjectgroupid[subjectgroupid] = screenids
		temp_map_subjectgrouptype[subjectgroupid] = subject_type
	}

	//}
	//if temp_map_screengroupid[screengroupid] != nil{
	//if temp_map_screengroupid[screengroupid] == nil{
		//fmt.Println("get_data --------------------------------------")
		//screenids = temp_map_screengroupid[screengroupid]
	//}else{
	//	screenids,err = koala.GetAccess(screengroupid)
	//	if err != nil{
	//		log4go.Error(err.Error())
	//		return nil
	//	}
	//	temp_map_screengroupid[screengroupid] = screenids
	//}

	//if temp_map_subjectgroupid[subjectgroupid] != nil{
	//if temp_map_subjectgroupid[subjectgroupid] == nil{
	//	subjectids = temp_map_subjectgroupid[subjectgroupid]
	//	subject_type = temp_map_subjectgrouptype[subjectgroupid]
	//}else{
	//	subjectids,subject_type,err = koala.GetPerson(subjectgroupid)
	//	if err != nil{
	//		log4go.Error(err.Error())
	//		return nil
	//	}
	//	temp_map_screengroupid[screengroupid] = screenids
	//	temp_map_subjectgrouptype[subjectgroupid] = subject_type
	//}
	//subjectids,subject_type,err = koala.GetPerson(subjectgroupid)
	//if err != nil{
	//	log4go.Error(err.Error())
	//	return nil
	//}
	items_map:=map[string]interface{}{
		"subjects_group":subjectgroupid,
		"screens_group":screengroupid,
		"screens":screenids,
		"subjects":subjectids,
		"subject_type":subject_type,
		"start_time":item.Start_time,
		"end_time":item.End_time,
		"name":item.Name,
		"id":item_id,
	}
	log4go.Info(items_map)
	return items_map

}
func GetChildren2(item model.Statistic_Item)(map[string]interface{}){
	var(
		err error
		subject_type int
	)
	screenids := make([]int,0)
	screenids_all := make([]int,0)
	subjectids := make([]int,0)
	item_id := item.Id
	screengroupid := item.Screens_groups
	subjectgroupid := item.Subjects_group
	fmt.Println("---------------------------------------------get screengroupid",screengroupid)
	fmt.Println("---------------------------------------------get subjectgroupid",subjectgroupid)

	index := strings.Index(screengroupid,"-1")
	if index == -1{
		group_ids := strings.Split(screengroupid,",")
		if len(group_ids) > 0{
			for i:= 0;i<len(group_ids) ;i++  {
				group_id,_ := strconv.Atoi(group_ids[i])
				screenids,err = koala.GetAccess(group_id)
				if err != nil{
					log4go.Error(err.Error())
					return nil
				}
				temp_map_screengroupid[group_id] = screenids
				screenids_all = append(screenids_all,screenids...)
			}
		}
	} else{
		fmt.Println("---------------------------------------------get screengroupid is -1")
		screen_list ,err := koala.GetScreenList()
		if err!=nil{
			log4go.Error("get children failed, reason is :",err)
			return nil
		}
		fmt.Println("---------------------------------------------get screen_list")
		data_list,err := screen_list.Get("data").Array()
		if err!=nil{
			log4go.Error("screen list is not array, reason is :",err)
			return nil
		}
		if len(data_list) == 0 {
			log4go.Error("no screen")
			return nil
		}
		for _,v := range  data_list{
			id ,err := v.(interface{}).(map[string]interface{})["id"].(json.Number).Int()
			if err!=nil{
				log4go.Error("get screen id failed, reason is :",err)
				return nil
			}
			screenids = append(screenids,id)
		}
		temp_map_screengroupid[-1] = screenids
		screenids_all = append(screenids_all,screenids...)
	}
	fmt.Println("---------------------------------------------get subjectgroupid")
	if subjectgroupid == -1{
		temp_map_subjectgroupid[subjectgroupid] = subjectids
		temp_map_subjectgrouptype[subjectgroupid] = 0
	}else{
		subjectids,subject_type,err = koala.GetPerson(subjectgroupid)
		if err != nil{
			log4go.Error(err.Error())
			return nil
		}
		temp_map_subjectgroupid[subjectgroupid] = screenids
		temp_map_subjectgrouptype[subjectgroupid] = subject_type
	}
	items_map:=map[string]interface{}{
		"subjects_group":subjectgroupid,
		"screens_group":screengroupid,
		"screens":screenids_all,
		"subjects":subjectids,
		"subject_type":subject_type,
		"start_time":item.Start_time,
		"end_time":item.End_time,
		"name":item.Name,
		"id":item_id,
	}
	log4go.Info(items_map)
	return items_map

}



func (s *Service)GetStatistic(date string,)(model.Statistic , error){
	statistic,err := s.dao.GetStatistic(date)
	if err != nil{
		return model.Statistic{}, err
	}
	return statistic, nil
}

func (s *Service)GetStatistic2(date string,grandparent_map map[int][]model.Statistic_Item,grandparent []model.Statistic_Main)([]map[string]interface{} , error){

	var (
		//count int
		//reupdate = 0
		//counter int
	)
	resp_map_list := make([]map[string]interface{},0)
	//LOOP:for _,v := range grandparent{
	for _,v := range grandparent{
		//resp_map_list = make([]map[string]interface{},0)
		parentids := grandparent_map[v.Id]
		resp_map := make(map[string]interface{})
		resp_map_items := make([]map[string]interface{},0)

		resp_map["id"] = v.Id
		resp_map["name"] = v.Name

		for _,i := range parentids{

			statistics,err := s.Getonestatistic(date,v.Id,i.Id)
			if err != nil{
				log4go.Error("获取子项统计数出错,错误原因为:",err)
			}
			//fmt.Println(statistics.Number)
			resp_map_item := map[string]interface{}{
				"id":i.Id,
				"name":i.Name,
				"date":statistics.Date,
				"subjects_group":i.Subjects_group,
				"screens_group":i.Screens_groups,
				"start_time":i.Start_time,
				"end_time":i.End_time,
				"number":statistics.Number,
				"parent_id":i.Parent_id,
			}
			resp_map_items = append(resp_map_items,resp_map_item)

			//if statistics.Number == 0 {
			//	count +=1
			//}
		}
		//if count == len(resp_map_list){
		//	reupdate += 1
		//}
		resp_map["items"] = resp_map_items
		resp_map_list = append(resp_map_list, resp_map)
	}
	//for _,v := range resp_map_list{
	//	items_len := len(v["items"].([]map[string]interface{}))
	//	if items_len != 0{
	//
	//	}
	//}
	//counter += 1
	//fmt.Println(count)
	//fmt.Println(len(resp_map_list))
	//if reupdate != 0 && counter < 2{
	//	s.PrepareForstatistic2(date,nil)
	//	//time.Sleep(3*time.Second)
	//	fmt.Println("is null so we check again")
	//	goto LOOP
	//}

	return resp_map_list, nil
}


func (s *Service)Getonestatistic(date string,grandparent_id,parent_id int)(model.Statistics , error){
	statistics,err := s.dao.GetStatistic2(date,grandparent_id,parent_id)
	if err != nil{
		return model.Statistics{}, err
	}
	return statistics, nil
}



func (s *Service)GetTime()(*model.Time , error){
	var(
		t struct {
			Time *model.Time
		}
	)
	path := "./time.toml"
	toml.DecodeFile(path,&t)

	return t.Time,nil
}

func (s *Service)DesayunoCount(time1,time2 string,Position_map map[int]interface{})(int , error){
	EmployeeRecords := model.EmployeeRecords_json{
		Snap_begin_time:time1,
		Snap_end_time:time2,
		Page:1,
		Size:1000,
	}
	res, err, _, _ := s.GetIdentificationRecords_GroupBy(EmployeeRecords,Position_map)
	if err != nil{
		log4go.Error(err.Error())
		return 0 , errors.New("查询日间食堂统计数据失败")
	}
	count := len(res)
	return count,nil
}
func(s *Service) AlmuerzoCount(time1,time2 string,Position_map map[int]interface{})(int , error){
	EmployeeRecords := model.EmployeeRecords_json{
		Snap_begin_time:time1,
		Snap_end_time:time2,
		Page:1,
		Size:1000,
	}
	res, err, _, _ := s.GetIdentificationRecords_GroupBy(EmployeeRecords,Position_map)
	if err != nil{
		log4go.Error(err.Error())
		return 0 , errors.New("查询午间食堂统计数据失败")
	}
	count := len(res)
	return count,nil
}
func(s *Service) JantarCount(time1,time2 string,Position_map map[int]interface{})(int , error){
	EmployeeRecords := model.EmployeeRecords_json{
		Snap_begin_time:time1,
		Snap_end_time:time2,
		Page:1,
		Size:1000,
	}
	res, err, _, _ := s.GetIdentificationRecords_GroupBy(EmployeeRecords,Position_map)
	if err != nil{
		log4go.Error(err.Error())
		return 0 , errors.New("查询夜间食堂统计数据失败")
	}
	count := len(res)
	return count,nil
}

func (s *Service)GetBetweenDates(sdate, edate string) []string {
	d := []string{}
	timeFormatTpl := "2006-01-02 15:04:05"
	if len(timeFormatTpl) != len(sdate) {
		timeFormatTpl = timeFormatTpl[0:len(sdate)]
	}
	date, err := time.Parse(timeFormatTpl, sdate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	date2, err := time.Parse(timeFormatTpl, edate)
	if err != nil {
		// 时间解析，异常
		return d
	}
	if date2.Before(date) {
		// 如果结束时间小于开始时间，异常
		return d
	}
	// 输出日期格式固定
	timeFormatTpl = "2006-01-02"
	date2Str := date2.Format(timeFormatTpl)
	d = append(d, date.Format(timeFormatTpl))
	for {
		date = date.AddDate(0, 0, 1)
		dateStr := date.Format(timeFormatTpl)
		d = append(d, dateStr)
		if dateStr == date2Str {
			break
		}
	}
	return d
}
func (s *Service)UTCchange(UTC string)(string){
	//Bdate := "2020-07-18T23:00:00.000Z"//时间字符串
	//Bdate := "2020-07-18 23:00:00"//时间字符串
	if strings.Contains(UTC,"Z"){
		t, _ := time.Parse(time.RFC3339, UTC)

		local_time := t.In(time.Local)
		return local_time.String()[0:19]
	}else{
		return UTC
	}
	//t, err := time.ParseInLocation("2006-01-02 15:04", Bdate, time.Local)//t被转为本地时间的time.Time
	//t,err := time.Parse("2006-01-02 15:04", Bdate)//t被转为UTC时间的time.Time
}

func (s *Service) Statistic_Project_Create(name string)(results model.Statistic_Main,err error){

	results,err = s.dao.Statistic_Project_Create(name)
	if err!=nil{
		log4go.Error(err)
		return	results,err
	}
	return results,nil
}

func (s *Service) Statistic_Project_Item_Create(Parent_id int,Statistic_Item_Json model.Statistic_Item_Json)(err error){

	err = s.dao.Statistic_Project_Item_Create(Parent_id,Statistic_Item_Json)
	if err!=nil{
		log4go.Error(err)
		return	err
	}
	return nil
}

func (s *Service) Statistic_Project_Update(name string,parent_id int)(results model.Statistic_Main,err error){
	results,err = s.dao.Statistic_Project_Update(name,parent_id)
	if err!=nil{
		return	results,err
	}
	return results,nil
}

func (s *Service) Statistic_Project_Items_Get(Parent_id int)(Statistic_Items []model.Statistic_Item,err error){

	Statistic_Items , err = s.dao.Statistic_Project_Items_Get(Parent_id)
	if err!=nil{
		log4go.Error(err)
		return	Statistic_Items,err
	}
	return Statistic_Items,nil
}

func (s *Service) Statistic_Project_Delete(id int)(err error){

	err = s.dao.Statistic_Project_Delete(id)
	if err!=nil{
		log4go.Error(err)
		return	err
	}
	return nil
}
func (s *Service) Statistics_Project_Delete(id int)(err error){

	err = s.dao.Statistics_Project_Delete(id)
	if err!=nil{
		log4go.Error(err)
		return	err
	}
	return nil
}


func (s *Service) Statistic_Project_Items_Delete(Parent_id int)(err error){

	err = s.dao.Statistic_Project_Items_Delete(Parent_id)
	if err!=nil{
		log4go.Error(err)
		return	err
	}
	return nil
}


func (s *Service) Statistic_Project_Parent_Get()( []model.Statistic_Main, error){
	var (
		Parent []model.Statistic_Main
		err error
	)
	Parent,err = s.dao.Statistic_Project_Get()
	if err!=nil{
		log4go.Error(err)
		return	Parent,err
	}
	return Parent,nil
}

func (s *Service) Statistic_Project_Get(Persongrouplist_g,Accesslist_g map[int]interface{})( []map[string]interface{}, error){

	results := make([]map[string]interface{},0)
	Parent_results,err := s.Statistic_Project_Parent_Get()
	if err!=nil{
		log4go.Error(err)
		return	nil,errors.New("获取统计项失败")
	}
	for _,v := range Parent_results{
		Statistic_Items ,err := s.Statistic_Project_Items_Get(v.Id)
		if err!=nil{
			log4go.Error(err)
			s := fmt.Sprintf("获取%s子项统计项失败,编号为%d", v.Name,v.Id)
			return	nil,errors.New(s)
		}
		for k,v := range Statistic_Items{
			ids_str := strings.Split(v.Screens_groups,",")
			if len(ids_str) == 0 {
				continue
			}
			screens_name := make([]string,0)
			for i := range ids_str{
				group_id,_ := strconv.Atoi(ids_str[i])

				if value,_ := Accesslist_g[group_id];value!=nil{
					temp_screen_name := Accesslist_g[group_id].(string)
					screens_name = append(screens_name,temp_screen_name)
				}
				if group_id == -1{
					screens_name = append(screens_name,"全相机")
				}
			}
			all_screen_name :=strings.Join(screens_name,",")
			Statistic_Items[k].Screens_group_name = all_screen_name
			if Statistic_Items[k].Subjects_group == -1{
				Statistic_Items[k].Subjects_group_name = "全人员"
			}
			if value,_ := Persongrouplist_g[v.Subjects_group];value!=nil{
				Statistic_Items[k].Subjects_group_name = Persongrouplist_g[v.Subjects_group].(string)
			}
		}
		var res = map[string]interface{}{
			"id": v.Id,
			"name": v.Name,
			"items":Statistic_Items,
		}
		results = append(results, res)
	}
	return results,nil
}

func(s *Service)GetRecordById(screen_id,subject_id int)(result model.IdentificationRecord,err error){
	result,err = s.dao.GetIdentificationRecordById(screen_id,subject_id)
	return result,err
}

func(s *Service)TruncateTables()(err error){
   err = s.dao.TruncateTable()
   return err
}
func(s *Service)DeleteRecord(day int)(err error){
	err = s.dao.DeleteRecord(day)
	return err
}
