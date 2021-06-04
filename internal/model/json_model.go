package model




type User_json struct {
	Name				string `gorm:"column:name" json:"name"`
	Phone				string `gorm:"not null" json:"phone"`
	Password 			string `gorm:"not null" json:"password "`
	Password_check 			string `gorm:"not null" json:"password_check "`
}











type Room_json struct {
	Room_name string	`json:"room_name"`
	Principal string	`json:"principal"`
	Principal_Phone string	`json:"principal_phone"`
	Remarks string	`json:"remarks"`
	Camera_Ids []int	`json:"camera_ids"`
}

type Add_Bedroom struct {
	Room_name string	`json:"room_name"`
	Principal string	`json:"principal"`
	Principal_phone string	`json:"principal_phone"`
	Remarks string	`json:"remarks"`
	Screen_bind_ids []int	`json:"screen_bind_ids"`
}

type Update_Bedroom struct {
	Room_name string	`json:"room_name"`
	Principal string	`son:"principal"`
	Principal_phone string	`json:"principal_phone"`
	Remarks string	`json:"remarks"`
	Screen_bind_ids []int	`json:"screen_bind_ids"`
}

type Modify_bed struct {
	Room_id		int		`json:"room_id"`
	Bed_name	string	`json:"bed_name"`
	Bed_status	string	`json:"bed_status"`
	Remarks	string	`json:"remarks"`
}
type Modify_bed2 struct {
	Room_id		int		`json:"room_id"`
	Bed_name	string	`json:"bed_name"`
	Bed_status	int	`json:"bed_status"`
	Remarks	string	`json:"remarks"`
}

//type Bed_status struct {
//	Bed_status	int	`json:"bed_status"`
//	Subject_id	int	`json:"subject_id"`
//}
type Subject_id struct {
	//Bed_status	int	`json:"bed_status"`
	Subject_id	int	`json:"subject_id"`
}
type Checktime struct {
	Checkin	string	`json:"check_in"`
	Checkout	string	`json:"check_out"`
}

type Bedroom_json struct {
	Room_name string	`json:"room_name"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}

type Normal_bed struct {
	Subject_name string	`json:"subject_name"`
	Bed_name string	`json:"bed_name"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}
type Abnormal_bed struct {

	Page int 	`json:"page"`
	Size int	`json:"size"`
}
type Getcheckinsubjects struct {
	Subject_name string	`json:"subject_name"`
	Checkin	string	`json:"check_in"`
	Checkout	string	`json:"check_out"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}


type Records_json struct {
	Room_id int	`json:"room_id"`
	//Bed_id int	`json:"bed_id"`
	Subject_name string	`json:"subject_name"`
	Bed_name string	`json:"bed_name"`
	Checkin	string	`json:"check_in"`
	Checkout	string	`json:"check_out"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}

type EmployeeRecords_json struct {
	//Snap_position string	`json:"snap_position"`
	//Bed_id int	`json:"bed_id"`
	Subject_id int	`json:"subject_id"`
	Screen_id int	`json:"screen_id"`
	Snap_position string	`json:"snap_position"`
	User_role int	`json:"user_role"`
	Name	string	`json:"name"`
	Snap_begin_time	string	`json:"snap_begin_time"`
	Snap_end_time	string	`json:"snap_end_time"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}

type StrangerRecords_json struct {
	//Snap_position string	`json:"snap_position"`
	//Bed_id int	`json:"bed_id"`
	Screen_id int	`json:"screen_id"`
	Snap_begin_time	string	`json:"snap_begin_time"`
	Snap_end_time	string	`json:"snap_end_time"`
	Page int 	`json:"page"`
	Size int	`json:"size"`
}


type Statistic_json struct {
	Snap_begin_time	string	`json:"snap_begin_time"`
	Snap_end_time	string	`json:"snap_end_time"`
}

type Statistic_Main_Json struct {
	Name	string	`json:"name"`
	Items   []Statistic_Item_Json `json:"items"`
}

type Statistic_Item_Json struct {
	Name	string	`json:"name"`
	Subjects_group	int	`json:"subjects_group"`
	Screens_group	int	`json:"screens_group"`
	Screens_groups	[]int	`json:"screens_groups"`
	Start_time	string	`json:"start_time"`
	End_time	string	`json:"end_time"`
}