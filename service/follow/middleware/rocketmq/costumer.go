package rocketmq

import (
	"context"
	request2 "douyin/proto/comment/request"
	request1 "douyin/proto/user/request"
	"douyin/proto/video/request"
	"douyin/service/follow/dao/mysql"
	"douyin/service/follow/dao/redis"
	"douyin/service/follow/initialize/grpc_client"
	"douyin/service/follow/initialize/rocketmq/producer"
	"douyin/service/follow/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.uber.org/zap"
)

const (
	errorConnectToGRPCServer           = "connect to grpc server failed"
	errorPushVActiveFollowerUserinfo   = "push v or active follower user information failed"
	errorAddUserFollowFollowerCount    = "add user follow follower count failed"
	errorCreateFollow                  = "create follow failed"
	errorDeleteVActiveFollowerUserinfo = "delete v active follower user information failed"
	errorSubUserFollowFollowerCount    = "sub user follow follower count failed"
	errorDeleteFollow                  = "delete follow failed"
	errorPushVSet                      = "push v user set failed"
	errorAddUserFollowUserCountSet     = "add user follow user count set failed"
	errorPushVActiveBasicInfoInit      = "push v or active basic info init failed"
	errorPushVCommentBasicInfoInit     = "push v comment basic info init failed"
	errorGetUserFollowList             = "get user follow list failed"
	errorJsonUnmarshal                 = "json unmarshal failed"
	errorSendMessage                   = "生产者3消息发送失败"
)

const (
	successCostumer1   = "消费者1消息执行完毕"
	successCostumer2   = "消费者2消息执行完毕"
	successCostumer3   = "消费者3消息执行完毕"
	successSendMessage = "生产者3消息发送成功"
)

const (
	topic3 = "followTopic3"
)

const (
	followCount    = 10
	vfollowerCount = 20
)

type ProducerMessage1 struct {
	LoginUserId   int64 `json:"loginUserId"`
	UserId        int64 `json:"userId"`
	IsV           bool  `json:"isV"`
	IsActive      bool  `json:"isActive"`
	LoginIsV      bool  `json:"loginIsV"`
	LoginIsActive bool  `json:"loginIsActive"`
}

type ProducerMessage2 struct {
	LoginUserId   int64 `json:"loginUserId"`
	UserId        int64 `json:"userId"`
	IsV           bool  `json:"isV"`
	IsActive      bool  `json:"isActive"`
	LoginIsV      bool  `json:"loginIsV"`
	LoginIsActive bool  `json:"loginIsActive"`
}

type ProducerMessage3 struct {
	UserId         int64   `json:"userId"`
	FollowIdList   []int64 `json:"followIdList"`
	FollowerIdList []int64 `json:"followerIdList"`
}

func FollowCustomer1CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer1Message ProducerMessage1
	if err := json.Unmarshal(msgs[0].Body, &producer1Message); err != nil {
		zap.L().Error(errorJsonUnmarshal, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	res1, err := grpc_client.UserClient.PushVActiveFollowerUserinfo(context.Background(), &request1.DouyinVActiveFollowFollowerUserinfoRequest{
		LoginUserId:      producer1Message.LoginUserId,
		UserId:           producer1Message.UserId,
		IsV:              producer1Message.IsV,
		IsActive:         producer1Message.IsActive,
		IsFollowFollower: 2,
		LoginIsActive:    producer1Message.LoginIsActive,
		LoginIsV:         producer1Message.LoginIsV,
	})
	if res1 == nil {
		zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	if res1.Code == 2 {
		zap.L().Error(errorPushVActiveFollowerUserinfo, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//关注关系写入mysql关注表
	var f = model.Follow{
		FollowID:   producer1Message.UserId,
		FollowerID: producer1Message.LoginUserId,
	}
	if err := mysql.CreateFollow(&f); err != nil {
		zap.L().Error(errorCreateFollow, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//redis用户关注数和粉丝数加一
	userFollowCount, userFollowerCount, err := redis.AddUserFollowFollowerCount(producer1Message.LoginUserId, producer1Message.UserId)
	if err != nil {
		zap.L().Error(errorAddUserFollowFollowerCount, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	if userFollowerCount == vfollowerCount && producer1Message.IsActive == true {
		fs1 := make([]model.Follow, 0)
		fs2 := make([]model.Follow, 0)
		if err = mysql.GetUserFollowList(&fs1, producer1Message.UserId); err != nil {
			zap.L().Error(errorGetUserFollowList, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		if err = mysql.GetUserFollowerList(&fs2, producer1Message.UserId); err != nil {
			zap.L().Error(errorGetUserFollowList, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		var followIdList []int64
		for i := 0; i < len(fs1); i++ {
			followIdList = append(followIdList, fs1[i].FollowID)
		}
		var followerIdList []int64
		for i := 0; i < len(fs2); i++ {
			followerIdList = append(followerIdList, fs2[i].FollowerID)
		}
		//走消息队列
		var producerMessage3 = ProducerMessage3{
			UserId:         producer1Message.UserId,
			FollowIdList:   followIdList,
			FollowerIdList: followerIdList,
		}
		data, _ := json.Marshal(producerMessage3)
		msg := &primitive.Message{
			Topic: topic3,
			Body:  data,
		}
		_, err = producer.Producer3.SendSync(context.Background(), msg)
		if err != nil {
			zap.L().Error(errorSendMessage, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		zap.L().Info(successSendMessage)
	}
	if userFollowCount == followCount {
		resAddUserFollowUserCountSet, err := grpc_client.UserClient.AddUserFollowUserCountSet(context.Background(), &request1.DouyinUserFollowCountSetRequest{
			UserId: producer1Message.LoginUserId,
		})
		if err != nil {
			if resAddUserFollowUserCountSet == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
			if resAddUserFollowUserCountSet.Code == 2 {
				zap.L().Error(errorAddUserFollowUserCountSet, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
		}
	}
	zap.L().Info(successCostumer1)
	return consumer.ConsumeSuccess, nil
}

func FollowCustomer2CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer2Message ProducerMessage2
	if err := json.Unmarshal(msgs[0].Body, &producer2Message); err != nil {
		zap.L().Error(errorJsonUnmarshal, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	res, err := grpc_client.UserClient.DeleteVActiveFollowerUserinfo(context.Background(), &request1.DouyinDeleteVActiveFollowFollowerUserinfoRequest{
		LoginUserId:      producer2Message.LoginUserId,
		UserId:           producer2Message.UserId,
		IsV:              producer2Message.IsV,
		IsActive:         producer2Message.IsActive,
		LoginIsActive:    producer2Message.LoginIsActive,
		LoginIsV:         producer2Message.LoginIsV,
		IsFollowFollower: 2,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		if res.Code == 2 {
			zap.L().Error(errorDeleteVActiveFollowerUserinfo, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
	}
	//关注关系从mysql关注表删除
	if err := mysql.DeleteFollow(producer2Message.LoginUserId, producer2Message.UserId); err != nil {
		zap.L().Error(errorDeleteFollow, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//redis用户关注数和粉丝数减一
	if err := redis.SubUserFollowFollowerCount(producer2Message.LoginUserId, producer2Message.UserId); err != nil {
		zap.L().Error(errorSubUserFollowFollowerCount, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info(successCostumer2)
	return consumer.ConsumeSuccess, err
}

func FollowCustomer3CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer3Message ProducerMessage3
	if err := json.Unmarshal(msgs[0].Body, &producer3Message); err != nil {
		zap.L().Error(errorJsonUnmarshal, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//！！！有一种可能性是用户由活跃用户升级为大V，这个时候不光要保存相关信息，还要将该用户之前在活跃用户那存的所有信息删掉
	//可能是掉下来再存的，实际在列表中有，别忘了加判断
	//在用户服务中存大V列表，大V用户信息、大V粉丝信息、大V关注的信息
	var isActive bool
	if len(producer3Message.FollowIdList) == 0 {
		resPushVSet1, err := grpc_client.UserClient.PushVUserRelativeInfoInit(context.Background(), &request1.DouyinPushVSetRequest{
			UserId:         producer3Message.UserId,
			FollowerIdList: producer3Message.FollowerIdList,
		})
		if err != nil {
			if resPushVSet1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
			if resPushVSet1.Code == 2 {
				zap.L().Error(errorPushVSet, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
		}
		if resPushVSet1.IsExist == 1 {
			//zap.L().Info("要写入大V的用户已经写入过了")
			return consumer.ConsumeSuccess, nil
		}
		isActive = resPushVSet1.IsActive
	} else {
		resPushVSet2, err := grpc_client.UserClient.PushVUserRelativeInfoInit(context.Background(), &request1.DouyinPushVSetRequest{
			UserId:         producer3Message.UserId,
			FollowerIdList: producer3Message.FollowerIdList,
			FollowIdList:   producer3Message.FollowIdList,
		})
		if err != nil {
			if resPushVSet2 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
			if resPushVSet2.Code == 2 {
				zap.L().Error(errorPushVSet, zap.Error(err))
				return consumer.ConsumeRetryLater, err
			}
		}
		if resPushVSet2.IsExist == 1 {
			return consumer.ConsumeSuccess, nil
		}
		isActive = resPushVSet2.IsActive
	}
	//在视频服务存大V的用户id与视频id哈希，大V点赞的视频信息，大V发布的视频信息
	res, err := grpc_client.VideoClient.PushVBasicInfoInit(context.Background(), &request.DouyinPushVInfoInitRequest{
		UserId:   producer3Message.UserId,
		IsActive: isActive,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		if res.Code == 2 {
			zap.L().Error(errorPushVActiveBasicInfoInit, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
	}
	//可能没有发布视频 处理
	if len(res.VideoIdList) == 0 {
		return consumer.ConsumeSuccess, nil
	}
	//在评论服务中保存大V发布的各视频的评论信息
	res1, err := grpc_client.CommentClient.PushVCommentBasicInfoInit(context.Background(), &request2.DouyinPushVCommentBasicInfoInitRequest{
		UserId:      producer3Message.UserId,
		VideoIdList: res.VideoIdList,
	})
	if err != nil {
		if res1 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
		if res1.Code == 2 {
			zap.L().Error(errorPushVCommentBasicInfoInit, zap.Error(err))
			return consumer.ConsumeRetryLater, err
		}
	}
	zap.L().Info(successCostumer3)
	return consumer.ConsumeSuccess, nil
}
