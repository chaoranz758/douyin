package redis

const (
	errorJsonMarshal                 = "json marshal failed"
	errorAddUserInfoToZSet           = "add user information to zset failed"
	errorGetVVideoInfo               = "get v video information from redis failed"
	errorGetActiveVideoInfo          = "get active video information from redis failed"
	errorComputeZSetCount            = "compute zset count failed"
	errorSetKeyVVideoID              = "set key v video id failed"
	errorSetKeyActiveVideoID         = "set key active video id failed"
	errorGetVVideoInfoByVideoId      = "get v video information by video id failed"
	errorGetActiveVideoInfoByVideoId = "get active video information by video id failed"
	errorJsonUnMarshal               = "json unmarshal failed"
	errorZRevRangeByScore            = "z rev range by score failed"
	errorZCard                       = "zCard failed"
	warnKey                          = "redis监听的key正在被其它线程执行"
	errorWatchKey                    = "redis watch key failed"
	errorConnectRedis                = "connect redis failed"
	errorTypeTurnFailed              = "type turn failed"
	errorPushVBasicInfoInit          = "push v basic information init phase failed"
	errorDeleteActiveVideoInfo       = "delete active video information failed"
	errorPushActiveBasicInfoInit     = "push active basic information init phase failed"
)
