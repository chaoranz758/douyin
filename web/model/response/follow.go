package response

type FollowUserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
}

type GetFollowListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMSG  string     `json:"status_msg"`
	UserList   []UserInfo `json:"user_list"`
}

type GetFollowerListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMSG  string     `json:"status_msg"`
	UserList   []UserInfo `json:"user_list"`
}

type GetFriendListResponse struct {
	StatusCode int32      `json:"status_code"`
	StatusMSG  string     `json:"status_msg"`
	UserList   []UserInfo `json:"user_list"`
}
