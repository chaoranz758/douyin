package redis

const (
	KeyPrefix            = "douyin:"
	KeyVVideoComment     = "v_video_comment:"
	KeyVideoCommentCount = "comment_video_count:"
)

func getKey(part string) string {
	return KeyPrefix + part
}
