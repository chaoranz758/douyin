package service

import (
	"context"
	"douyin/proto/follow/request"
	response1 "douyin/proto/follow/response"
	request2 "douyin/proto/message/request"
	request1 "douyin/proto/user/request"
	"douyin/proto/user/response"
	"douyin/service/follow/dao/mysql"
	"douyin/service/follow/dao/redis"
	"douyin/service/follow/initialize/grpc_client"
	"douyin/service/follow/initialize/rocketmq/producer"
	"douyin/service/follow/model"
	"encoding/json"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/dtm-labs/dtmcli"
	"github.com/dtm-labs/dtmgrpc/workflow"
	"github.com/lithammer/shortuuid/v3"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	wfName1 = "workflow-followUser"
)

const (
	topic1 = "followTopic1"
	topic2 = "followTopic2"
	topic3 = "followTopic3"
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

func FollowUserDtm(req *request.DouyinRelationActionRequest) error {
	if req.ActionType == 1 {
		//判断当前登录用户是否为大V或活跃用户
		res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
			UserId: req.LoginUserId,
		})
		if err != nil {
			if res == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return err
			}
			if res.Code == 2 {
				zap.L().Error(errorExeVActiveSet, zap.Error(err))
				return err
			}
		}
		//判断关注用户是否为大V或活跃用户
		res1, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
			UserId: req.ToUserId,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return err
			}
			if res1.Code == 2 {
				zap.L().Error(errorExeVActiveSet, zap.Error(err))
				return err
			}
		}
		wfName := wfName1 + shortuuid.New()
		errWF := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
			//大V或活跃用户分支事务
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				if res.IsInfluencer == true || res.IsActiver == true {
					res1, err := grpc_client.UserClientDtm.PushVActiveFollowerUserinfoRevert(wf.Context, &request1.DouyinVActiveFollowFollowerUserinfoRequest{
						LoginUserId:      req.LoginUserId,
						UserId:           req.ToUserId,
						IsV:              res1.IsInfluencer,
						IsActive:         res1.IsActiver,
						LoginIsV:         res.IsInfluencer,
						LoginIsActive:    res.IsActiver,
						IsFollowFollower: 1,
					})
					if err != nil {
						if res1 == nil {
							zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
							return err
						}
						if res1.Code == 2 {
							zap.L().Error(errorPushVActiveFollowerUserinfo, zap.Error(err))
							return err
						}
					}
				}
				//zap.L().Info("大V或活跃用户分支事务回滚完成")
				return nil
			})
			if res.IsInfluencer == true || res.IsActiver == true {
				res1, err := grpc_client.UserClientDtm.PushVActiveFollowerUserinfo(wf.Context, &request1.DouyinVActiveFollowFollowerUserinfoRequest{
					LoginUserId:      req.LoginUserId,
					UserId:           req.ToUserId,
					IsV:              res1.IsInfluencer,
					IsActive:         res1.IsActiver,
					LoginIsV:         res.IsInfluencer,
					LoginIsActive:    res.IsActiver,
					IsFollowFollower: 1,
				})
				if err != nil {
					if res1 == nil {
						zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
						return err
					}
					if res1.Code == 2 {
						zap.L().Error(errorPushVActiveFollowerUserinfo, zap.Error(err))
						return err
					}
				}
			}
			//zap.L().Info("大V或活跃用户分支事务执行完成")
			//用户粉丝集合分支事务
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				resAddUserFollowUserCountSet, err := grpc_client.UserClientDtm.AddUserFollowUserCountSetRevert(wf.Context, &request1.DouyinUserFollowCountSetRequest{
					UserId: req.LoginUserId,
				})
				if err != nil {
					if resAddUserFollowUserCountSet == nil {
						zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
						return err
					}
					if resAddUserFollowUserCountSet.Code == 2 {
						zap.L().Error(errorAddUserFollowUserCountSet, zap.Error(err))
						return err
					}
				}
				//zap.L().Info("用户粉丝集合分支事务回滚完成")
				return nil
			})
			userFollowCount, userFollowerCount, err := redis.GetUserFollowFollowerCountInner(req.LoginUserId, req.ToUserId)
			if err != nil {
				zap.L().Error(errorGetUserFollowFollowerCount, zap.Error(err))
				return status.Error(codes.Aborted, err.Error())
			}
			if res.IsActiver == false && res.IsInfluencer == false {
				if userFollowCount == 9 {
					resAddUserFollowUserCountSet, err := grpc_client.UserClientDtm.AddUserFollowUserCountSet(wf.Context, &request1.DouyinUserFollowCountSetRequest{
						UserId: req.LoginUserId,
					})
					if err != nil {
						if resAddUserFollowUserCountSet == nil {
							zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
							return err
						}
						if resAddUserFollowUserCountSet.Code == 2 {
							zap.L().Error(errorAddUserFollowUserCountSet, zap.Error(err))
							return err
						}
					}
				}
			}
			//zap.L().Info("用户粉丝集合分支事务执行完成")
			//用户关注者集合分支事务
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				resAddUserFollowerUserCountSet, err := grpc_client.UserClientDtm.AddUserFollowerUserCountSetRevert(wf.Context, &request1.DouyinUserFollowerCountSetRequest{
					UserId: req.ToUserId,
				})
				if err != nil {
					if resAddUserFollowerUserCountSet == nil {
						zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
						return err
					}
					if resAddUserFollowerUserCountSet.Code == 2 {
						zap.L().Error(errorAddUserFollowerUserCountSet, zap.Error(err))
						return err
					}
				}
				//zap.L().Info("用户关注者集合分支事务回滚完成")
				return nil
			})
			if res1.IsActiver == false && res1.IsInfluencer == false {
				println(userFollowerCount)
				if userFollowerCount == 9 {
					resAddUserFollowerUserCountSet, err := grpc_client.UserClientDtm.AddUserFollowerUserCountSet(wf.Context, &request1.DouyinUserFollowerCountSetRequest{
						UserId: req.ToUserId,
					})
					if err != nil {
						if resAddUserFollowerUserCountSet == nil {
							zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
							return err
						}
						if resAddUserFollowerUserCountSet.Code == 2 {
							zap.L().Error(errorAddUserFollowerUserCountSet, zap.Error(err))
							return err
						}
					}
				}
			}
			//zap.L().Info("用户关注者集合分支事务执行完成")
			_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
				if res1.IsInfluencer == true || res1.IsActiver == true {
					//走消息队列
					var producerMessage1 = ProducerMessage1{
						LoginUserId:   req.LoginUserId,
						UserId:        req.ToUserId,
						IsV:           res1.IsInfluencer,
						IsActive:      res1.IsActiver,
						LoginIsV:      res.IsInfluencer,
						LoginIsActive: res.IsActiver,
					}
					data, _ := json.Marshal(producerMessage1)
					msg := &primitive.Message{
						Topic: topic1,
						Body:  data,
					}
					_, err = producer.Producer1.SendSync(context.Background(), msg)
					if err != nil {
						zap.L().Error(errorSendMessage1, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info(successSendMessage1)
					return nil, nil
				}
				//关注关系写入mysql关注表
				var f = model.Follow{
					FollowID:   req.ToUserId,
					FollowerID: req.LoginUserId,
				}
				if err := mysql.CreateFollow(&f); err != nil {
					zap.L().Error(errorCreateFollow, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//redis用户关注数和粉丝数加一
				_, userFollowerCount, err := redis.AddUserFollowFollowerCount(req.LoginUserId, req.ToUserId)
				if err != nil {
					zap.L().Error(errorAddUserFollowFollowerCount, zap.Error(err))
					if err := mysql.DeleteFollow(req.LoginUserId, req.ToUserId); err != nil {
						zap.L().Error(errorDeleteFollow, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, status.Error(codes.Aborted, err.Error())
				}
				if userFollowerCount == 20 {
					fs1 := make([]model.Follow, 0)
					fs2 := make([]model.Follow, 0)
					if err = mysql.GetUserFollowList(&fs1, req.ToUserId); err != nil {
						zap.L().Error(errorGetUserFollowList, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					if err = mysql.GetUserFollowerList(&fs2, req.ToUserId); err != nil {
						zap.L().Error(errorGetUserFollowList, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
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
						UserId:         req.ToUserId,
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
						zap.L().Error(errorSendMessage3, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info(successSendMessage3)
					return nil, nil
				}
				//zap.L().Info("本地事务执行完成")
				return nil, nil
			})
			return err
		})
		if errWF != nil {
			zap.L().Error(errorRegisterWorkflow, zap.Error(err))
			return err
		}
		if err = workflow.Execute(wfName, shortuuid.New(), nil); err != nil {
			zap.L().Error(errorExecuteWorkflow, zap.Error(err))
			return err
		}
		//zap.L().Info("关注用户dtm事务执行完成")
		return nil
	}
	//判断当前登录用户是否为大V或活跃用户
	res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.LoginUserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return err
		}
	}
	//判断关注用户是否为大V或活跃用户
	res3, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.ToUserId,
	})
	if err != nil {
		if res3 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res3.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return err
		}
	}
	wfName := wfName1 + shortuuid.New()
	errWorkflow := workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		//删除大V活跃用户粉丝信息分支事务
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if res.IsInfluencer == true || res.IsActiver == true {
				//删除redis对应部分
				res1, err := grpc_client.UserClientDtm.PushVActiveFollowerUserinfo(wf.Context, &request1.DouyinVActiveFollowFollowerUserinfoRequest{
					LoginUserId:      req.LoginUserId,
					UserId:           req.ToUserId,
					IsV:              res3.IsInfluencer,
					IsActive:         res3.IsActiver,
					LoginIsV:         res.IsInfluencer,
					LoginIsActive:    res.IsActiver,
					IsFollowFollower: 1,
				})
				if err != nil {
					if res1 == nil {
						zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
						return err
					}
					if res1.Code == 2 {
						zap.L().Error(errorDeleteVActiveFollowerUserinfo, zap.Error(err))
						return err
					}
				}
			}
			//zap.L().Info("删除大V活跃用户粉丝信息分支事务回滚")
			return nil
		})
		if res.IsInfluencer == true || res.IsActiver == true {
			//删除redis对应部分
			res1, err := grpc_client.UserClientDtm.DeleteVActiveFollowerUserinfo(wf.Context, &request1.DouyinDeleteVActiveFollowFollowerUserinfoRequest{
				LoginUserId:      req.LoginUserId,
				UserId:           req.ToUserId,
				IsV:              res3.IsInfluencer,
				IsActive:         res3.IsActiver,
				LoginIsV:         res.IsInfluencer,
				LoginIsActive:    res.IsActiver,
				IsFollowFollower: 1,
			})
			if err != nil {
				if res1 == nil {
					zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
					return err
				}
				if res1.Code == 2 {
					zap.L().Error(errorDeleteVActiveFollowerUserinfo, zap.Error(err))
					return err
				}
			}
		}
		//zap.L().Info("删除大V活跃用户粉丝信息分支事务执行成功")
		//本地事务
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			if res3.IsInfluencer == true || res3.IsActiver == true {
				//走消息队列
				var producerMessage2 = ProducerMessage2{
					LoginUserId:   req.LoginUserId,
					UserId:        req.ToUserId,
					IsV:           res3.IsInfluencer,
					IsActive:      res3.IsActiver,
					LoginIsV:      res.IsInfluencer,
					LoginIsActive: res.IsActiver,
				}
				data, _ := json.Marshal(producerMessage2)
				msg := &primitive.Message{
					Topic: topic2,
					Body:  data,
				}
				_, err = producer.Producer2.SendSync(context.Background(), msg)
				if err != nil {
					zap.L().Error(errorSendMessage2, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				zap.L().Info(successSendMessage2)
				return nil, nil
			}
			//关注关系从mysql关注表删除
			if err := mysql.DeleteFollow(req.LoginUserId, req.ToUserId); err != nil {
				zap.L().Error(errorDeleteFollow, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis用户关注数和粉丝数减一
			if err := redis.SubUserFollowFollowerCount(req.LoginUserId, req.ToUserId); err != nil {
				zap.L().Error(errorSubUserFollowFollowerCount, zap.Error(err))
				//关注关系写入mysql关注表
				var f = model.Follow{
					FollowID:   req.ToUserId,
					FollowerID: req.LoginUserId,
				}
				if err := mysql.CreateFollow(&f); err != nil {
					zap.L().Error(errorCreateFollow, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//zap.L().Info("本地事务执行成功")
			return nil, nil
		})
		return err
	})
	if errWorkflow != nil {
		zap.L().Error(errorRegisterWorkflow, zap.Error(errWorkflow))
		return err
	}
	if err = workflow.Execute(wfName, shortuuid.New(), nil); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return err
	}
	//zap.L().Info("删除关注用户dtm事务执行完成")
	return nil
}

func GetFollowList(req *request.DouyinRelationFollowListRequest) ([]*response.User, error) {
	//判断用户是否为大V或活跃用户
	res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
	}
	if res.IsInfluencer == true || res.IsActiver == true {
		//根据判断条件选择是从redis还是mysql读用户信息
		res1, err := grpc_client.UserClient.GetVActiveFollowerUserinfo(context.Background(), &request1.DouyinGetVActiveFollowFollowerUserinfoRequest{
			LoginUserId:      req.LoginUserId,
			UserId:           req.UserId,
			IsV:              res.IsInfluencer,
			IsActive:         res.IsActiver,
			IsFollowFollower: 1,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetVActiveFollowerUserinfo, zap.Error(err))
				return nil, err
			}
		}
		return res1.User, nil
	}
	//先从mysql读出关注列表，再根据用户id批量从用户服务读用户信息
	fs := make([]model.Follow, 0)
	if err := mysql.GetUserFollowList(&fs, req.UserId); err != nil {
		zap.L().Error(errorGetUserFollowList, zap.Error(err))
		return nil, err
	}
	if len(fs) == 0 {
		return nil, nil
	}
	var idList []int64
	for i := 0; i < len(fs); i++ {
		idList = append(idList, fs[i].FollowID)
	}
	res2, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
		UserId:      idList,
		LoginUserId: req.LoginUserId,
	})
	if err != nil {
		if res2 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res2.Code == 2 {
			zap.L().Error(errorGetUserInfoList, zap.Error(err))
			return nil, err
		}
	}
	return res2.User, nil
}

func GetFollowerList(req *request.DouyinRelationFollowerListRequest) ([]*response.User, error) {
	//判断用户是否为大V或活跃用户
	res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
	}
	if res.IsInfluencer == true || res.IsActiver == true {
		//根据判断条件选择是从redis还是mysql读用户信息
		res1, err := grpc_client.UserClient.GetVActiveFollowerUserinfo(context.Background(), &request1.DouyinGetVActiveFollowFollowerUserinfoRequest{
			LoginUserId:      req.LoginUserId,
			UserId:           req.UserId,
			IsV:              res.IsInfluencer,
			IsActive:         res.IsActiver,
			IsFollowFollower: 2,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetVActiveFollowerUserinfo, zap.Error(err))
				return nil, err
			}
		}
		return res1.User, nil
	}
	//先从mysql读出关注列表，再根据用户id批量从用户服务读用户信息
	fs := make([]model.Follow, 0)
	if err := mysql.GetUserFollowerList(&fs, req.UserId); err != nil {
		zap.L().Error(errorGetUserFollowList, zap.Error(err))
		return nil, err
	}
	if len(fs) == 0 {
		return nil, nil
	}
	var idList []int64
	for i := 0; i < len(fs); i++ {
		idList = append(idList, fs[i].FollowerID)
	}
	res2, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
		UserId:      idList,
		LoginUserId: req.LoginUserId,
	})
	if err != nil {
		if res2 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res2.Code == 2 {
			zap.L().Error(errorGetUserInfoList, zap.Error(err))
			return nil, err
		}
	}
	return res2.User, nil
}

func GetFriendList(req *request.DouyinRelationFriendListRequest) ([]*response1.FriendUser, error) {
	//判断用户是否为大V或活跃用户
	res, err := grpc_client.UserClient.UserIsInfluencerActiver(context.Background(), &request1.DouyinUserIsInfluencerActiverRequest{
		UserId: req.UserId,
	})
	if err != nil {
		if res == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res.Code == 2 {
			zap.L().Error(errorExeVActiveSet, zap.Error(err))
			return nil, err
		}
	}
	if res.IsInfluencer == true || res.IsActiver == true {
		//根据判断条件选择是从redis还是mysql读用户信息
		res1, err := grpc_client.UserClient.GetVActiveFollowerUserinfo(context.Background(), &request1.DouyinGetVActiveFollowFollowerUserinfoRequest{
			LoginUserId:      req.LoginUserId,
			UserId:           req.UserId,
			IsV:              res.IsInfluencer,
			IsActive:         res.IsActiver,
			IsFollowFollower: 2,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetVActiveFollowerUserinfo, zap.Error(err))
				return nil, err
			}
		}
		//取相应消息
		if len(res1.User) != 0 {
			var friendIdList []int64
			var index []int
			for i := 0; i < len(res1.User); i++ {
				if res1.User[i].IsFollow == true {
					friendIdList = append(friendIdList, res1.User[i].Id)
					index = append(index, i)
				}
			}
			if len(friendIdList) == 0 {
				return nil, nil
			}
			res2, err := grpc_client.MessageClient.GetUserFriendMessage(context.Background(), &request2.DouyinGetUserFriendMessageRequest{
				LoginUserId: req.LoginUserId,
				ToUserId:    friendIdList,
			})
			if err != nil {
				zap.L().Error(errorGetUserFriendMessage, zap.Error(err))
				return nil, err
			}
			var results []*response1.FriendUser
			for i := 0; i < len(friendIdList); i++ {
				var result = response1.FriendUser{
					Id:            res1.User[index[i]].Id,
					Name:          res1.User[index[i]].Name,
					FollowCount:   res1.User[index[i]].FollowCount,
					FollowerCount: res1.User[index[i]].FollowerCount,
					IsFollow:      res1.User[index[i]].IsFollow,
					Message:       res2.Message[i],
					MsgType:       res2.MsgType[i],
				}
				results = append(results, &result)
			}
			return results, nil
		}
		return nil, nil
	}
	//先从mysql读出关注列表，再根据用户id批量从用户服务读用户信息
	fs := make([]model.Follow, 0)
	if err := mysql.GetUserFollowerList(&fs, req.UserId); err != nil {
		zap.L().Error(errorGetUserFollowList, zap.Error(err))
		return nil, err
	}
	fs1 := make([]model.Follow, 0)
	if err := mysql.GetUserFollowList(&fs1, req.UserId); err != nil {
		zap.L().Error(errorGetUserFollowList, zap.Error(err))
		return nil, err
	}
	var friendIdList1 []int64
	for i := 0; i < len(fs); i++ {
		friendIdList1 = append(friendIdList1, fs[i].FollowerID)
	}
	var friendIdList2 []int64
	for i := 0; i < len(fs1); i++ {
		friendIdList2 = append(friendIdList2, fs1[i].FollowID)
	}
	m := make(map[int64]bool, 0)
	friendIdList := make([]int64, 0)
	index := make([]int, 0)
	for _, v1 := range friendIdList1 {
		if _, ok := m[v1]; !ok {
			m[v1] = true
		}
	}
	for i, v1 := range friendIdList2 {
		if _, ok := m[v1]; ok {
			friendIdList = append(friendIdList, v1)
			index = append(index, i)
		}
	}
	if len(friendIdList) == 0 {
		return nil, nil
	}
	res2, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
		UserId:      friendIdList,
		LoginUserId: req.LoginUserId,
	})
	if err != nil {
		if res2 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return nil, err
		}
		if res2.Code == 2 {
			zap.L().Error(errorGetUserInfoList, zap.Error(err))
			return nil, err
		}
	}
	//从消息服务取出相关消息
	if len(res2.User) != 0 {
		var friendIdListF []int64
		for i := 0; i < len(res2.User); i++ {
			friendIdListF = append(friendIdListF, res2.User[i].Id)
		}
		res3, err := grpc_client.MessageClient.GetUserFriendMessage(context.Background(), &request2.DouyinGetUserFriendMessageRequest{
			LoginUserId: req.LoginUserId,
			ToUserId:    friendIdListF,
		})
		if err != nil {
			zap.L().Error(errorGetUserFriendMessage, zap.Error(err))
			return nil, err
		}
		var results []*response1.FriendUser
		for i := 0; i < len(res2.User); i++ {
			var result = response1.FriendUser{
				Id:            res2.User[i].Id,
				Name:          res2.User[i].Name,
				FollowCount:   res2.User[i].FollowCount,
				FollowerCount: res2.User[i].FollowerCount,
				IsFollow:      res2.User[i].IsFollow,
				Message:       res3.Message[i],
				MsgType:       res3.MsgType[i],
			}
			results = append(results, &result)
		}
		return results, nil
	}
	return nil, nil
}
