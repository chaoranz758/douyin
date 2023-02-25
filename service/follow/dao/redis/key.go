package redis

const (
	KeyPrefix            = "douyin:"
	KeyUserFollowCount   = "user_follow_count:"
	KeyUserFollowerCount = "user_follower_count:"
)

func getKey(part string) string {
	return KeyPrefix + part
}
