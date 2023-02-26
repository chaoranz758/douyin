package util

import "douyin/service/user/initialize/snowflake"

func GenID() int64 {
	return snowflake.Node.Generate().Int64()
}
