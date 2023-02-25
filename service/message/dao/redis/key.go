package redis

const (
	KeyPrefix                   = "douyin:"
	KeyMessageActiveUser        = "message_active_user"
	KeyMessageActiveUserMessage = "message_active_user_message:"
	KeyUserSendMessageCount     = "user_send_message_count:"
	KeyUserGetMessageCount      = "user_get_message_count:"
	KeyGetCountSet              = "user_get_count_set"
	KeySendCountSet             = "user_send_count_set"
)

func getKey(part string) string {
	return KeyPrefix + part
}
