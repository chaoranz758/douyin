package redis

const (
	KeyPrefix               = "douyin:"
	KeyUserLoginCount       = "user_login_count:"
	KeyUserLoginSet         = "user_login_set"
	KeyVSet                 = "v_set"
	KeyActiveSet            = "active_set"
	KeyVUserInfo            = "user_v_info:"
	KeyActiveUserInfo       = "user_active_info:"
	KeyVFollowerInfo        = "user_v_follower_info:"
	KeyActiveFollowerInfo   = "user_active_follower_info:"
	KeyVFollowInfo          = "user_v_follow_info:"
	KeyActiveFollowInfo     = "user_active_follow_info:"
	KeyUserVideoSet         = "user_video_set"
	KeyUserFavoriteVideoSet = "user_favorite_video_set"
	KeyUserFollowCountSet   = "user_follow_count_set"
	KeyUserFollowerCountSet = "user_follower_count_set"
)

func getKey(part string) string {
	return KeyPrefix + part
}
