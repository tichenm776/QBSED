package model

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/alecthomas/log4go"
	xtime "github.com/bilibili/kratos/pkg/time"
	"io"
	"os"
	"sync"
	"zhiyuan/zyutil"
)


// Kratos hello kratos.
type KoalaPerson_ struct {
	Id               int    `json:"id"`
	Company_id       int    `json:"company_id"`
	Create_time      int64  `json:"create_time"`
	InterView_pinyin string `json:"interView_pinyin"`
	Visitor_type     string `json:"visitor_type"`
	Subject_type     int    `json:"subject_type"`
	Email            string `json:"email"`
	Password_reseted string `json:"password_reseted"`
	Name             string `json:"name"`
	Pinyin           string `json:"pinyin"`
	Gender           int    `json:"gender"`
	Photo_ids        string `json:"photo_ids"`
	Phone            string `json:"phone"`
	Avatar           string `json:"avatar"`
	Department       string `json:"department"`  //部门
	Title            string `json:"title"`       //职位
	Description      string `json:"description"` //签名
	Job_number       string `json:"job_number"`  //工号
	Remark           string `json:"remark"`
	Birthday         int64  `json:"birthday"`     //生日	时间戳（秒）
	Entry_date       int64  `json:"entry_date"`   //入职时间	时间戳（秒）
	Purpose          int64  `json:"purpose"`      //(访客属性) 来访目的	0: 其他, 1: 面试, 2: 商务, 3: 亲友, 4: 快递送货
	Interviewee      int64  `json:"interviewee"`  //(访客属性) 受访人
	Come_from        int64  `json:"come_from"`    //(访客属性) 来访单位
	Start_time       int64  `json:"start_time"`   //(访客属性) 预定来访时间	时间戳（秒）
	End_time         int64  `json:"end_time"`     //(访客属性) 预定离开时间	时间戳（秒）
	Visit_notify     string `json:"visit_notify"` //(访客属性) 来访是否发APP消息推送
}

type Koala_person struct {
	Avatar             string                `json:"avatar"`
	Birthday           int64                 `json:"birthday"`
	Come_from          string                `json:"come_from"`
	Company_id         int                   `json:"company_id"`
	Create_time        int64                 `json:"create_time"`
	Department         string                `json:"department"`
	Description        string                `json:"description"`
	Email              string                `json:"email"`
	End_time           int64                 `json:"end_time"`
	Entry_date         int64                 `json:"entry_date"`
	Gender             int                   `json:"gender"`
	Id                 int                   `json:"id"`
	Interviewee        string                `json:"interviewee"`
	Interviewee_pinyin string                `json:"interviewee_pinyin"`
	Job_number         string                `json:"job_number"`
	Name               string                `json:"name"`
	Password_reseted   bool                  `json:"password_reseted"`
	Phone              string                `json:"phone"`
	Photos             []Koala_person_photos `json:"photos"`
	Pinyin             string                `json:"pinyin"`
	Purpose            int                   `json:"purpose"`
	Remark             string                `json:"remark"`
	Start_time         int64                 `json:"start_time"`
	Subject_type       int                   `json:"subject_type"`
	Title              string                `json:"title"`
	Visit_notify       bool                  `json:"password_reseted"`
}
type Koala_person_photos struct {
	Origin_url string  `json:"origin_url"`
	Company_id int     `json:"company_id"`
	Id         int     `json:"id"`
	Quality    float32 `json:"quality"`
	Subject_id int     `json:"subject_id"`
	Url        string  `json:"url"`
	Version    int     `json:"version"`
}

type ServerConfig struct {
	Network      string         `dsn:"network"`
	Addr         string         `dsn:"address"`
	Timeout      xtime.Duration `dsn:"query.timeout"`
	ReadTimeout  xtime.Duration `dsn:"query.readTimeout"`
	WriteTimeout xtime.Duration `dsn:"query.writeTimeout"`
}

type FaceInfo struct {
	ObjectID    int32
	BoundingBox []int
	PosePitch   float32
	PoseRoll    float32
	Blur        float32
	Sex         string
	Age         int
	Minority    int
}

type FaceRecognizeInfo struct {
	FaceToken       string
	SearchScore     float32
	SearchThreshold float32
	PersonInfo      *FaceInfo
}

type FaceRecognizeGroup struct {
	GroupAlias string
}

type PersonInfo struct {
	Name            string
	Birthday        string
	Sex             string
	CertificateType string
	ID              string
	Country         string
	Province        string
	City            string
}

type EventCommInfo struct {
	Resolution     []int
	PictureType    int
	MachineAddress string
	SerialNo       string
}

type IDCardInfo struct {
	Name           string
	Sex            string
	Nation         string
	Number         string
	Address        string
	Office         string
	ValidTimeStart string
	ValidTimeStop  string
	ProfilePic     string
}

//type Addperson struct {
//	GroupID	int	`json:"GroupID"`
//	PersonInfo	interface{}	`json:"PersonInfo"`
//	ImageInfo interface{}	`json:"ImageInfo"`
//}

type Addperson struct {
	//GroupID	int	`json:"GroupID"`
	PersonInfo interface{} `json:"Person"`
	//ImageInfo interface{}	`json:"ImageInfo"`
}
type CreatePhoto struct {
	GroupID    int         `json:"GroupID"`
	PersonInfo interface{} `json:"PersonInfo"`
	ImageInfo  interface{} `json:"ImageInfo"`
}

type Rpcmodel struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type Rpcmodel_findperson struct {
	Id     int         `json:"id"`
	Method string      `json:"method"`
	Object int         `json:"object"`
	Params interface{} `json:"params"`
}

type NormalResponsemodel struct {
	Code    int         `json:"code"`
	Message string      `json:"err_msg"`
	Data    interface{} `json:"data"`
}

type RYPerson struct {
	Code int    `json:"Code"`
	Name string `json:"Name"`
	Sex  string `json:"Sex"`
	Type int    `json:"Type"`
	//Country				string	`json:"Country"`
	//Province					string	`json:"Province"`
	//City					string	`json:"City"`
	CertificateType string `json:"CertificateType"`
	GroupName       string `json:"GroupName"`
	Birthday        string `json:"Birthday"`
}
type RYPerson_2 struct {
	ID   string `json:"ID"`
	Name string `json:"Name"`
	Sex  string `json:"Sex"`
	//Type	int	`json:"Type"`
	Country         string `json:"Country"`
	Province        string `json:"Province"`
	City            string `json:"City"`
	CertificateType string `json:"CertificateType"`
	//GroupName				string	`json:"GroupName"`
	//Birthday	string	`json:"Birthday"`
}
type RYimg struct {
	Lengths [1]int64 `json:"Lengths"`
	Amount  int      `json:"Amount"`
}

type Pair struct {
	a, b, c interface{}
}

//统一的请求
type Person struct {
	Subject_type int    `json:"subject_type"`
	Subject_Id   int    `json:"subject_id"`
	Name         string `json:"name"`
	Photo        string `json:"photo"`
	Ip           string `json:"ip"`
}

//
type Resp4Device struct {
	Code    int         `json:"code"`
	Err_msg string      `json:"err_msg"`
	Data    interface{} `json:"data"`
	Page    interface{} `json:"Page"`
}

type EventPerson struct {
	PersonName string	`json:"personName"`
	SubjectId string	`json:"subjectId"`
	PersonType string	`json:"personType"`
	FacePicture string	`json:"facePicture"`
	Similarity string	`json:"similarity"`
	CameraIp string		`json:"cameraIp"`
	Timestamp int64		`json:"timestamp"`
}

type Access_Setting struct {
	Status int	`json:"status"`
	Name string	`json:"name"`
	Comment string	`json:"comment"`
	Subject_group_id int	`json:"subject_group_id"`
	Screen_group_ids []int	`json:"screen_group_ids"`
	Schedule_ids []int	`json:"schedule_ids"`
	Calendar_ids []int	`json:"calendar_ids"`
}

//type AccessScreen struct {
//	Screen_ids []int	`json:"screen_ids"`
//}

var G_map = Demo{
	Data: make(map[string]interface{},0),
	Lock: &sync.Mutex{},
	RLock:&sync.RWMutex{},
}
var G_map_time = Demo{
	Data: make(map[string]interface{},0),
	Lock: &sync.Mutex{},
	RLock:&sync.RWMutex{},
}
var Resp = make(map[string]interface{},0)

type Demo struct {
	Data map[string]interface{}
	Lock *sync.Mutex//golang struct 在赋值时是浅拷贝,会导致lock的不是同一份锁 所以需要加上指针操作
	RLock *sync.RWMutex//golang struct 在赋值时是浅拷贝,会导致lock的不是同一份锁 所以需要加上指针操作
}

func (d *Demo) Get(k string) interface{}{
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	d.RLock.RLock()
	defer d.RLock.RUnlock()
	return d.Data[k]
}

func (d *Demo) Set(k string,v interface{}) {
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	d.RLock.Lock()
	defer d.RLock.Unlock()
	d.Data[k]=v
	//log4go.Info("sety value",k,v)
}
func (d *Demo) Getmap()(map[string]interface{}) {
	defer zyutil.Recover()
	d.RLock.RLock()
	defer d.RLock.RUnlock()
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	copy_map := d.Data
	return copy_map
}
func (d *Demo) Clean() {
	d.RLock.Lock()
	defer d.RLock.Unlock()
	d.Data = nil //释放内存
	d.Data = make(map[string]interface{},0)
}
func (d *Demo) Delete(k string) {
	if d.Get(k) != nil{
		d.RLock.Lock()
		defer d.RLock.Unlock()
		delete(d.Data,k)
	}
}


type SubjectFile struct {
	Subjects	[]map[string]interface{}
}

func SaveFile(datetime string,datalist []map[string]interface{})(error){

	tofile :=SubjectFile{
		Subjects:datalist,
	}
	//file_path := fmt.Sprintf("F:/program/code/go/go_project/src/zhiyuan/koalamate_backup_server/subjects.toml")
	file_path := fmt.Sprintf("./date/"+datetime+".toml")
	//file_path := fmt.Sprintf("./%s_subjects.toml",version)
	file, err := os.Create(file_path)
	if err != nil {
		// failed to create/open the file
		log4go.Error("创建toml文件失败",err)
		return err
		//return false
	}
	if err := toml.NewEncoder(file).Encode(tofile); err != nil {
		// failed to encode
		log4go.Error("写入数据失败",err)
		//return false
		return err
	}
	return nil
}

func ReadFile2(datetime string)([]map[string]interface{},error){
	toSubjectFile := SubjectFile{}
	file_path := fmt.Sprintf("./date/"+datetime+".toml")
	subjectsbyte,err := ReadFile(file_path)
	if err != nil{
		log4go.Error("读取文件错误，原因：",err)
		return []map[string]interface{}{},err
	}
	//toSubjectFile := make([]map[string]interface{},0)
	if err := toml.Unmarshal(subjectsbyte,&toSubjectFile); err != nil {
		log4go.Error("读取数据失败",err)
		return []map[string]interface{}{},err
	}

	return toSubjectFile.Subjects,nil
}

func ReadFile(path string)([]byte,error){
	filesize,err := os.Stat(path)
	if err != nil{
		log4go.Error("判断文件大小出错：",err)
		return nil,err
	}
	ff := filesize.Size()
	if ff == 0{
		log4go.Info("file is nil")
		return nil, nil
	}
	f,err := os.Open(path)
	if err != nil{
		log4go.Error("读取文件出错：",err)
		return nil,err
	}
	defer f.Close()
	var chunk []byte

	buf := make([]byte,ff)
	for{
		n,err := f.Read(buf)
		if err != nil && err != io.EOF{
			log4go.Error("读取文件buf出错：",err)
			return nil, err
		}
		if n == 0{
			break
		}
		chunk = append(chunk,buf[:n]...)
	}
	return chunk,nil

}