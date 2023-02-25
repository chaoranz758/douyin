package request

type FavoriteVideoRequest struct {
	Token      string `form:"token" binding:"required"`
	VideoID    string `form:"video_id" binding:"required"`
	ActionType string `form:"action_type" binding:"required,oneof=1 2"`
}

type GetFavoriteVideoListRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token"`
}
