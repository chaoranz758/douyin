package service

const (
	errorConnectToGRPCServer            = "connect to grpc server failed"
	errorExeVActiveSet                  = "exe if request is V or active from redis set failed"
	errorPushVActiveFollowerUserinfo    = "push v active follower user information failed"
	errorAddUserFollowFollowerCount     = "add user follow follower count failed"
	errorCreateFollow                   = "create follow failed"
	errorDeleteVActiveFollowerUserinfo  = "delete v active follower user information failed"
	errorSubUserFollowFollowerCount     = "sub user follow follower count failed"
	errorDeleteFollow                   = "delete follow failed"
	errorGetVActiveFollowerUserinfo     = "get v active follower user information failed"
	errorGetUserFollowList              = "get user follow list failed"
	errorGetUserInfoList                = "get user information list failed"
	errorJudgeUserIsFollow              = "judge user is follow failed"
	errorJudgeUserIsFollowList          = "judge user is follow list failed"
	errorGetUserFollowFollowerCount     = "get user follow follower count failed"
	errorGetUserFollowFollowerCountList = "get user follow follower count list failed"
	errorAddUserFollowUserCountSet      = "add user follow user count set failed"
	errorAddUserFollowerUserCountSet    = "add user follower user count set failed"
	errorGetUserFriendMessage           = "get user friend message failed"
	errorRegisterWorkflow               = "workflow register failed"
	errorExecuteWorkflow                = "result of workflow.Execute is"
	errorSendMessage1                   = "生产者1消息发送失败"
	errorSendMessage2                   = "生产者2消息发送失败"
	errorSendMessage3                   = "生产者3消息发送失败"
)

const (
	successSendMessage1 = "生产者1消息发送成功"
	successSendMessage2 = "生产者2消息发送成功"
	successSendMessage3 = "生产者3消息发送成功"
)
