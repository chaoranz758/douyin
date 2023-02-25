package util

import (
	"bytes"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"go.uber.org/zap"
	"os"
)

const (
	errorGetSnapshot = "get snap shot failed"
)

func GetSnapshot(videoPath string, frameNum int) (image *bytes.Buffer, err error) {
	buf := bytes.NewBuffer(nil)
	err = ffmpeg.Input(videoPath).Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).WithOutput(buf, os.Stdout).Run()
	if err != nil {
		zap.L().Error(errorGetSnapshot, zap.Error(err))
		return nil, err
	}
	return buf, err
}
