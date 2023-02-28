package util

import (
	"douyin/web/initialize/oss"
	"fmt"
	"go.uber.org/zap"
	"os"
)

const (
	errorPutVideoDataToOSS = "put video data to oss failed"
	errorPutObject         = "put object failed"
	errorOpenFile          = "open file"
)

const baseUrl = "https://simpledouyin.oss-cn-qingdao.aliyuncs.com/"

func Upload(videoURLLocal, videoName, imageName string) error {
	videoName1 := fmt.Sprintf("videos/%s.mp4", videoName)
	videoData1, err := os.Open(videoURLLocal)
	if err != nil {
		zap.L().Error(errorOpenFile, zap.Error(err))
		return err
	}
	if err = oss.OssBucket.PutObject(videoName1, videoData1); err != nil {
		zap.L().Error(errorPutVideoDataToOSS, zap.Error(err))
		return err
	}
	imageName1 := fmt.Sprintf("images/%s.jpg", imageName)
	videoURL := baseUrl + videoName1
	imageData1, err := GetSnapshot(videoURL, 1)
	if err != nil {
		zap.L().Error(errorGetSnapshot, zap.Error(err))
		return err
	}
	if err = oss.OssBucket.PutObject(imageName1, imageData1); err != nil {
		zap.L().Error(errorPutObject, zap.Error(err))
		return err
	}
	return nil
}
