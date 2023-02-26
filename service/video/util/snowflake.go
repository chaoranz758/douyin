package util

import "douyin/service/video/initialize/snowflake"

func GenID() int64 {
	return snowflake.Node.Generate().Int64()
}
