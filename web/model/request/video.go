package request

type GetVideoListRequest struct {
	LastTime string `form:"latest_time"`
	Token    string `form:"token"`
}

type PublishVideoRequest struct {
	Token string `form:"token" binding:"required"`
	Title string `form:"title" binding:"required"`
}

type GetPublishVideoRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token"`
}
