package request

type SendMessageRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserID   string `form:"to_user_id" binding:"required"`
	ActionType string `form:"action_type" binding:"required,oneof=1"`
	Content    string `form:"content" binding:"required"`
}

type GetMessageRequest struct {
	Token      string `form:"token" binding:"required"`
	ToUserID   string `form:"to_user_id" binding:"required"`
	PreMsgTime int64  `form:"pre_msg_time"`
}
