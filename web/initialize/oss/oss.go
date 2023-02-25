package oss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"go.uber.org/zap"
)

const (
	errorOSSNew               = "create oss new failed"
	errorGetBucket            = "get bucket failed"
	errorSetBucketTransferAcc = "set bucket transfer acc failed"
)

var OssBucket *oss.Bucket

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
	OssBucket = bucket
	return nil
}
