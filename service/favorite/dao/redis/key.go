package redis

const (
	KeyPrefix                  = "douyin:"
	KeyUserFavoriteVideoCount  = "user_favorite_video_count:"
	KeyVideoFavoriteCount      = "video_favorite_count:"
	KeyUserFavoritedTotalCount = "user_favorited_total_count:"
)

func getKey(part string) string {
	return KeyPrefix + part
}
