// koala api
// 访问考拉服务器
package koala

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"
	"zhiyuan/QBSED/internal/model"
)

var baseUrl = ""
var koalaHost = ""
var jar, _ = cookiejar.New(nil)


//func InitLogin(){
//
//}
var G_koala = Koala{}

type Koala struct {

	Flag bool

}
func (k *Koala) Init(host string, koalaProt int)error{

	log4go.Info("koala init -----------------------------------------")
	//baseUrl = "http://" + host + ":" + strconv.Itoa(koalaProt)
	baseUrl = "http://" + host
	//baseUrl = "http://" + host
	log4go.Info("koala url: " + baseUrl)
	k.Flag =false
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	//jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		k.Flag =false
		//log4go.Crash(err)
		log4go.Error(err)
		return err
	}
	log4go.Info("koala init success -----------------------------------------")
	k.Flag =true
	return nil
}

func Init(url string) {
	//baseUrl = "http://" + config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	baseUrl = "http://"+url
	//baseUrl = "http://hz91zo.oicp.vip:10880"
	//baseUrl = "http://hz91zo.oicp.vip:59350"
	//koalaHost = "hz91zo.oicp.vip:10880"
	//koalaHost = config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	//log.Info("koala url: " + baseUrl)
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	//jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	/*
		var err error
		jar, err = cookiejar.New(nil)
		if err != nil {
			log.Crash(err)
		}

		if err := KoalaLogin(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword); err != nil {
			log.Crash(err)
		}
	*/

}
func Init2() {
	//baseUrl = "http://" + config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	//baseUrl = "http://"+url
	//baseUrl = "http://hz91zo.oicp.vip:10880"
	baseUrl = "http://192.168.18.51:80"
	//koalaHost = "hz91zo.oicp.vip:10880"
	//koalaHost = config.Gconf.KoalaHost + ":" + strconv.Itoa(config.Gconf.KoalaPort)
	log.Info("koala url: " + baseUrl)
	KoalaLogin("admin@91zo.com","admin123")
	// All users of cookiejar should import "golang.org/x/net/publicsuffix"
	//jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	/*
		var err error
		jar, err = cookiejar.New(nil)
		if err != nil {
			log.Crash(err)
		}

		if err := KoalaLogin(config.Gconf.KoalaUsername, config.Gconf.KoalaPassword); err != nil {
			log.Crash(err)
		}
	*/

}

//func doResponse(body *[]byte) error {
func doResponse(resp *http.Response) (*simplejson.Json, error) {
	log.Debug(resp.Status)
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("Read response body error")
	}
	log.Debug(string(body))

	jdata, err := simplejson.NewJson(body)
	if err != nil {
		log.Error(err.Error())
		return nil, errors.New("Face++返回报文错误")
	}

	code, _ := jdata.Get("code").Int()
	if code != 0 {
		desc, _ := jdata.Get("desc").String()
		log.Error(desc)
		return nil, errors.New(desc)
	}
	return jdata, nil
}

func KoalaLogin(username string, password string) (bool,error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}

	data := url.Values{}
	data.Set("username", username)
	data.Add("password", password)

	req, err := http.NewRequest("POST", baseUrl+"/auth/login", bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Error(err)
		return false,errors.New("New request error")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Koala Admin")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return false,err
	}

	res, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return false,err
	}
	codenum,_ := res.Get("code").Int()
	if codenum != 0{
		return false,nil
	}
	return true,nil

}
func GetEventsUser(startTime int64, endTime int64,page int64) ([]map[string]interface{},int64, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	start := strconv.FormatInt(startTime, 10)
	end := strconv.FormatInt(endTime, 10)
	page_num :=strconv.FormatInt(page, 10)
	log4go.Debug(" page is ----------------------------------------------")
	log4go.Debug(page)
	//req, err := http.NewRequest("GET", baseUrl+"/event/events?category=user&user_role=0&start="+start+"&end="+end+"&page="+page_num+"&size=1000", nil)
	req, err := http.NewRequest("GET", baseUrl+"/event/events?category=user&start="+start+"&end="+end+"&page="+page_num+"&size=1000", nil)
	if err != nil {
		log4go.Error(err)
		return nil,1, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil,1, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,1, err
	}

	arr, err := resp_json.Get("data").Array()
	pages := resp_json.Get("page")
	total,err:= pages.Get("total").Int64()
	//log.Info("adduni failed(%v)", arr2)
	if err != nil {
		log4go.Error(err.Error())
		return nil,1, err
	}

	var data = make([]map[string]interface{}, 0)
	for _, jdata := range arr {
		res := jdata.(map[string]interface{})
		data = append(data,res)
	}

	return data,total, nil
}
func AddPhoto(photo *multipart.File) (int, error) {
	client := &http.Client{
		Jar: jar,
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile3("photo", "photo.jpg","image/jpeg")
	if err != nil {
		return -1, err
	}
	_, err = io.Copy(part, *photo)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	//writer.WriteField("aa", "aa")

	err = writer.Close()
	if err != nil {
		return -1, err
	}

	request, err := http.NewRequest("POST", baseUrl+"/subject/photo", body)
	request.Header.Add("Content-Type", writer.FormDataContentType())
	log.Debug(request.URL)
	log.Debug(request.Method)
	log.Debug(request.Header)

	resp, err := client.Do(request)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	photo_id, _ := resp_json.Get("data").Get("id").Int()

	return photo_id, nil
}

func AddSubject(params *map[string]interface{}, photo *multipart.File) (*map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 上传底库照片
	photo_id, err := AddPhoto(photo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	jsdata := simplejson.New()
	for key, val := range *params {
		jsdata.Set(key, val)
	}
	photo_ids := []int{photo_id}
	jsdata.Set("photo_ids", photo_ids)
	byte_data, _ := jsdata.MarshalJSON()
	log.Debug(string(byte_data))

	// 新增人员
	req, err := http.NewRequest("POST", baseUrl+"/subject", strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var res_data = make(map[string]interface{})
	res_data["id"] = resp_json.Get("data").Get("id")
	res_data["photo_id"] = photo_id
	res_data["name"] = resp_json.Get("data").Get("name")

	return &res_data, nil
}

// 修改subject
func ModSubject(params *map[string]interface{}, photo *multipart.File) (*map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 上传底库照片
	photo_id, err := AddPhoto(photo)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	jsdata := simplejson.New()
	for key, val := range *params {
		jsdata.Set(key, val)
	}
	photo_ids := []int{photo_id}
	jsdata.Set("photo_ids", photo_ids)
	byte_data, _ := jsdata.MarshalJSON()
	log.Debug(string(byte_data))

	// 修改人员
	subject_id := (*params)["subject_id"].(int)
	req, err := http.NewRequest("PUT", baseUrl+"/subject/"+strconv.Itoa(subject_id), strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	var res_data = make(map[string]interface{})
	res_data["id"] = resp_json.Get("data").Get("id")
	res_data["photo_id"] = photo_id
	res_data["name"] = resp_json.Get("data").Get("name")

	return &res_data, nil
}

func DeleteSubject(subject_id int) error {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("DELETE", baseUrl+"/subject/"+strconv.Itoa(subject_id), nil)
	if err != nil {
		log.Error(err)
		return errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	_, err = doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func GetSubjects(category string) (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("GET", baseUrl+"/mobile-admin/subjects/list?category="+category+"&size=8000", nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp_json, nil
}

func GetDisplayDevice(device_token string) (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
	}

	// 新增人员
	req, err := http.NewRequest("GET", baseUrl+"/screen/get-display-config?device_token="+device_token, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp_json, nil
}

func GetDepartment() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", baseUrl+"/subject/department/list?size=10000", nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	arr, err := resp_json.Get("data").Get("department").Array()
	if err != nil {
		return nil, err
	}

	var data = make([]map[string]interface{},0)
	for k, jdata := range arr {
		res := jdata.(string)
		tmpdata := map[string]interface{}{
			"id":k+1,
			"department":res,
		}
		data = append(data, tmpdata)
	}


	return data, nil
}

func GetSubjectsByCondition(category string, name string, department string, page int, size int, koalaIP string) (map[string]interface{}, []map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", baseUrl+"/subject/list?category=" + category +
		"&name=" + name + "&department=" + department + "&size=" + strconv.Itoa(size) + "&page=" + strconv.Itoa(page), nil)
	if err != nil {
		log.Error(err)
		return nil, nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, nil, err
	}

	arr, err := resp_json.Get("data").Array()
	if err != nil {
		return nil, nil, err
	}

	var datas = make(map[string]interface{})
	pageInfo, err1 := resp_json.Get("page").Map()
	if err1 != nil {
		return nil, nil, err1
	}
	datas["count"] = pageInfo["count"]
	datas["current"] = pageInfo["current"]
	datas["size"] = pageInfo["size"]
	datas["total"] = pageInfo["total"]


	var data = make([]map[string]interface{}, 0)
	for _, jdata := range arr {
		res := jdata.(map[string]interface{})
		var res_data = make(map[string]interface{})
		res_data["id"] = res["id"]
		res_data["name"] = res["name"]
		res_data["description"] = res["description"]
		res_data["department"] = res["department"]
		if res["photos"].([]interface{})[0].(map[string]interface{})["url"] != "" {
			res_data["avatar"] = "http://" + koalaIP +res["photos"].([]interface{})[0].(map[string]interface{})["url"].(string)
		} else {
			res_data["avatar"] = ""
		}
		data = append(data, res_data)
	}


	return datas, data, nil
}

func GetSubjectPhotosById(id int) (*simplejson.Json, error){
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", baseUrl+"/subject/" + strconv.Itoa(id), nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	//req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return resp_json, nil
}

func GetScreenList() (*simplejson.Json, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	//baseUrl = "http://hz91zo.oicp.vip:59350"
	req, err := http.NewRequest("GET", baseUrl+"/system/screen", nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Add("Connection", "keep-alive");
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_json, nil
}

/**
获取单个考拉门禁
*/
func GetScreenById(id int) (map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", baseUrl+"/system/screen/" + strconv.Itoa(id), nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{})
	resp_data["id"], _ = resp_json.Get("data").Get("id").Int()
	//resp_data["allow_all_subjects"], _ = resp_json.Get("data").Get("allow_all_subjects").Bool()
	//resp_data["allow_visitor"], _ = resp_json.Get("data").Get("allow_visitor").Bool()
	resp_data["camera_position"], _ = resp_json.Get("data").Get("camera_position").String()
	resp_data["screen_token"],_ = resp_json.Get("data").Get("screen_token").String()

	return resp_data, nil
}

//3.2API 门禁
func GetScreenListV3() (map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}
	funcpath := "devices/screens/group/list"
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{})
	resp_data["id"], _ = resp_json.Get("data").Get("id").Int()
	//resp_data["allow_all_subjects"], _ = resp_json.Get("data").Get("allow_all_subjects").Bool()
	//resp_data["allow_visitor"], _ = resp_json.Get("data").Get("allow_visitor").Bool()
	resp_data["camera_position"], _ = resp_json.Get("data").Get("camera_position").String()
	resp_data["screen_token"],_ = resp_json.Get("data").Get("screen_token").String()

	return resp_data, nil
}

func GetTimer() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
	}
	//Init2()
	funcpath := "/access/schedule/list"
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make([]interface{},0)
	resp_data,ok := resp_json.Get("data").Array()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data_arr = make([]map[string]interface{},0)
	for _,v := range resp_data{
		//data_map,ok := interface{}(i).(map[string]interface{})
		if v, ok := interface{}(v).(map[string]interface{}); ok {
			resp_data_arr = append(resp_data_arr,v)
		}
	}
	return resp_data_arr, nil
}

func GetAccessList() ([]map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	//Init2()
	funcpath := "/devices/screens/group/list?page=1&size=10000"
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data = make([]interface{},0)
	resp_data,ok := resp_json.Get("data").Array()
	if ok != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data_arr = make([]map[string]interface{},0)
	for i := range resp_data{
		//data_map,ok := interface{}(i).(map[string]interface{})
		if v, ok := interface{}(resp_data[i]).(map[string]interface{}); ok {
			resp_data_arr = append(resp_data_arr,v)
		}
	}
	return resp_data_arr, nil
}

func GetAccess(id int) ([]int, error){
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	//Init2()
	///devices/screens/group/1?page=1&size=8000
	transform := strconv.Itoa(id)
	funcpath := "/devices/screens/group/"+transform+"?page=1&size=10000"
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{})
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	//var resp_data_arr = make([]map[string]interface{},0)
	//for i := range resp_data{
		//data_map,ok := interface{}(i).(map[string]interface{})
	if v, ok := interface{}(resp_data).(map[string]interface{}); ok {
		screenids := make([]int,0)
		if l,ok := v["screens"].([]interface{});ok{
			for k,_ := range l{
				id ,_ := l[k].(map[string]interface{})["id"].(json.Number).Int()
				screenids = append(screenids,id )
			}
		}else{
			return nil, err
		}
		return screenids, nil
	}
	return nil,err
}

func GetPerson(id int) ([]int,int, error){
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	//Init2()
	///devices/screens/group/1?page=1&size=8000
	transform := strconv.Itoa(id)
	funcpath := "/subjects/group/"+transform+"?page=1&size=8000"
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log4go.Error(err)
		return nil,-1, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil,-1, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil,-1, err
	}
	var resp_data = make(map[string]interface{})
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log4go.Error(err.Error())
		return nil,-1, err
	}
	//var resp_data_arr = make([]map[string]interface{},0)
	//for i := range resp_data{
	//data_map,ok := interface{}(i).(map[string]interface{})
	if v, ok := interface{}(resp_data).(map[string]interface{}); ok {
		subjectids := make([]int,0)
		subject_type,_ := v["subject_type"].(json.Number).Int()
		if l,ok := v["subjects"].([]interface{});ok{
			for k,_ := range l{
				id ,_ := l[k].(map[string]interface{})["id"].(json.Number).Int()
				subjectids = append(subjectids,id)
			}
		}else{
			return nil,-1, err
		}
		return subjectids,subject_type, nil
	}
	return nil,-1,err
}

func GetPersonData(resp_json *simplejson.Json)([]int,int, error){
	//fmt.Println(resp_json)
	//var resp_data = make(map[string]interface{})
	subjectids := make([]int,0)
	var subject_type int
	resp_data_,err := resp_json.Get("data").Array()
	if err != nil {
		fmt.Println(err)
		fmt.Println("get data err --------------------------------------")
		log4go.Error(err.Error())
		return nil,-1, err
	}
	//var resp_data_arr = make([]map[string]interface{},0)
	//for i := range resp_data{
	//data_map,ok := interface{}(i).(map[string]interface{})
	//fmt.Println(resp_json)
	//fmt.Println("--------------------------------------------------------------resp_json")
	for _,t := range resp_data_{
		//fmt.Println("for loop ")
		if v, ok := interface{}(t).(interface{}); ok {
			//fmt.Println("is interface")
			subject_type,err = v.(map[string]interface{})["subject_type"].(json.Number).Int()
			if err != nil {
				return nil, -1, err
			}
			id ,err:= v.(map[string]interface{})["id"].(json.Number).Int()
			if err != nil {
				return nil, -1, err
			}
			subjectids = append(subjectids, id)
		}
	}
	//fmt.Println(subjectids)
	return subjectids,subject_type, nil
}

func CreateAccess(name,comment string) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()


	jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	jsdata.Set("name", name)
	jsdata.Set("comment", comment)
	byte_data, _ := jsdata.MarshalJSON()
	log.Debug(string(byte_data))
	funcpath := "/devices/screens/group"
	req, err := http.NewRequest("PUT", baseUrl+funcpath, strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

func DeleteAccess(id int) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()


	jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	extra := map[string]interface{}{
		"screens_count":10000,
		"screens":"a/b/c",
	}

	jsdata.Set("extra", extra)
	//jsdata.Set("comment", comment)
	byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/devices/screens/group/"+strconv.Itoa(id)
	req, err := http.NewRequest("DELETE", baseUrl+funcpath, strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

//添加门禁到门禁分组
func InsertScreen2Access(id int,screenids []int ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}

	jsdata.Set("screen_ids", screenids)
	//jsdata.Set("comment", comment)
	byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/devices/screens/group/"+strconv.Itoa(id)+"/insert"
	req, err := http.NewRequest("POST", baseUrl+funcpath, strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}


//添加门禁到门禁分组
func DeleteScreen4Access(id int,screenids []int ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}

	jsdata.Set("screen_ids", screenids)
	//jsdata.Set("comment", comment)
	byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/devices/screens/group/"+strconv.Itoa(id)+"/delete"
	req, err := http.NewRequest("POST", baseUrl+funcpath, strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}



func GetPersonGroupList(sub_type string) ([]map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
		Timeout:30*time.Second,
	}
	//Init2()
	funcpath := "/subjects/group/list?page=1&size=10000&subject_type="+sub_type
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Info(req.URL)
	log4go.Info(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	log4go.Info(resp_json)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data = make([]interface{},0)
	resp_data,ok := resp_json.Get("data").Array()
	log4go.Info(resp_data)
	if ok != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data_arr = make([]map[string]interface{},0)
	for i := range resp_data{
		//data_map,ok := interface{}(i).(map[string]interface{})
		if v, ok := interface{}(resp_data[i]).(map[string]interface{}); ok {
			resp_data_arr = append(resp_data_arr,v)
		}
	}
	return resp_data_arr, nil
}

func CreatePersonGroup(name,comment string) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()


	jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	jsdata.Set("name", name)
	jsdata.Set("comment", comment)
	jsdata.Set("subject_type", 0)
	byte_data, _ := jsdata.MarshalJSON()
	log.Debug(string(byte_data))
	funcpath := "/subjects/group"
	req, err := http.NewRequest("PUT", baseUrl+funcpath, strings.NewReader(string(byte_data)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}
//koala无对外接口
//测试接口（无法删除）
func DeletePersonGroup(id int) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()
	//jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"

	//json := "[]"


	//extra := map[string]interface{}{
	//	"subject_type":"员工组",
	//	"subject_count":1,
	//}
	data_delete := map[string]interface{}{
		"id":id,
		//"extra":extra,
	}
	data_delete_list := []map[string]interface{}{data_delete}
	jsonBytes, err := json.Marshal(data_delete_list)
	//[{"id":3,"extra":{"subject_type":"员工组","subject_count":1}}]
	log.Info(string(jsonBytes))
	//jsdata.Set("", extra)
	//jsdata.Set("a", data_delete)
	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/access/subjects/list"
	req, err := http.NewRequest("DELETE", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

//添加门禁到门禁分组
func InsertPerson2PersonGroup(id int,subjectids []int ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	//jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}
	addid := map[string]interface{}{
		"subject_ids": subjectids,
	}

	jsonBytes, _ := json.Marshal(addid)
	//jsdata.Set("subject_ids", subjectids)
	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/subjects/group/"+strconv.Itoa(id)+"/insert"
	req, err := http.NewRequest("POST", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}


//删除人员从人员分组
func DeletePerson2PersonGroup(id int,subject_ids []int ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	//jsdata := simplejson.New()

	deleteid := map[string]interface{}{
		"subject_ids": subject_ids,
	}

	jsonBytes, _ := json.Marshal(deleteid)
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}

	//jsdata.Set("subject_ids", subject_ids)
	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	log.Info(string(jsonBytes))
	//log.Debug(string(byte_data))
	funcpath := "/subjects/group/"+strconv.Itoa(id)+"/delete"
	req, err := http.NewRequest("POST", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}



func CreateAccessControlRules(setting model.Access_Setting ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	//jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}
	jsonBytes, err := json.Marshal(setting)
	//jsdata.Set("subject_ids", subject_ids)
	log.Info(string(jsonBytes))
	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/access/setting"
	req, err := http.NewRequest("PUT", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

func UpdateAccessControlRules(id int , setting model.Access_Setting ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	//jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	//extra := map[string]interface{}{
	//	//	"screen_ids": screenids,
	//	//}
	jsonBytes, err := json.Marshal(setting)
	//jsdata.Set("subject_ids", subject_ids)
	log.Info(string(jsonBytes))
	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/access/setting/"+strconv.Itoa(id)
	req, err := http.NewRequest("POST", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

func DeleteAccessControlRules(id int ) (map[string]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	//Init2()

	//jsdata := simplejson.New()
	//name = "test"
	//comment = "lavida"
	deleteid := map[string]interface{}{
			"id": id,
		}

	jsonBytes, err := json.Marshal([]interface{}{deleteid})
	//jsdata.Set("subject_ids", subject_ids)

	//jsdata.Set("comment", comment)
	//byte_data, _ := jsdata.MarshalJSON()
	//log.Debug(string(byte_data))
	funcpath := "/access/setting/list"
	req, err := http.NewRequest("DELETE", baseUrl+funcpath, strings.NewReader(string(jsonBytes)))
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	log.Debug(req.URL)
	log.Debug(req.Method)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	var resp_data = make(map[string]interface{},0)
	resp_data,ok := resp_json.Get("data").Map()
	if ok != nil {
		log.Error(err.Error())
		return nil, err
	}
	return resp_data, nil
}

func GetPersonGroupListByGroupId(id int) ([]interface{}, error){
	client := &http.Client{
		Jar: jar,
	}
	Init2()
	funcpath := "/subjects/group/"+strconv.Itoa(id)
	req, err := http.NewRequest("GET", baseUrl+funcpath, nil)
	if err != nil {
		log4go.Error(err)
		return nil, errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data = make([]interface{},0)
	resp_data,ok := resp_json.Get("data").Get("subjects").Array()
	if ok != nil {
		log4go.Error(err.Error())
		return nil, err
	}
	var resp_data_arr = make([]interface{},0)
	for i := range resp_data{
		if _, ok := interface{}(i).(interface{}); ok {
			resp_data_arr = append(resp_data_arr,resp_data[i].(map[string]interface{})["id"])
		}
	}
	return resp_data_arr, nil
}
func Constants() (*simplejson.Json,error){
	client := &http.Client{
		Jar: jar,
	}
	fmt.Println(baseUrl)
	req, err := http.NewRequest("GET", baseUrl+"/event/constants?category=user", nil)
	if err != nil {
		log4go.Error(err)
		return nil,errors.New("New request error")
	}
	log4go.Debug(req.URL)
	log4go.Debug(req.Method)

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log4go.Error(err.Error())
		fmt.Println(err.Error())
		return nil,err
	}

	resp_json, err := doResponse(resp)
	if err != nil {
		log4go.Error(err.Error())
		fmt.Println(err.Error())
		return nil,err
	}
	//fmt.Println(resp_json)
	data_value := resp_json.Get("data")
	//if err != nil{
	//	log4go.Error(err.Error())
	//	fmt.Println(err.Error())
	//	return nil,err
	//}
	return data_value,nil
}