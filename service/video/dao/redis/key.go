package redis

const (
	KeyPrefix                 = "douyin:"
	KeyVPublishVideoInfo      = "v_publish_video_info:"
	KeyActivePublishVideoInfo = "active_publish_video_info:"
	KeyVideoInfoZSet          = "video_list_info_zset"
	KeyUserVideoCount         = "user_publish_video_count:"
	KeyVVideoID               = "v_video_id"
	KeyActiveVideoID          = "active_video_id"
	KeyVFavoriteVideo         = "v_favorite_video:"
	KeyActiveFavoriteVideo    = "active_favorite_video:"
)

func getKey(part string) string {
	return KeyPrefix + part
}
