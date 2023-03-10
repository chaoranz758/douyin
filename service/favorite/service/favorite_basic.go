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
	//????????????
	if req.ActionType == 1 {
		//???????????????????????????????????????V???????????????
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
		//??????????????????????????????V???????????????,?????????????????????ID??????
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
			//????????????????????????
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
				//zap.L().Info("????????????????????????????????????")
				return nil
			})
			//??????????????????
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
			//zap.L().Info("????????????????????????????????????????????????")
			//????????????????????????
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
				//zap.L().Info("????????????????????????????????????")
				return nil
			})
			//????????????
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
			//zap.L().Info("????????????????????????????????????????????????")
			//????????????
			_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
				//???????????????????????????V???????????????
				if res.IsInfluencer == true || res.IsActiver == true {
					//???V??????????????????????????????????????????????????????redis ????????????
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
					//??????????????????Mysql?????????
					var f = model.Favorite{
						VideoID: req.VideoId,
						UserID:  req.LoginUserId,
					}
					if err := mysql.CreateFavoriteRelation(&f); err != nil {
						zap.L().Error(errorFavoriteRelation, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					//redis?????????????????????????????????????????????
					_, err := redis.AddFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId)
					if err != nil {
						zap.L().Error(errorAddFavoriteCount, zap.Error(err))
						//???????????????Mysql???????????????
						if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
							zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
							return nil, status.Error(codes.Aborted, err.Error())
						}
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, nil
				}
				//???V?????????????????????????????????????????????
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
				//??????????????????Mysql?????????
				var f = model.Favorite{
					VideoID: req.VideoId,
					UserID:  req.LoginUserId,
				}
				if err := mysql.CreateFavoriteRelation(&f); err != nil {
					zap.L().Error(errorFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//redis?????????????????????????????????????????????
				_, err := redis.AddFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId)
				if err != nil {
					zap.L().Error(errorAddFavoriteCount, zap.Error(err))
					//???????????????Mysql???????????????
					if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
						zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
						return nil, status.Error(codes.Aborted, err.Error())
					}
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//zap.L().Info("????????????????????????")
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
		//zap.L().Info("????????????dtm??????????????????")
		return nil
	}
	//???????????????????????????????????????V???????????????
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
	//??????????????????????????????V???????????????
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
		//????????????????????????
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
			//zap.L().Info("??????????????????????????????")
			return nil
		})
		//??????????????????
		if res.IsInfluencer == true || res.IsActiver == true {
			if res1.IsV == false && res1.IsActive == false {
				//??????V?????????????????????????????????????????????redis??????
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
		//zap.L().Info("????????????????????????")
		//????????????
		_, err = wf.NewBranch().Do(func(bb *dtmcli.BranchBarrier) ([]byte, error) {
			if res.IsInfluencer == true || res.IsActiver == true {
				if res1.IsV == true || res1.IsActive == true {
					//???V??????????????????????????????
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
				//???????????????Mysql???????????????
				if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
					zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
					return nil, status.Error(codes.Aborted, err.Error())
				}
				//redis?????????????????????????????????????????????
				if err := redis.SubFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId); err != nil {
					zap.L().Error(errorSubFavoriteCount, zap.Error(err))
					//??????????????????Mysql?????????
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
			//???????????????Mysql???????????????
			if err := mysql.DeleteFavoriteRelation(req.VideoId, req.LoginUserId); err != nil {
				zap.L().Error(errorDeleteFavoriteRelation, zap.Error(err))
				return nil, status.Error(codes.Aborted, err.Error())
			}
			//redis?????????????????????????????????????????????
			if err := redis.SubFavoriteCount(req.LoginUserId, req.VideoId, res1.AuthorId); err != nil {
				zap.L().Error(errorSubFavoriteCount, zap.Error(err))
				//??????????????????Mysql?????????
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
			//zap.L().Info("????????????????????????")
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
	//zap.L().Info("??????????????????dtm??????????????????")
	return nil
}

func GetFavoriteVideoList(req *request.DouyinFavoriteListRequest) ([]*response.Video, error) {
	//?????????????????????????????????V???????????????
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
	//?????????????????????V?????????????????????redis????????????????????????
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
		//???????????????????????????
		favoriteCountList, err := redis.GetVideoFavoriteCount(idList)
		if err != nil {
			zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
			return nil, err
		}
		//mysql?????????????????????????????????????????????
		userFavoriteBool, err := mysql.GetUserFavoriteBool(req.LoginUserId, idList)
		if err != nil {
			zap.L().Error(errorGetUserFavoriteBool, zap.Error(err))
			return nil, err
		}
		//???????????????????????????
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
		//????????????????????????????????????
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
		//??????????????????mysql??????????????????????????????ID,???????????????????????????????????????
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
		//???????????????????????????
		favoriteCountList, err := redis.GetVideoFavoriteCount(idList)
		if err != nil {
			zap.L().Error(MGetVideoFavoriteCount, zap.Error(err))
			return nil, err
		}
		//mysql?????????????????????????????????????????????
		userFavoriteBool, err := mysql.GetUserFavoriteBool(req.LoginUserId, idList)
		if err != nil {
			zap.L().Error(errorGetUserFavoriteBool, zap.Error(err))
			return nil, err
		}
		//???????????????????????????
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
		//????????????????????????????????????
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
