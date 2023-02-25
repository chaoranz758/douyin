package response

type CommentVideoResponse struct {
	StatusCode int32       `json:"status_code"`
	StatusMSG  string      `json:"status_msg"`
	Comment    CommentInfo `json:"comment"`
}

type GetCommentVideoListResponse struct {
	StatusCode  int32         `json:"status_code"`
	StatusMSG   string        `json:"status_msg"`
	CommentList []CommentInfo `json:"comment_list"`
}

type CommentInfo struct {
	ID         int64    `json:"id"`
	User       UserInfo `json:"user"`
	Content    string   `json:"content"`
	CreateDate string   `json:"create_date"`
}
