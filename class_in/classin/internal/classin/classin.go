package classin

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"strconv"
	"time"
	log "github.com/alecthomas/log4go"
	"errors"
)

const(
	SID="33234258"
	SECRET="9gbEIchq"
	class_in_url="https://api.eeo.cn/partner/api/course.api.php?action="
)
var jar, _ = cookiejar.New(nil)
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
	fmt.Println(string(body))
	body = bytes.TrimPrefix(body, []byte("\xef\xbb\xbf"))
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

type Classin struct {}

var G_Class	=Classin{}

func (C *Classin)Md5V() string  {
	h := md5.New()
	string:=strconv.FormatInt(time.Now().Unix(),10)
	h.Write([]byte(SECRET+string))
	return hex.EncodeToString(h.Sum(nil))
}
func (C *Classin)RPC (action string,formdata map[string]interface{})(*simplejson.Json, error){
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, val := range formdata {
		_ = writer.WriteField(key, val.(string),"")
		//_ = writer.WriteField(key, val.(string))
	}
	err := writer.Close()
	if err != nil {
		log.Error(err)
		return  nil,err
	}
	log.Info(writer)
	req, err := http.NewRequest("POST", class_in_url+action, body)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	log.Debug(req.URL)
	log.Debug(req.Method)
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
	return resp_json,nil
}












//func GetAllMember() ([]map[string]interface{}, error) {
//	client := &http.Client{
//		Jar: jar,
//		Timeout:3*time.Second,
//	}
//
//
//
//
//
//	req, err := http.NewRequest("POST", getSchoolStudentListByPage_url, nil)
//	if err != nil {
//		log.Error(err)
//		return nil, errors.New("New request error")
//	}
//	log.Debug(req.URL)
//	log.Debug(req.Method)
//	resp, err := client.Do(req)
//	if err != nil {
//		log.Error(err.Error())
//		return nil, err
//	}
//
//	resp_json, err := doResponse(resp)
//	if err != nil {
//		log.Error(err.Error())
//		return nil, err
//	}
//	fmt.Println(resp_json)
//	var data = make([]map[string]interface{}, 0)
//	if value,err :=resp_json.Get("data").Array(); err == nil{
//		for _,mapdata := range value{
//			res := mapdata.(map[string]interface{})
//			data = append(data,res)
//		}
//	}
//	return data, nil
//}


