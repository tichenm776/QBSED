package oss_server


//import (
//	"fmt"
//	"os"
//	"github.com/aliyun/aliyun-oss-go-sdk/oss"
//)
//func handleError(err error) {
//	fmt.Println("Error:", err)
//	os.Exit(-1)
//}
//func main() {
//	// Endpoint以杭州为例，其它Region请按实际情况填写。
//	endpoint := "http://oss-cn-hangzhou.aliyuncs.com"
//	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
//	accessKeyId := "<yourAccessKeyId>"
//	accessKeySecret := "<yourAccessKeySecret>"
//	bucketName := "<yourBucketName>"
//	// 创建OSSClient实例。
//	client, err := oss.New(endpoint, accessKeyId, accessKeySecret)
//	if err != nil {
//		handleError(err)
//	}
//	// 创建存储空间。
//	err = client.CreateBucket(bucketName)
//	if err != nil {
//		handleError(err)
//	}
//
//}
