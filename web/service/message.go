package service

import (
	"context"
	request1 "douyin/proto/message/request"
	"douyin/proto/message/response"
	"douyin/web/client/grpc"
	"douyin/web/model/request"
	"go.uber.org/zap"
	"strconv"
)

const (
	errorSendMessage = "send message failed"
	errorGetMessage  = "get message failed"
)

func SendMessage(sendMessageRequest request.SendMessageRequest, loginUserId int64) error {
	actionType, _ := strconv.Atoi(sendMessageRequest.ActionType)
	toUserId1, _ := strconv.ParseInt(sendMessageRequest.ToUserID, 10, 64)
	res, err := grpc.MessageClient.SendMessage(context.Background(), &request1.DouyinMessageActionRequest{
		ToUserId:    toUserId1,
		ActionType:  int32(actionType),
		Content:     sendMessageRequest.Content,
		LoginUserId: loginUserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorSendMessage, zap.Error(err))
			return err
		}
	}
	return nil
}

func GetMessage(getMessageRequest request.GetMessageRequest, loginUserId int64) ([]*response.Message, error) {
	toUserId1, _ := strconv.ParseInt(getMessageRequest.ToUserID, 10, 64)
	res, err := grpc.MessageClient.GetMessage(context.Background(), &request1.DouyinMessageChatRequest{
		LoginUserId: loginUserId,
		ToUserId:    toUserId1,
		PreMsgTime:  getMessageRequest.PreMsgTime,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
	}
	return res.MessageList, nil
}
