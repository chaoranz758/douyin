package handler

import (
	"context"
	"douyin/proto/message/api"
	"douyin/proto/message/request"
	"douyin/proto/message/response"
	"douyin/service/message/service"
	"go.uber.org/zap"
)

type Message struct {
	api.UnimplementedMessageServer
}

func (message *Message) SendMessage(ctx context.Context, req *request.DouyinMessageActionRequest) (*response.DouyinMessageActionResponse, error) {
	if err := service.SendMessage(req); err != nil {
		zap.L().Error(errorSendMessage, zap.Error(err))
		return &response.DouyinMessageActionResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinMessageActionResponse{
		Code: 1,
	}, nil
}

func (message *Message) GetMessage(ctx context.Context, req *request.DouyinMessageChatRequest) (*response.DouyinMessageChatResponse, error) {
	data, err := service.GetMessage(req)
	if err != nil {
		zap.L().Error(errorGetMessage, zap.Error(err))
		return &response.DouyinMessageChatResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinMessageChatResponse{
		Code:        1,
		MessageList: data,
	}, nil
}

func (message *Message) GetUserFriendMessage(ctx context.Context, req *request.DouyinGetUserFriendMessageRequest) (*response.DouyinGetUserFriendMessageResponse, error) {
	msg, msgType, err := service.GetUserFriendMessage(req)
	if err != nil {
		zap.L().Error(errorGetUserFriendMessage, zap.Error(err))
		return &response.DouyinGetUserFriendMessageResponse{
			Code: 2,
		}, err
	}
	return &response.DouyinGetUserFriendMessageResponse{
		Message: msg,
		MsgType: msgType,
		Code:    1,
	}, nil
}
