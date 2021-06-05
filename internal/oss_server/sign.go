package oss_server

import (
	"github.com/alecthomas/log4go"

	//"io"
	//"io/ioutil"
	//"os"
	//"strings"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

func HandleError(err error) {
	log4go.Info("oss err handle",err.Error())
}

const(
	region = "oss-cn-hangzhou"
	endpoint = "oss-cn-hangzhou.aliyuncs.com"
	AccessKeyId = "LTAI5tEKeNGQuXGcUCjCSMP3"
	AccessKeySecret = "VLUFD7uK3kdMtqF21TCKLBqVQ1fwoF"
	bucketName = "lawtest"
)
// http:// + region +"."+endpoint+filepath
func GetSignUrl(filepath string)(map[string]string,error){
	client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		HandleError(err)
		return nil,err
	}
	objectName := filepath
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		HandleError(err)
		return nil,err
	}
	// 签名直传。
	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, 600)
	if err != nil {
		HandleError(err)
		return nil,err
	}
	info := map[string]string{
		"region":region,
		"endpoint":endpoint,
		"accesskeyid":AccessKeyId,
		"accesskeysecret":AccessKeySecret,
		"bucketname":bucketName,
		"signedurl":signedURL,
		"signedurlwithoutsecuritytoken":"http://" +bucketName+"." +endpoint+"/"+filepath,
	}
	return info,nil
}


//func main() {
//	//region := "oss-cn-hangzhou"
//	//endpoint := "oss-cn-hangzhou.aliyuncs.com"
//	//AccessKeyId := "LTAI5tEKeNGQuXGcUCjCSMP3"
//	//AccessKeySecret := "VLUFD7uK3kdMtqF21TCKLBqVQ1fwoF"
//
//
//	client, err := oss.New(endpoint, AccessKeyId, AccessKeySecret)
//	if err != nil {
//		HandleError(err)
//	}
//
//	objectName := "minshi/3.jpeg"
//	// 获取存储空间。
//	bucket, err := client.Bucket(bucketName)
//	if err != nil {
//		HandleError(err)
//	}
//
//	// 签名直传。
//	signedURL, err := bucket.SignURL(objectName, oss.HTTPGet, 600)
//	if err != nil {
//		HandleError(err)
//	}
//	fmt.Println("signedURL",signedURL)
//
//
//}
