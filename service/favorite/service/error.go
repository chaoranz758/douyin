package service

const (
	errorConnectToGRPCServer          = "connect to grpc server failed"
	errorExeVActiveSet                = "exe if request is V or active from redis set failed"
	errorJudgeVideoAuthor             = "judge video author is v or active failed"
	errorPushVActiveFavoriteVideo     = "push v or active favorite video information failed"
	errorAddFavoriteCount             = "add favorite count failed"
	errorFavoriteRelation             = "create favorite relation failed"
	errorDeleteVActiveFavoriteVideo   = "delete v active favorite video information failed"
	errorSubFavoriteCount             = "sub favorite count failed"
	errorDeleteFavoriteRelation       = "delete favorite relation failed"
	errorGetVActiveFavoriteVideo      = "get v or active favorite video information failed"
	errorGetUserFavoriteID            = "get user favorite id failed"
	errorGetVideoListInner            = "get video list inner failed"
	MGetVideoFavoriteCount            = "get video favorite count failed"
	errorGetUserFavoriteBool          = "get user favorite bool failed"
	errorGetCommentCount              = "get comment count failed"
	errorGetUserInfoList              = "get user information list failed"
	errorAddUserFavoriteVideoCountSet = "add user favorite video count set failed"
	errorGetUserFavoritedCount        = "get user favorited count failed"
	errorGetFavoriteCount             = "get user favorite videos count failed"
	errorJsonUnmarshal                = "json unmarshal failed"
	errorSendMessage1                 = "生产者1消息发送失败"
	errorSendMessage2                 = "生产者2消息发送失败"
	errorSendMessage3                 = "生产者3消息发送失败"
	errorSendMessage4                 = "生产者4消息发送失败"
	errorExecuteWorkflow              = "result of workflow.Execute is"
	errorRpcInput                     = "rpc input failed"
)

const (
	successSendMessage1 = "生产者1消息发送成功"
	successSendMessage2 = "生产者2消息发送成功"
	successSendMessage3 = "生产者3消息发送成功"
	successSendMessage4 = "生产者4消息发送成功"
)
