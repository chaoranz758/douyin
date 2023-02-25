package util

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
	"os"
)

const (
	errorOSSNew               = "create oss new failed"
	errorGetBucket            = "get bucket failed"
	errorOpenVideoData        = "open video data failed"
	errorPutVideoDataToOSS    = "put video data to oss failed"
	errorPutObject            = "put object failed"
	errorSetBucketTransferAcc = "set bucket transfer acc failed"
	errorOpenFile             = "open file"
)

var ossBucket *oss.Bucket

func InitOSS() error {
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := ""
	// 阿里云账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM用户进行API访问或日常运维，请登录RAM控制台创建RAM用户。
	accessKeyId := ""
	accessKeySecret := ""
	bucketName := ""
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 创建OSSClient实例。
	client, err := oss.New(endpoint, accessKeyId, accessKeySecret, oss.EnableCRC(false))
	if err != nil {
		zap.L().Error(errorOSSNew, zap.Error(err))
		return err
	}
	accConfig := oss.TransferAccConfiguration{}
	accConfig.Enabled = true
	err = client.SetBucketTransferAcc(bucketName, accConfig)
	if err != nil {
		zap.L().Error(errorSetBucketTransferAcc, zap.Error(err))
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		zap.L().Error(errorGetBucket, zap.Error(err))
		return err
	}
	ossBucket = bucket
	return nil
}

func Upload(videoURLLocal, videoName, imageName string) error {
	videoName1 := fmt.Sprintf("videos/%s.mp4", videoName)
	//fmt.Printf("文件名：%v\n", *videoData.Filename)
	//fmt.Printf("文件大小：%v\n", videoData.Size)
	videoData1, err := os.Open(videoURLLocal)
	if err != nil {
		zap.L().Error(errorOpenFile, zap.Error(err))
		return err
	}
	if err = ossBucket.PutObject(videoName1, videoData1); err != nil {
		zap.L().Error(errorPutVideoDataToOSS, zap.Error(err))
		return err
	}
	imageName1 := fmt.Sprintf("images/%s.jpg", imageName)
	videoURL := "https://simpledouyin.oss-cn-qingdao.aliyuncs.com/" + videoName1
	imageData1, err := GetSnapshot(videoURL, 1)
	if err != nil {
		zap.L().Error(errorGetSnapshot, zap.Error(err))
		return err
	}
	if err = ossBucket.PutObject(imageName1, imageData1); err != nil {
		zap.L().Error(errorPutObject, zap.Error(err))
		return err
	}
	return nil
}
