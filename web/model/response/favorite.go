package response

type FavoriteVideoResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
}

type GetFavoriteVideoListResponse struct {
	StatusCode int32   `json:"status_code"`
	StatusMSG  string  `json:"status_msg"`
	VideoList  []Video `json:"video_list"`
}
