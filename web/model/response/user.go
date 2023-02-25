package response

type UserRegisterResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type UserLoginResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
	UserID     int64  `json:"user_id"`
	Token      string `json:"token"`
}

type GetUserInfoResponse struct {
	StatusCode int32    `json:"status_code"`
	StatusMSG  string   `json:"status_msg"`
	User       UserInfo `json:"user"`
}

type UserInfo struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int64  `json:"follow_count"`
	FollowerCount int64  `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}
