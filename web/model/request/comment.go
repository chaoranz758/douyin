package request

type CommentVideoRequest struct {
	Token       string `form:"token" binding:"required"`
	VideoID     string `form:"video_id" binding:"required"`
	CommentText string `form:"comment_text"`
	ActionType  string `form:"action_type" binding:"required,oneof=1 2"`
	CommentID   string `form:"comment_id"`
}

type GetCommentVideoListRequest struct {
	Token   string `form:"token" binding:"required"`
	VideoID string `form:"video_id" binding:"required"`
}
