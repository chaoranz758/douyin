package response

type GetVideoListResponse struct {
	StatusCode int32   `json:"status_code"`
	StatusMSG  string  `json:"status_msg"`
	NextTime   int64   `json:"next_time"`
	VideoList  []Video `json:"video_list"`
}

type PublishVideoResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
}

type GetPublishVideoResponse struct {
	StatusCode int32   `json:"status_code"`
	StatusMSG  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
}

type Video struct {
	ID            int64    `json:"id"`
	Author        UserInfo `json:"author"`
	PlayURL       string   `json:"play_url"`
	CoverURL      string   `json:"cover_url"`
	FavoriteCount int64    `json:"favorite_count"`
	CommentCount  int64    `json:"comment_count"`
	IsFavorite    bool     `json:"is_favorite"`
	Title         string   `json:"title"`
}
