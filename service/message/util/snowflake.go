package util

import "douyin/service/message/initialize/snowflake"

func GenID() int64 {
	return snowflake.Node.Generate().Int64()
}
