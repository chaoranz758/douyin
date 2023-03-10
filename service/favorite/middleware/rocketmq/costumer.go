package rocketmq

import (
	"context"
	request1 "douyin/proto/user/request"
	request2 "douyin/proto/video/request"
	"douyin/service/favorite/dao/mysql"
	"douyin/service/favorite/dao/redis"
	"douyin/service/favorite/initialize/grpc_client"
	"douyin/service/favorite/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errorConnectToGRPCServer          = "connect to grpc server failed"
	errorPushVActiveFavoriteVideo     = "push v or active favorite video information failed"
	errorAddFavoriteCount             = "add favorite count failed"
	errorFavoriteRelation             = "create favorite relation failed"
	errorDeleteVActiveFavoriteVideo   = "delete v active favorite video information failed"
	errorSubFavoriteCount             = "sub favorite count failed"
	errorDeleteFavoriteRelation       = "delete favorite relation failed"
	errorAddUserFavoriteVideoCountSet = "add user favorite video count set failed"
	errorGetFavoriteCount             = "get favorite count failed"
	errorJsonUnmarshal                = "json unmarshal failed"
	errorExecuteWorkflow              = "result of workflow.Execute is"
)

const (
	successCostumer1 = "消费者1消息执行完毕"
	successCostumer2 = "消费者2消息执行完毕"
	successCostumer3 = "消费者3消息执行完毕"
	successCostumer4 = "消费者4消息执行完毕"
)

const (
	wfName1            = "workflow-favoriteVideoCustomer1"
	wfName2            = "workflow-favoriteVideoCustomer2"
	wfName3            = "workflow-favoriteVideoCustomer3"
	favoriteVideoCount = 4
)

type Producer1Message struct {
	LoginUserID       int64 `json:"loginUserID"`
	VideoID           int64 `json:"videoID"`
	IsV               bool  `json:"isV"`
	IsActive          bool  `json:"isActive"`
	LoginUserIsActive bool  `json:"loginUserIsActive"`
	LoginUserIsV      bool  `json:"loginUserIsV"`
	AuthorId          int64 `json:"authorId"`
}

type Producer2Message struct {
	LoginUserID int64 `json:"loginUserID"`
	VideoID     int64 `json:"videoID"`
	AuthorId    int64 `json:"authorId"`
}

type Producer3Message struct {
	LoginUserID       int64 `json:"loginUserID"`
	VideoID           int64 `json:"videoID"`
	LoginUserIsActive bool  `json:"loginUserIsActive"`
	LoginUserIsV      bool  `json:"loginUserIsV"`
	AuthorId          int64 `json:"authorId"`
	IsV               bool  `json:"isV"`
	IsActive          bool  `json:"isActive"`
}

type Producer4Message struct {
	LoginUserID int64 `json:"loginUserID"`
	VideoID     int64 `json:"videoID"`
	AuthorId    int64 `json:"authorId"`
}

func FavoriteCustomer1CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	wfName := wfName1 + shortuuid.New()
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var producer1Message Producer1Message
		if err := json.Unmarshal(msgs[0].Body, &producer1Message); err != nil {
			zap.L().Error(errorJsonUnmarshal, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideoRevert(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
				VideoId:           producer1Message.VideoID,
				IsV:               producer1Message.IsV,
				IsActive:          producer1Message.IsActive,
				LoginUserIsV:      producer1Message.LoginUserIsV,
				LoginUserIsActive: producer1Message.LoginUserIsActive,
				LoginUserId:       producer1Message.LoginUserID,
				AuthorId:          producer1Message.AuthorId,
			})
			if err != nil {
				if res2 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res2.Code == 2 {
					zap.L().Error(errorPushVActiveFavoriteVideo, zap.Error(err))
					return err
				}
			}
			//zap.L().Info("分支事务回滚")
			return nil
		})
		//读视频基本信息并将视频信息存入大V或活跃用户关注的基本视频信息redis
		res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideo(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
			VideoId:           producer1Message.VideoID,
			IsV:               producer1Message.IsV,
			IsActive:          producer1Message.IsActive,
			LoginUserIsActive: producer1Message.LoginUserIsActive,
			LoginUserIsV:      producer1Message.LoginUserIsV,
			LoginUserId:       producer1Message.LoginUserID,
			AuthorId:          producer1Message.AuthorId,
		})
		if err != nil {
			if res2 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return err
			}
			if res2.Code == 2 {
				zap.L().Error(errorPushVActiveFavoriteVideo, zap.Error(err))
				return err
			}
		}
		//zap.L().Info("分支事务执行完成")
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			//点赞关系写入Mysql点赞表
			var f = model.Favorite{
				VideoID: producer1Message.VideoID,
				UserID:  producer1Message.LoginUserID,
			}
			if err := mysql.CreateFavoriteRelation(&f); err != nil {
				zap.L().Error(errorFavoriteRelation, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis用户点赞视频数和视频点赞数加一
			_, err = redis.AddFavoriteCount(producer1Message.LoginUserID, producer1Message.VideoID, producer1Message.AuthorId)
			if err != nil {
				zap.L().Error(errorAddFavoriteCount, zap.Error(err))
				//点赞关系从Mysql点赞表删除
				if err := mysql.DeleteFavoriteRelation(producer1Message.VideoID, producer1Message.LoginUserID); err != nil {
					zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//zap.L().Info("本地事务执行完成")
			return nil, nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return consumer.ConsumeRetryLater, err
	}
	if err = workflow.Execute(wfName, shortuuid.New(), nil); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info(successCostumer1)
	return consumer.ConsumeSuccess, nil
}

func FavoriteCustomer2CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	wfName := wfName2 + shortuuid.New()
	err := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var producer2Message Producer2Message
		if err := json.Unmarshal(msgs[0].Body, &producer2Message); err != nil {
			zap.L().Error(errorJsonUnmarshal, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			res5, err := grpc_client.UserClientDtm.AddUserFavoriteVideoCountSetRevert(wf.Context, &request1.DouyinUserVideoCountSetRequest{
				UserId: producer2Message.LoginUserID,
			})
			if err != nil {
				if res5 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res5.Code == 2 {
					zap.L().Error(errorAddUserFavoriteVideoCountSet, zap.Error(err))
					return err
				}
			}
			//zap.L().Info("调用用户服务分支事务回滚完成")
			return nil
		})
		count, err := redis.GetFavoriteCount(producer2Message.LoginUserID)
		if err != nil {
			zap.L().Error(errorGetFavoriteCount, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		if count == favoriteVideoCount {
			res5, err := grpc_client.UserClientDtm.AddUserFavoriteVideoCountSet(wf.Context, &request1.DouyinUserVideoCountSetRequest{
				UserId: producer2Message.LoginUserID,
			})
			if err != nil {
				if res5 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res5.Code == 2 {
					zap.L().Error(errorAddUserFavoriteVideoCountSet, zap.Error(err))
					return err
				}
			}
		}
		//zap.L().Info("调用用户服务分支事务执行完成")
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			//点赞关系写入Mysql点赞表
			var f = model.Favorite{
				VideoID: producer2Message.VideoID,
				UserID:  producer2Message.LoginUserID,
			}
			if err := mysql.CreateFavoriteRelation(&f); err != nil {
				zap.L().Error(errorFavoriteRelation, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis用户点赞视频数和视频点赞数加一
			_, err := redis.AddFavoriteCount(producer2Message.LoginUserID, producer2Message.VideoID, producer2Message.AuthorId)
			if err != nil {
				zap.L().Error(errorAddFavoriteCount, zap.Error(err))
				//点赞关系从Mysql点赞表删除
				if err := mysql.DeleteFavoriteRelation(producer2Message.VideoID, producer2Message.LoginUserID); err != nil {
					zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//zap.L().Info("本地事务执行完成")
			return nil, nil
		})
		return err
	})
	if err != nil {
		return consumer.ConsumeRetryLater, err
	}
	//data, _ := json.Marshal(&msgs[0].Body)
	if err = workflow.Execute(wfName, shortuuid.New(), nil); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info(successCostumer2)
	return consumer.ConsumeSuccess, nil
}

func FavoriteCustomer3CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	wfName := wfName3 + shortuuid.New()
	errWorkflow := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var producer3Message Producer3Message
		if err := json.Unmarshal(msgs[0].Body, &producer3Message); err != nil {
			zap.L().Error(errorJsonUnmarshal, zap.Error(err))
			return status.Error(codes.Aborted, err.Error())
		}
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideo(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
				VideoId:           producer3Message.VideoID,
				IsV:               producer3Message.IsV,
				IsActive:          producer3Message.IsActive,
				LoginUserIsV:      producer3Message.LoginUserIsV,
				LoginUserIsActive: producer3Message.LoginUserIsActive,
				LoginUserId:       producer3Message.LoginUserID,
				AuthorId:          producer3Message.AuthorId,
			})
			if err != nil {
				if res2 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res2.Code == 2 {
					zap.L().Error(errorPushVActiveFavoriteVideo, zap.Error(err))
					return err
				}
			}
			//zap.L().Info("分支事务回滚成功")
			return nil
		})
		//将大V或活跃用户喜欢的基本视频信息从redis删除
		res2, err := grpc_client.VideoClientDtm.DeleteVActiveFavoriteVideo(wf.Context, &request2.DouyinDeleteVActiveFavoriteVideoRequest{
			UserId:            producer3Message.LoginUserID,
			VideoId:           producer3Message.VideoID,
			LoginUserIsV:      producer3Message.LoginUserIsV,
			LoginUserIsActive: producer3Message.LoginUserIsActive,
		})
		if err != nil {
			if res2 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return err
			}
			if res2.Code == 2 {
				zap.L().Error(errorDeleteVActiveFavoriteVideo, zap.Error(err))
				return err
			}
		}
		//zap.L().Info("分支事务执行成功")
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			//点赞关系从Mysql点赞表删除
			if err := mysql.DeleteFavoriteRelation(producer3Message.VideoID, producer3Message.LoginUserID); err != nil {
				zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis用户点赞视频数和视频点赞数减一
			if err := redis.SubFavoriteCount(producer3Message.LoginUserID, producer3Message.VideoID, producer3Message.AuthorId); err != nil {
				zap.L().Error(errorSubFavoriteCount, zap.Error(err))
				//点赞关系写入Mysql点赞表
				var f = model.Favorite{
					VideoID: producer3Message.VideoID,
					UserID:  producer3Message.LoginUserID,
				}
				if err := mysql.CreateFavoriteRelation(&f); err != nil {
					zap.L().Error(errorFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//zap.L().Info("本地事务执行完成")
			return nil, nil
		})
		return err
	})
	if errWorkflow != nil {
		return consumer.ConsumeRetryLater, errWorkflow
	}
	data, _ := json.Marshal(&msgs[0].Body)
	if err := workflow.Execute(wfName, shortuuid.New(), data); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info(successCostumer3)
	return consumer.ConsumeSuccess, nil
}

func FavoriteCustomer4CallBack(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	var producer4Message Producer4Message
	if err := json.Unmarshal(msgs[0].Body, &producer4Message); err != nil {
		zap.L().Error(errorJsonUnmarshal, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//fmt.Printf("%v\n", producer4Message)
	//点赞关系从Mysql点赞表删除
	if err := mysql.DeleteFavoriteRelation(producer4Message.VideoID, producer4Message.LoginUserID); err != nil {
		zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	//redis用户点赞视频数和视频点赞数减一
	if err := redis.SubFavoriteCount(producer4Message.LoginUserID, producer4Message.VideoID, producer4Message.AuthorId); err != nil {
		zap.L().Error(errorSubFavoriteCount, zap.Error(err))
		return consumer.ConsumeRetryLater, err
	}
	zap.L().Info(successCostumer4)
	return consumer.ConsumeSuccess, nil
}
