package response

type SendMessageResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMSG  string `json:"status_msg"`
}

type GetMessageRequest struct {
	StatusCode  int32         `json:"status_code"`
	StatusMSG   string        `json:"status_msg"`
	MessageList []MessageInfo `json:"message_list"`
}

type MessageInfo struct {
	ID         int64  `json:"id"`
	Content    string `json:"content"`
	CreateTime string `json:"create_time"`
}
