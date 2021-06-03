package model

import (
	"time"
)

type User struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name				string `gorm:"column:name" json:"name"`
	Created_at			int64 `gorm:"not null" json:"created_at"`
	Phone				string `gorm:"not null" json:"phone"`
	Password 			string `gorm:"not null" json:"password "`
	Role				int `gorm:"not null" json:"role "`		//律师1，管理2，司法局3 大数字权限兼容小数字
}

type Case struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name				string `gorm:"column:name" json:"name"`
	Created_at			int64 `gorm:"not null" json:"created_at"`
	Archiving_time 			int64 `gorm:"not null" json:"archiving_time "`
	Archiving_status  			bool `gorm:"not null" json:"archiving_status  "`
	User_id				string `gorm:"not null" json:"user_id"`
	Case_number  			string `gorm:"not null" json:"case_number"`
	Case_type 				int `gorm:"not null" json:"case_type"`    //1 民事 2刑事 3非诉 4顾问
	Money				float64 `gorm:"not null" json:"money"`
	Get_money  				bool `gorm:"not null" json:"get_money"`
}

type Files struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Case_id  			string `gorm:"not null" json:"case_id"`
}

type File struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Files_id  			int `gorm:"not null" json:"files_id"`
	File_number  			string `gorm:"not null" json:"file_number"`
	File_name  			string `gorm:"not null" json:"file_name"`
	File_url  			string `gorm:"not null" json:"file_url"`
}

type IdentificationRecord struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Event_type			int	`gorm:"column:event_type" json:"event_type"`
	Subject_type		int	`gorm:"column:subject_type" json:"subject_type"`
	Name				string `gorm:"column:name" json:"name"`
	Subject_id			int `gorm:"not null" json:"subject_id"`
	Screen_id			int `gorm:"not null" json:"screen_id"`
	Snap_position		string	`gorm:"column:snap_position" json:"snap_position"`
	Photo				string `gorm:"column:photo" json:"photo"`
	Snap_photo			string `gorm:"not null" json:"snap_photo"`
	Snap_time			string `gorm:"not null" json:"snap_time"`
	Date_time			time.Time `gorm:"not null" json:"date_time"`
	Recognition_times	int	`gorm:"column:recognition_times" json:"recognition_times"`
	Come_from			string `gorm:"not null" json:"come_from"`
	Remark			string `gorm:"not null" json:"remark"`
	Start_time			string `gorm:"not null" json:"start_time"`
	End_time			string `gorm:"not null" json:"end_time"`
}
type Stranger struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Event_type			int	`gorm:"column:event_type" json:"event_type"`
	Name				string	`gorm:"column:name" json:"name"`
	Subject_id			int	`gorm:"column:subject_id" json:"subject_id"`
	Screen_id			int `gorm:"not null" json:"screen_id"`
	Snap_position		string	`gorm:"column:snap_position" json:"snap_position"`
	Photo				string	`gorm:"column:photo" json:"photo"`
	Snap_photo			string `gorm:"not null" json:"snap_photo"`
	Snap_time			string `gorm:"not null" json:"snap_time"`
	Date_time			time.Time `gorm:"not null" json:"date_time"`
}
type Statistic struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Desayuno			int	`gorm:"column:desayuno" json:"desayuno"`
	Almuerzo			int	`gorm:"column:almuerzo" json:"almuerzo"`
	Jantar				int	`gorm:"column:jantar" json:"jantar"`
	Date				string `gorm:"not null" json:"date"`
}
type Time struct {
	Desayuno_begin      string	`toml:"desayuno_begin" json:"desayuno_begin"`
	Desayuno_end         string	`toml:"desayuno_end" json:"desayuno_end"`
	Almuerzo_begin      string	`toml:"almuerzo_begin" json:"almuerzo_begin"`
	Almuerzo_end  string	`toml:"almuerzo_end" json:"almuerzo_end"`
	Jantar_begin  string	`toml:"jantar_begin" json:"jantar_begin"`
	Jantar_end  string	`toml:"jantar_end" json:"jantar_end"`
	Screen_id int `toml:"jantar_end" json:"screen_id"`
}
type Statistic_Main struct {
	Id		int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name	string	`gorm:"column:name" json:"name"`
}
type Statistic_Item struct {

	Id		int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Parent_id	int	`gorm:"column:parent_id"json:"parent_id"`
	Name	string	`gorm:"column:name" json:"name"`
	Subjects_group	int	`gorm:"column:subjects_group" json:"subjects_group"`
	Subjects_group_name	string	`gorm:"column:subjects_group_name" json:"subjects_group_name"`
	Screens_group	int	`gorm:"column:screens_group" json:"screens_group"`
	Screens_groups	string	`gorm:"column:screens_groups" json:"screens_groups"`
	Screens_group_name	string	`gorm:"column:screens_group_name" json:"screens_group_name"`
	Start_time	string	`gorm:"column:start_time" json:"start_time"`
	End_time	string	`gorm:"column:end_time" json:"end_time"`
	//Number      int     `gorm:"column:number" json:"number"`
}
type Statistics struct {
	Id					int `gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Date				string `gorm:"not null" json:"date"`
	Parent_id	int	`gorm:"column:parent_id"json:"parent_id"`
	Grandparent_id	int	`gorm:"column:grandparent_id"json:"grandparent_id"`
	Subjects_group	int	`gorm:"column:subjects_group" json:"subjects_group"`
	Screens_group	int	`gorm:"column:screens_group" json:"screens_group"`
	Screens_groups	string	`gorm:"column:screens_groups" json:"screens_groups"`
	Start_time	string	`gorm:"column:start_time" json:"start_time"`
	End_time	string	`gorm:"column:end_time" json:"end_time"`
	Number      int     `gorm:"column:number" json:"number"`
}
type Statistical struct {
	Id int	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name string	 `gorm:"column:name" json:"name"`
}
type Statistical_item struct {
	Id int	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name string	 `gorm:"column:name" json:"name"`
	Screens_Group int	`gorm:"column:screens_group" json:"screens_group"`
	Subjects_Group int	`gorm:"column:subjects_Group" json:"subjects_Group"`
	Foreign_id int	`gorm:"column:foreign_id" json:"foreign_id"`
	Start_time string	`gorm:"column:start_time" json:"start_time"`
	End_time string	`gorm:"column:end_time" json:"end_time"`
}
type Statistic_list struct {
	Id int	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Name string	 `gorm:"column:name" json:"name"`
	Come_from			string `gorm:"not null" json:"come_from"`
	Remark			string `gorm:"not null" json:"remark"`
	Start_time			string `gorm:"not null" json:"start_time"`
	End_time			string `gorm:"not null" json:"end_time"`
	Enter_time			string `gorm:"not null" json:"enter_time"`
	Leave_time			string `gorm:"not null" json:"leave_time"`
	Remain_12			string `gorm:"not null" json:"remain_12"`
	Remain_24			string `gorm:"not null" json:"remain_24"`
}
type Member struct {
	//Id int	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Id int	`gorm:"primary_key;not null" json:"id"`
	Name string	 `gorm:"column:name" json:"name"`
	Phone			string `gorm:"not null" json:"phone"`
	Inert_time			string `gorm:"not null" json:"inert_time"`
	SID			string `gorm:"not null" json:"sid"`
	IsChange	int	   `gorm:"not null" json:"ischange"` //1 2
	Member_Type	string	   `gorm:"not null" json:"member_type"`
	Using	string	   `gorm:"not null" json:"using"` //teacher using or not
	//IsChange	int	   `gorm:"not null" json:"ischange"`
	//End_time			string `gorm:"not null" json:"end_time"`
	//Enter_time			string `gorm:"not null" json:"enter_time"`
	//Leave_time			string `gorm:"not null" json:"leave_time"`
	//Remain_12			string `gorm:"not null" json:"remain_12"`
	//Remain_24			string `gorm:"not null" json:"remain_24"`
}
type Course struct {
	//Id int	`gorm:"primary_key;AUTO_INCREMENT;not null" json:"id"`
	Id int	`gorm:"primary_key;not null" json:"id"`
	Course_id string	 `gorm:"column:courseid" json:"courseid"`
	Class_id string	 `gorm:"column:classid" json:"classid"`
	Name string	 `gorm:"column:name" json:"name"`
	Starttime			string `gorm:"not null" json:"starttime"`
	Endtime			string `gorm:"not null" json:"endtime"`
	Teacher_id			string `gorm:"not null" json:"teacherid"`
	Teacher_sid			string `gorm:"not null" json:"teachersid"`
	Student_id			string `gorm:"not null" json:"studentid"`
	Student_sid			string `gorm:"not null" json:"studentsid"`
	SeatNum			string `gorm:"not null" json:"seatnum"`

}