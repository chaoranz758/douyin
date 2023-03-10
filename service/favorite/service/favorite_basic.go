package service

import (
	"context"
	request3 "douyin/proto/comment/request"
	"douyin/proto/favorite/request"
	request1 "douyin/proto/user/request"
	request2 "douyin/proto/video/request"
	"douyin/proto/video/response"
	"douyin/service/favorite/dao/mysql"
	"douyin/service/favorite/dao/redis"
	"douyin/service/favorite/initialize/grpc_client"
	"douyin/service/favorite/initialize/rocketmq"
	"douyin/service/favorite/model"
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
	wfName1            = "workflow-favoriteVideo"
	favoriteVideoCount = 4
)

const (
	topic1 = "favoriteTopic1"
	topic2 = "favoriteTopic2"
	topic3 = "favoriteTopic3"
	topic4 = "favoriteTopic4"
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

type FavoriteVideoDtmMessage struct {
	VideoId           int64 `json:"videoId"`
	IsV               bool  `json:"isV"`
	IsActive          bool  `json:"isActive"`
	LoginUserIsV      bool  `json:"loginUserIsV"`
	LoginUserIsActive bool  `json:"loginUserIsActive"`
	LoginUserId       int64 `json:"loginUserId"`
	AuthorId          int64 `json:"authorId"`
}

type DeleteFavoriteVideoDtmMessage struct {
	VideoId           int64 `json:"videoId"`
	UserId            int64 `json:"userId"`
	LoginUserIsV      bool  `json:"loginUserIsV"`
	LoginUserIsActive bool  `json:"loginUserIsActive"`
	IsV               bool  `json:"isV"`
	IsActive          bool  `json:"isActive"`
	AuthorId          int64 `json:"authorId"`
}

func FavoriteVideoDtm(req *request.DouyinFavoriteActionRequest) error {
	//点赞接口
	if req.ActionType == 1 {
		//判断点赞视频的用户是否为大V或活跃用户
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
		//判断视频作者是否为大V或活跃用户,若并将视频作者ID返回
		res1, err := grpc_client.VideoClient.JudgeVideoAuthor(context.Background(), &request2.DouyinJudgeVideoAuthorRequest{
			VideoId: req.VideoId,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return err
			}
			if res1.Code == 2 {
				zap.L().Error(errorJudgeVideoAuthor, zap.Error(err))
				return err
			}
		}
		wfName := wfName1 + shortuuid.New()
		err = workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
			var req1 FavoriteVideoDtmMessage
			if err := json.Unmarshal(data, &req1); err != nil {
				zap.L().Error(errorJsonUnmarshal, zap.Error(err))
				return status.New(codes.Aborted, err.Error()).Err()
			}
			//视频服务回滚分支
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				if res.IsInfluencer == true || res.IsActiver == true {
					if res1.IsV == false && res1.IsActive == false {
						res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideoRevert(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
							VideoId:           req1.VideoId,
							IsV:               req1.IsV,
							IsActive:          req1.IsActive,
							LoginUserIsV:      req1.LoginUserIsV,
							LoginUserIsActive: req1.LoginUserIsActive,
							LoginUserId:       req1.LoginUserId,
							AuthorId:          req1.AuthorId,
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
					}
				}
				//zap.L().Info("点赞接口调用视频服务回滚")
				return nil
			})
			//视频服务分支
			if res.IsInfluencer == true || res.IsActiver == true {
				if res1.IsV == false && res1.IsActive == false {
					res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideo(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
						VideoId:           req1.VideoId,
						IsV:               req1.IsV,
						IsActive:          req1.IsActive,
						LoginUserIsV:      req1.LoginUserIsV,
						LoginUserIsActive: req1.LoginUserIsActive,
						LoginUserId:       req1.LoginUserId,
						AuthorId:          req1.AuthorId,
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
				}
			}
			//zap.L().Info("点赞接口调用视频服务正向事务成功")
			//用户服务回滚分支
			wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
				res5, err := grpc_client.UserClientDtm.AddUserFavoriteVideoCountSetRevert(wf.Context, &request1.DouyinUserVideoCountSetRequest{
					UserId: req1.LoginUserId,
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
				//zap.L().Info("点赞接口调用用户服务回滚")
				return nil
			})
			//用户服务
			if res.IsInfluencer == false && res.IsActiver == false && res1.IsV == false && res1.IsActive == false {
				count, err := redis.GetFavoriteCount(req.LoginUserId)
				if err != nil {
					zap.L().Error(errorGetFavoriteCount, zap.Error(err))
					return err
				}
				if count == favoriteVideoCount {
					res5, err := grpc_client.UserClientDtm.AddUserFavoriteVideoCountSet(wf.Context, &request1.DouyinUserVideoCountSetRequest{
						UserId: req1.LoginUserId,
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
			}
			//zap.L().Info("点赞接口调用用户服务正向事务成功")
			//本地事务
			_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
				//点赞视频的用户是大V或活跃用户
				if res.IsInfluencer == true || res.IsActiver == true {
					//大V或活跃用户走消息队列：将视频信息存入redis 后续流程
					if res1.IsV == true || res1.IsActive == true {
						var producer1Message = Producer1Message{
							LoginUserID:       req.LoginUserId,
							VideoID:           req.VideoId,
							IsV:               res1.IsV,
							IsActive:          res1.IsActive,
							LoginUserIsV:      res.IsInfluencer,
							LoginUserIsActive: res.IsActiver,
							AuthorId:          res1.AuthorId,
						}
						data, _ := json.Marshal(producer1Message)
						msg := &primitive.Message{
							Topic: topic1,
							Body:  data,
						}
						_, err = rocketmq.Producer1.SendSync(context.Background(), msg)
						if err != nil {
							zap.L().Error(errorSendMessage1, zap.Error(err))
							return nil, status.Error(codes.Aborted, err.Error())
						}
						zap.L().Info(successSendMessage1)
						return nil, nil
					}
					//点赞关系写入Mysql点赞表
					var f = model.Favorite{
						VideoID: req.VideoId,
						UserID:  req.LoginUserId,
					}
					if err := mysql.CreateFavoriteRelation(&f); err != nil {
						zap.L().Error(errorFavoriteRelation, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//redis用户点赞视频数和视频点赞数加一
					_, err := redis.AddFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId)
					if err != nil {
						zap.L().Error(errorAddFavoriteCount, zap.Error(err))
						//点赞关系从Mysql点赞表删除
						if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
							zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
							return nil, status.Error(codes.Aborted, err.Error())
						}
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, nil
				}
				//大V或活跃用户走消息队列：后续流程
				if res1.IsV == true || res1.IsActive == true {
					var producer2Message = Producer2Message{
						LoginUserID: req.LoginUserId,
						VideoID:     req.VideoId,
						AuthorId:    res1.AuthorId,
					}
					data, _ := json.Marshal(producer2Message)
					msg := &primitive.Message{
						Topic: topic2,
						Body:  data,
					}
					_, err = rocketmq.Producer2.SendSync(context.Background(), msg)
					if err != nil {
						zap.L().Error(errorSendMessage2, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info(successSendMessage2)
					return nil, nil
				}
				//点赞关系写入Mysql点赞表
				var f = model.Favorite{
					VideoID: req.VideoId,
					UserID:  req.LoginUserId,
				}
				if err := mysql.CreateFavoriteRelation(&f); err != nil {
					zap.L().Error(errorFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//redis用户点赞视频数和视频点赞数加一
				_, err := redis.AddFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId)
				if err != nil {
					zap.L().Error(errorAddFavoriteCount, zap.Error(err))
					//点赞关系从Mysql点赞表删除
					if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
						zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//zap.L().Info("本地事务执行成功")
				return nil, nil
			})
			if err != nil {
				return err
			}
			return nil
		})
		var req1 = FavoriteVideoDtmMessage{
			VideoId:           req.VideoId,
			IsV:               false,
			IsActive:          false,
			LoginUserIsActive: res.IsActiver,
			LoginUserIsV:      res.IsInfluencer,
			LoginUserId:       req.LoginUserId,
			AuthorId:          res1.AuthorId,
		}
		data, _ := json.Marshal(&req1)
		if err = workflow.Execute(wfName, shortuuid.New(), data); err != nil {
			zap.L().Error(errorExecuteWorkflow, zap.Error(err))
			return err
		}
		//zap.L().Info("视频点赞dtm事务执行完成")
		return nil
	}
	//判断点赞视频的用户是否为大V或活跃用户
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
	//判断视频作者是否为大V或活跃用户
	res1, err := grpc_client.VideoClient.JudgeVideoAuthor(context.Background(), &request2.DouyinJudgeVideoAuthorRequest{
		VideoId: req.VideoId,
	})
	if err != nil {
		if res1 == nil {
			zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
			return err
		}
		if res1.Code == 2 {
			zap.L().Error(errorJudgeVideoAuthor, zap.Error(err))
			return err
		}
	}
	wfName := wfName1 + shortuuid.New()
	err = workflow.Register(wfName, func(wf *workflow.Workflow, data []byte) error {
		var req1 DeleteFavoriteVideoDtmMessage
		if err := json.Unmarshal(data, &req1); err != nil {
			zap.L().Error(errorJsonUnmarshal, zap.Error(err))
			return status.New(codes.Aborted, err.Error()).Err()
		}
		//调用视频服务回滚
		wf.NewBranch().OnRollback(func(bb *dtmcli.BranchBarrier) error {
			if res.IsInfluencer == true || res.IsActiver == true {
				if res1.IsV == false && res1.IsActive == false {
					res2, err := grpc_client.VideoClientDtm.PushVActiveFavoriteVideo(wf.Context, &request2.DouyinPushVActiveFavoriteVideoRequest{
						VideoId:           req1.VideoId,
						IsV:               req1.IsV,
						IsActive:          req1.IsActive,
						LoginUserIsV:      req1.LoginUserIsV,
						LoginUserIsActive: req1.LoginUserIsActive,
						LoginUserId:       req1.UserId,
						AuthorId:          req1.AuthorId,
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
				}
			}
			//zap.L().Info("调用视频服务回滚完成")
			return nil
		})
		//调用视频服务
		if res.IsInfluencer == true || res.IsActiver == true {
			if res1.IsV == false && res1.IsActive == false {
				//将大V或活跃用户喜欢的基本视频信息从redis删除
				res2, err := grpc_client.VideoClientDtm.DeleteVActiveFavoriteVideo(wf.Context, &request2.DouyinDeleteVActiveFavoriteVideoRequest{
					UserId:            req.LoginUserId,
					VideoId:           req.VideoId,
					LoginUserIsV:      res.IsInfluencer,
					LoginUserIsActive: res.IsActiver,
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
			}
		}
		//zap.L().Info("调用视频服务完成")
		//本地事务
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			if res.IsInfluencer == true || res.IsActiver == true {
				if res1.IsV == true || res1.IsActive == true {
					//大V或活跃用户走消息队列
					var producer3Message = Producer3Message{
						LoginUserID:       req.LoginUserId,
						VideoID:           req.VideoId,
						LoginUserIsActive: res.IsActiver,
						LoginUserIsV:      res.IsInfluencer,
						AuthorId:          res1.AuthorId,
						IsV:               res1.IsV,
						IsActive:          res1.IsActive,
					}
					data, _ := json.Marshal(producer3Message)
					msg := &primitive.Message{
						Topic: topic3,
						Body:  data,
					}
					_, err = rocketmq.Producer3.SendSync(context.Background(), msg)
					if err != nil {
						zap.L().Error(errorSendMessage3, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					zap.L().Info(successSendMessage3)
					return nil, nil
				}
				//点赞关系从Mysql点赞表删除
				if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
					zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//redis用户点赞视频数和视频点赞数减一
				if err := redis.SubFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId); err != nil {
					zap.L().Error(errorSubFavoriteCount, zap.Error(err))
					//点赞关系写入Mysql点赞表
					var f = model.Favorite{
						VideoID: req.VideoId,
						UserID:  req.LoginUserId,
					}
					if err := mysql.CreateFavoriteRelation(&f); err != nil {
						zap.L().Error(errorFavoriteRelation, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, status.Error(codes.Aborted, err.Error())
				}
				return nil, nil
			}
			if res1.IsV == true || res1.IsActive == true {
				var producer4Message = Producer4Message{
					LoginUserID: req.LoginUserId,
					VideoID:     req.VideoId,
					AuthorId:    res1.AuthorId,
				}
				data, _ := json.Marshal(producer4Message)
				msg := &primitive.Message{
					Topic: topic4,
					Body:  data,
				}
				_, err = rocketmq.Producer4.SendSync(context.Background(), msg)
				if err != nil {
					zap.L().Error(errorSendMessage4, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				zap.L().Info(successSendMessage4)
				return nil, nil
			}
			//点赞关系从Mysql点赞表删除
			if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
				zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis用户点赞视频数和视频点赞数减一
			if err := redis.SubFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId); err != nil {
				zap.L().Error(errorSubFavoriteCount, zap.Error(err))
				//点赞关系写入Mysql点赞表
				var f = model.Favorite{
					VideoID: req.VideoId,
					UserID:  req.LoginUserId,
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
	var req1 = DeleteFavoriteVideoDtmMessage{
		VideoId:           req.VideoId,
		IsV:               res1.IsV,
		IsActive:          res1.IsActive,
		LoginUserIsActive: res.IsActiver,
		LoginUserIsV:      res.IsInfluencer,
		UserId:            req.LoginUserId,
		AuthorId:          res1.AuthorId,
	}
	data, _ := json.Marshal(&req1)
	if err = workflow.Execute(wfName, shortuuid.New(), data); err != nil {
		zap.L().Error(errorExecuteWorkflow, zap.Error(err))
		return err
	}
	//zap.L().Info("删除视频点赞dtm事务执行完成")
	return nil
}

func GetFavoriteVideoList(req *request.DouyinFavoriteListRequest) ([]*response.Video, error) {
	//判断请求的用户是否为大V或活跃用户
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
	//若请求用户是大V或活跃用户，从redis中取视频基本信息
	if res.IsInfluencer == true || res.IsActiver == true {
		res1, err := grpc_client.VideoClient.GetVActiveFavoriteVideo(context.Background(), &request2.DouyinGetVActiveFavoriteVideoRequest{
			UserId:   req.UserId,
			IsV:      res.IsInfluencer,
			IsActive: res.IsActiver,
		})
		if err != nil {
			if res1 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res1.Code == 2 {
				zap.L().Error(errorGetVActiveFavoriteVideo, zap.Error(err))
				return nil, err
			}
		}
		var idList []int64
		var userIdList []int64
		for i := 0; i < len(res1.VideoList); i++ {
			idList = append(idList, res1.VideoList[i].Id)
			userIdList = append(userIdList, res1.VideoList[i].UserId)
		}
		//读视频对应的点赞数
		favoriteCountList, err := redis.GetVideoFavoriteCount(idList)
		if err != nil {
			zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
			return nil, err
		}
		//mysql点赞表中查看用户是否对视频点赞
		userFavoriteBool, err := mysql.GetUserFavoriteBool(req.LoginUserId, idList)
		if err != nil {
			zap.L().Error(errorGetUserFavoriteBool, zap.Error(err))
			return nil, err
		}
		//读视频对应的评论数
		res3, err := grpc_client.CommentClient.GetCommentCount(context.Background(), &request3.DouyinCommentCountRequest{
			VideoId: idList,
		})
		if err != nil {
			if res3 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res3.Code == 2 {
				zap.L().Error(errorGetCommentCount, zap.Error(err))
				return nil, err
			}
		}
		//按照用户列表获取用户信息
		res4, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
			UserId:      userIdList,
			LoginUserId: req.LoginUserId,
		})
		if err != nil {
			if res4 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res4.Code == 2 {
				zap.L().Error(errorGetUserInfoList, zap.Error(err))
				return nil, err
			}
		}
		var result []*response.Video
		for i := 0; i < len(idList); i++ {
			var r = response.Video{
				Id:            res1.VideoList[i].Id,
				PlayUrl:       res1.VideoList[i].PlayUrl,
				CoverUrl:      res1.VideoList[i].CoverUrl,
				FavoriteCount: favoriteCountList[i],
				CommentCount:  res3.CommentCount[i],
				IsFavorite:    userFavoriteBool[i],
				Title:         res1.VideoList[i].Title,
				Author:        res4.User[i],
			}
			result = append(result, &r)
		}
		return result, nil
	} else {
		//否则先查本地mysql表用户点赞了哪些视频ID,再查用户点赞的视频基本信息
		f := make([]model.Favorite, 0)
		if err := mysql.GetUserFavoriteID(req.UserId, &f); err != nil {
			zap.L().Error(errorGetUserFavoriteID, zap.Error(err))
			return nil, err
		}
		if len(f) == 0 {
			return nil, nil
		}
		var idList []int64
		for i := 0; i < len(f); i++ {
			idList = append(idList, f[i].VideoID)
		}
		res2, err := grpc_client.VideoClient.GetVideoListInner(context.Background(), &request2.DouyinGetVideoListRequest{
			VideoId: idList,
		})
		if err != nil {
			if res2 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res2.Code == 2 {
				zap.L().Error(errorGetVideoListInner, zap.Error(err))
				return nil, err
			}
		}
		var userIdList []int64
		for i := 0; i < len(f); i++ {
			userIdList = append(userIdList, res2.VideoList[i].UserId)
		}
		//读视频对应的点赞数
		favoriteCountList, err := redis.GetVideoFavoriteCount(idList)
		if err != nil {
			zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
			return nil, err
		}
		//mysql点赞表中查看用户是否对视频点赞
		userFavoriteBool, err := mysql.GetUserFavoriteBool(req.LoginUserId, idList)
		if err != nil {
			zap.L().Error(errorGetUserFavoriteBool, zap.Error(err))
			return nil, err
		}
		//读视频对应的评论数
		res3, err := grpc_client.CommentClient.GetCommentCount(context.Background(), &request3.DouyinCommentCountRequest{
			VideoId: idList,
		})
		if err != nil {
			if res3 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res3.Code == 2 {
				zap.L().Error(errorGetCommentCount, zap.Error(err))
				return nil, err
			}
		}
		//按照用户列表获取用户信息
		res4, err := grpc_client.UserClient.GetUserInfoList(context.Background(), &request1.DouyinUserListRequest{
			UserId:      userIdList,
			LoginUserId: req.LoginUserId,
		})
		if err != nil {
			if res4 == nil {
				zap.L().Error(errorConnectToGRPCServer, zap.Error(err))
				return nil, err
			}
			if res4.Code == 2 {
				zap.L().Error(errorGetUserInfoList, zap.Error(err))
				return nil, err
			}
		}
		var result []*response.Video
		for i := 0; i < len(idList); i++ {
			var r = response.Video{
				Id:            res2.VideoList[i].Id,
				PlayUrl:       res2.VideoList[i].PlayUrl,
				CoverUrl:      res2.VideoList[i].CoverUrl,
				FavoriteCount: favoriteCountList[i],
				CommentCount:  res3.CommentCount[i],
				IsFavorite:    userFavoriteBool[i],
				Title:         res2.VideoList[i].Title,
				Author:        res4.User[i],
			}
			result = append(result, &r)
		}
		return result, nil
	}
}
