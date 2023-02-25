package service

import (
	"context"
	"douyin/proto/message/request"
	"douyin/proto/message/response"
	"douyin/service/message/dao/mysql"
	"douyin/service/message/dao/redis"
	"douyin/service/message/dao/rocketmq"
	"douyin/service/message/model"
	"douyin/service/message/pkg/snowflake"
	"encoding/json"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
	"sort"
	time2 "time"
)

const (
	errorIsMessageActiveUser    = "judge is message active user failed"
	errorCreateMessage          = "create message failed"
	errorJsonUnmarshal          = "json unmarshal failed"
	errorGetMessage             = "get message failed"
	errorAddUserSendMessage     = "add user send message count failed"
	errorAddUserGetMessage      = "add user get message count failed"
	errorGetMessageLast         = "get message last failed"
	errorJudgeGetMessageIsFirst = "judge get message is first failed"
)

type Messages []*response.Message

func (m Messages) Len() int {
	return len(m)
}

func (m Messages) Less(i, j int) bool {
	return m[i].CreateTime < m[j].CreateTime
}

func (m Messages) Swap(i, j int) {
	m[i], m[j] = m[j], m[i] //与上面的交换方法一致
}

type ProducerMessage1 struct {
	MessageId  int64  `json:"messageId"`
	UserId     int64  `json:"userId"`
	ToUserId   int64  `json:"toUserId"`
	Content    string `json:"content"`
	CreateTime int64  `json:"createTime"`
}

func SendMessage(req *request.DouyinMessageActionRequest) error {
	messageId := snowflake.GenID()
	time := time2.Now().UnixMilli()
	//查看登录用户是否是消息活跃用户
	b, err := redis.IsMessageActiveUser(req.LoginUserId)
	if err != nil {
		zap.L().Error(errorIsMessageActiveUser, zap.Error(err))
		return err
	}
	//是消息活跃用户
	if b == true {
		//进入消息队列：1.存入消息redis 2.存入mysql
		//走消息队列
		var producerMessage1 = ProducerMessage1{
			MessageId:  messageId,
			UserId:     req.LoginUserId,
			ToUserId:   req.ToUserId,
			Content:    req.Content,
			CreateTime: time,
		}
		data, _ := json.Marshal(producerMessage1)
		msg := &primitive.Message{
			Topic: "messageTopic1",
			Body:  data,
		}
		sync, err := rocketmq.Producer1.SendSync(context.Background(), msg)
		if err != nil {
			zap.L().Error("生产者1消息发送失败", zap.Error(err))
			return err
		}
		zap.L().Info("生产者1消息发送成功")
		fmt.Printf("生产者1发送的消息：%v\n", sync.String())
		return nil
	}
	//不是消息活跃用户
	//存入mysql
	var m = model.Message{
		MessageID: messageId,
		UserID:    req.LoginUserId,
		ToUserID:  req.ToUserId,
		Content:   req.Content,
		CreatedAt: time,
	}
	if err := mysql.CreateMessage(&m); err != nil {
		zap.L().Error(errorCreateMessage, zap.Error(err))
		return err
	}
	if err := redis.AddUserSendMessage(req.LoginUserId); err != nil {
		zap.L().Error(errorAddUserSendMessage, zap.Error(err))
		if err := mysql.DeleteMessage(messageId); err != nil {
			zap.L().Error("delete message failed", zap.Error(err))
			return err
		}
		return err
	}
	return nil
}

func GetMessage(req *request.DouyinMessageChatRequest) ([]*response.Message, error) {
	var isFirst bool
	if req.PreMsgTime == 0 {
		isFirst = true
	} else {
		isFirst = false
	}
	//查看登录用户是否是消息活跃用户
	b, err := redis.IsMessageActiveUser(req.LoginUserId)
	if err != nil {
		zap.L().Error(errorIsMessageActiveUser, zap.Error(err))
		return nil, err
	}
	//查看请求的用户是否是消息活跃用户
	b1, err := redis.IsMessageActiveUser(req.ToUserId)
	if err != nil {
		zap.L().Error(errorIsMessageActiveUser, zap.Error(err))
		return nil, err
	}
	//是消息活跃用户
	if b == true {
		var (
			rs Messages
		)
		var ps []ProducerMessage1
		result, err := redis.GetMessage(req.LoginUserId, req.ToUserId, req.PreMsgTime, isFirst)
		if err != nil {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
		if len(result) != 0 {
			if isFirst == true {
				for i := 0; i < len(result); i++ {
					var p ProducerMessage1
					if err := json.Unmarshal([]byte(result[i]), &p); err != nil {
						zap.L().Error(errorJsonUnmarshal, zap.Error(err))
						return nil, err
					}
					ps = append(ps, p)
				}
			} else {
				for i := 1; i < len(result); i++ {
					var p ProducerMessage1
					if err := json.Unmarshal([]byte(result[i]), &p); err != nil {
						zap.L().Error(errorJsonUnmarshal, zap.Error(err))
						return nil, err
					}
					ps = append(ps, p)
				}
			}
		}
		//var rs []*response.Message
		for i := 0; i < len(ps); i++ {
			var r = response.Message{
				Id:         ps[i].MessageId,
				Content:    ps[i].Content,
				CreateTime: ps[i].CreateTime,
				ToUserId:   ps[i].ToUserId,
				FromUserId: ps[i].UserId,
			}
			rs = append(rs, &r)
		}
		//请求的用户是活跃用户
		if b1 == true {
			var ps1 []ProducerMessage1
			result1, err := redis.GetMessage(req.ToUserId, req.LoginUserId, req.PreMsgTime, isFirst)
			if err != nil {
				zap.L().Error(errorGetMessage, zap.Error(err))
				return nil, err
			}
			//fmt.Printf("用户7消息1%v", result1)
			if len(result1) != 0 {
				if isFirst == true {
					for i := 0; i < len(result1); i++ {
						var p ProducerMessage1
						if err := json.Unmarshal([]byte(result1[i]), &p); err != nil {
							zap.L().Error(errorJsonUnmarshal, zap.Error(err))
							return nil, err
						}
						ps1 = append(ps1, p)
					}
				} else {
					for i := 1; i < len(result1); i++ {
						var p ProducerMessage1
						if err := json.Unmarshal([]byte(result1[i]), &p); err != nil {
							zap.L().Error(errorJsonUnmarshal, zap.Error(err))
							return nil, err
						}
						ps1 = append(ps1, p)
					}
				}
			}
			for i := 0; i < len(ps1); i++ {
				var r = response.Message{
					Id:         ps1[i].MessageId,
					Content:    ps1[i].Content,
					CreateTime: ps1[i].CreateTime,
					ToUserId:   ps1[i].ToUserId,
					FromUserId: ps1[i].UserId,
				}
				rs = append(rs, &r)
			}
		} else {
			ms1 := make([]model.Message, 0)
			if err := mysql.GetMessage(&ms1, req.ToUserId, req.LoginUserId, req.PreMsgTime, isFirst); err != nil {
				zap.L().Error(errorGetMessage, zap.Error(err))
				return nil, err
			}
			for i := 0; i < len(ms1); i++ {
				var r = response.Message{
					Id:         ms1[i].MessageID,
					Content:    ms1[i].Content,
					CreateTime: ms1[i].CreatedAt,
					ToUserId:   ms1[i].ToUserID,
					FromUserId: ms1[i].UserID,
				}
				rs = append(rs, &r)
			}
		}
		if err := redis.AddUserGetMessage(req.LoginUserId); err != nil {
			zap.L().Error(errorAddUserSendMessage, zap.Error(err))
			return nil, err
		}
		sort.Sort(rs)
		return rs, nil
	}
	//var rs []*response.Message
	var (
		rs Messages
	)
	if b1 == true {
		var ps1 []ProducerMessage1
		result1, err := redis.GetMessage(req.ToUserId, req.LoginUserId, req.PreMsgTime, isFirst)
		if err != nil {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
		//fmt.Printf("用户7消息2%v", result1)
		if len(result1) != 0 {
			if isFirst == true {
				for i := 0; i < len(result1); i++ {
					var p ProducerMessage1
					if err := json.Unmarshal([]byte(result1[i]), &p); err != nil {
						zap.L().Error(errorJsonUnmarshal, zap.Error(err))
						return nil, err
					}
					ps1 = append(ps1, p)
				}
			} else {
				for i := 1; i < len(result1); i++ {
					var p ProducerMessage1
					if err := json.Unmarshal([]byte(result1[i]), &p); err != nil {
						zap.L().Error(errorJsonUnmarshal, zap.Error(err))
						return nil, err
					}
					ps1 = append(ps1, p)
				}
			}
		}
		for i := 0; i < len(ps1); i++ {
			var r = response.Message{
				Id:         ps1[i].MessageId,
				Content:    ps1[i].Content,
				CreateTime: ps1[i].CreateTime,
				ToUserId:   ps1[i].ToUserId,
				FromUserId: ps1[i].UserId,
			}
			rs = append(rs, &r)
		}
		ms := make([]model.Message, 0)
		if err := mysql.GetMessage(&ms, req.LoginUserId, req.ToUserId, req.PreMsgTime, isFirst); err != nil {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
		for i := 0; i < len(ms); i++ {
			var r = response.Message{
				Id:         ms[i].MessageID,
				Content:    ms[i].Content,
				CreateTime: ms[i].CreatedAt,
				ToUserId:   ms[i].ToUserID,
				FromUserId: ms[i].UserID,
			}
			rs = append(rs, &r)
		}
		if err := redis.AddUserGetMessage(req.LoginUserId); err != nil {
			zap.L().Error(errorAddUserSendMessage, zap.Error(err))
			return nil, err
		}
		sort.Sort(rs)
		return rs, nil
	} else {
		ms := make([]model.Message, 0)
		if err := mysql.GetMessage(&ms, req.LoginUserId, req.ToUserId, req.PreMsgTime, isFirst); err != nil {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
		ms1 := make([]model.Message, 0)
		if err := mysql.GetMessage(&ms1, req.ToUserId, req.LoginUserId, req.PreMsgTime, isFirst); err != nil {
			zap.L().Error(errorGetMessage, zap.Error(err))
			return nil, err
		}
		for i := 0; i < len(ms); i++ {
			var r = response.Message{
				Id:         ms[i].MessageID,
				Content:    ms[i].Content,
				CreateTime: ms[i].CreatedAt,
				ToUserId:   ms[i].ToUserID,
				FromUserId: ms[i].UserID,
			}
			rs = append(rs, &r)
		}
		for i := 0; i < len(ms1); i++ {
			var r = response.Message{
				Id:         ms1[i].MessageID,
				Content:    ms1[i].Content,
				CreateTime: ms1[i].CreatedAt,
				ToUserId:   ms1[i].ToUserID,
				FromUserId: ms1[i].UserID,
			}
			rs = append(rs, &r)
		}
		if err := redis.AddUserGetMessage(req.LoginUserId); err != nil {
			zap.L().Error(errorAddUserSendMessage, zap.Error(err))
			return nil, err
		}
		sort.Sort(rs)
		return rs, nil
	}
}

func GetUserFriendMessage(req *request.DouyinGetUserFriendMessageRequest) ([]string, []int64, error) {
	var msg []string
	var msgType []int64
	for i := 0; i < len(req.ToUserId); i++ {
		//先查登录用户到登录用户的好友方向
		m := make([]model.Message, 0, 1)
		if err := mysql.GetMessageLast(&m, req.LoginUserId, req.ToUserId[i]); err != nil {
			zap.L().Error(errorGetMessageLast, zap.Error(err))
			return nil, nil, err
		}
		//再查另一个方向
		m1 := make([]model.Message, 0, 1)
		if err := mysql.GetMessageLast(&m1, req.ToUserId[i], req.LoginUserId); err != nil {
			zap.L().Error(errorGetMessageLast, zap.Error(err))
			return nil, nil, err
		}
		if len(m) == 0 && len(m1) == 0 {
			msg = append(msg, "")
			msgType = append(msgType, 0)
		}
		if len(m) == 0 && len(m1) != 0 {
			msg = append(msg, m1[0].Content)
			msgType = append(msgType, 0)
		}
		if len(m) != 0 && len(m1) == 0 {
			msg = append(msg, m[0].Content)
			msgType = append(msgType, 1)
		}
		if len(m) != 0 && len(m1) != 0 {
			if m[0].CreatedAt > m1[0].CreatedAt {
				msg = append(msg, m[0].Content)
				msgType = append(msgType, 1)
			} else {
				msg = append(msg, m1[0].Content)
				msgType = append(msgType, 0)
			}
		}
	}
	return msg, msgType, nil
}
