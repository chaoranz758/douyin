package oss

import (
	"douyin/web/initialize/config"
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
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	// 创建OSSClient实例。
	client, err := oss.New(config.Config.Oss.Endpoint, config.Config.Oss.AccessKeyId, config.Config.Oss.AccessKeySecret, oss.EnableCRC(false))
	if err != nil {
		zap.L().Error(errorOSSNew, zap.Error(err))
		return err
	}
	accConfig := oss.TransferAccConfiguration{}
	accConfig.Enabled = true
	err = client.SetBucketTransferAcc(config.Config.Oss.BucketName, accConfig)
	if err != nil {
		zap.L().Error(errorSetBucketTransferAcc, zap.Error(err))
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(config.Config.Oss.BucketName)
	if err != nil {
		zap.L().Error(errorGetBucket, zap.Error(err))
		return err
	}
	OssBucket = bucket
	return nil
}
