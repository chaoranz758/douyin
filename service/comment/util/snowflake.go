package util

import "douyin/service/comment/initialize/snowflake"

func GenID() int64 {
	return snowflake.Node.Generate().Int64()
}
