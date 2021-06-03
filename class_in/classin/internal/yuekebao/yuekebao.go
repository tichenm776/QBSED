package yuekebao

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/alecthomas/log4go"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	//"strconv"
	"time"
)

const (

	GetAllMember_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=member"
	GetAllTeacher_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=teacher"
	GetAllcourse_catalog_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=course_catalog"
	GetAllcourse_category_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=course_category"
	GetAllcard_once_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=card_once"
	GetAllcourse_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=course"
	GetAllcard_tag_url = "https://www.yuekebao.cn/website/api.php?bid=10653&dataName=card_tag"
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
func GetAllMember() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
		req, err := http.NewRequest("GET", GetAllMember_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllTeacher() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET", GetAllTeacher_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllCourse() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET",
		GetAllcourse_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllcourse_category() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET", GetAllcourse_category_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllcourse_catalog() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET", GetAllcourse_catalog_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllcard_once() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET", GetAllcard_once_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}

func GetAllcard_tag() ([]map[string]interface{}, error) {
	client := &http.Client{
		Jar: jar,
		Timeout:3*time.Second,
	}
	req, err := http.NewRequest("GET", GetAllcard_tag_url, nil)
	if err != nil {
		log.Error(err)
		return nil, errors.New("New request error")
	}
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
	fmt.Println(resp_json)
	var data = make([]map[string]interface{}, 0)
	if value,err :=resp_json.Get("data").Array(); err == nil{
		for _,mapdata := range value{
			res := mapdata.(map[string]interface{})
			data = append(data,res)
		}
	}
	return data, nil
}