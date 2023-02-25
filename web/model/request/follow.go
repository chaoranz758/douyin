package request

type FollowUserRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserID   string `form:"to_user_id" binding:"required"`
	ActionType string `form:"action_type" binding:"required,oneof=1 2"`
}

type GetFollowListRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type GetFollowerListRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}

type GetFriendListRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}
