package request

type UserRegisterRequest struct {
	UserName string `form:"username" binding:"required,max=32"`
	Password string `form:"password" binding:"required,min=5,max=32"`
}

type UserLoginRequest struct {
	UserName string `form:"username" binding:"required,max=32"`
	Password string `form:"password" binding:"required,min=5,max=32"`
}

type GetUserInfoRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Token  string `form:"token" binding:"required"`
}
