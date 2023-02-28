
# 极简版抖音架构设计
## 场景分析
## 整体架构设计
## 技术选型
文件上传oss
## 数据库、中间件设计
## 各个接口 http rpc
### Mysql设计
### Redis设计
#### 1.用户service
##### 1.记录用户登录次数
**数据结构：string**
**key命名规范：user_login_count:user_id**
**过期时间：一周**

##### 2.活跃者列表 定时任务
**数据结构：set**
**key命名规范：active_users_id**
##### 3.大V列表
**数据结构：set**
**key命名规范：usersV_id**

##### 4.消息活跃用户列表 定时任务
**数据结构：string**
**key:user_id scorce:0.5×get+0.5×post**
**key命名规范：commitV_info:video_id**

##### 5.记录大V的用户信息
**数据结构：string**
**key命名规范：userV_info:user_id**
##### 6.记录活跃用户的用户信息 定时任务
**数据结构：string**
**key命名规范：user_active_info:user_id**

##### 7.记录大V关注的用户信息
**数据结构：hash**
**key命名规范：userV_follow_info:user_id**
##### 8.记录活跃用户关注的用户信息
**数据结构：hash**
**key命名规范：userV_follow_info:user_id**

##### 9.记录大V的粉丝信息
**数据结构：hash**
**key命名规范：userV_follower_info:user_id**
##### 10.记录活跃用户的粉丝信息
**数据结构：hash**
**key命名规范：userV_follower_info:user_id**

##### 11.最近登录次数多的列表

#### 2.视频service

##### 1.记录最新发布的300个视频
**数据结构：list**
**key命名规范：video_list**

##### 2.记录用户发布视频数
**数据结构：string**
**key命名规范：user_video_count:user_id**

##### 3.记录大V发布的视频信息
**数据结构：hash**
**key命名规范：videoV_info:user_id**
##### 4.记录活跃者用户发布的视频信息
**数据结构：hash**
**key命名规范：active_users:user_id**

##### 5.记录大V喜欢的视频信息
**数据结构：hash**
**key命名规范：videoV_info_favorite:user_id**
##### 6.记录活跃者喜欢的视频信息
**数据结构：hash**
**key命名规范：active_users_favorite:user_id**

##### 7.记录所有大V发布的视频id
**数据结构：set**
**key命名规范：videoV_info_favorite:user_id**
##### 8.记录所有活跃者发布的视频id
**数据结构：set**
**key命名规范：active_users_favorite:user_id**

#### 3.评论service
##### 1.记录视频评论数
**数据结构：string**
**key命名规范：video_commit_count:video_id**
##### 2.记录大V发布的视频的评论信息
**数据结构：hash**
**key命名规范：commitV_info:video_id**

#### 4.点赞service
##### 1.记录视频点赞数
**数据结构：string**
**key命名规范：video_favorite_count:video_id**
##### 2.记录用户点赞视频数
**数据结构：string**
**key命名规范：user_favorite_video:user_id**

#### 5.关注service
##### 1.记录用户关注数
**数据结构：string**
**key命名规范：user_follow_count:user_id**
##### 2.记录用户粉丝数
**数据结构：string**
**key命名规范：user_follower_count:user_id**

#### 6.消息service
##### 1.统计用户消息次数 只缓存前百分之一 所以用zscore
**数据结构：zset**
**key:user_id scorce:0.5×get+0.5×post**
**key命名规范：commitV_info:video_id**
##### 2.消息活跃用户的消息
**数据结构：string**
**key:user_id value:用户消息**

##### 3.用户消息get次数
**数据结构：string**
**key:user_id value:用户消息**
##### 4.用户消息post次数
**数据结构：string**
**key:user_id value:用户消息**

### RocketMQ设计
#### 1.用户点赞
#### 2.用户关注
#### 3.用户发消息

## 接口设计
### 用户相关接口
#### 1.用户注册接口
**服务链路**
```mermaid
graph TD;
    Web层-->用户service层;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->生成用户ID;
    生成用户ID-->密码加密;
    参数校验--失败-->返回响应信息;
    密码加密-->生成token;
    生成token-->写入Mysql数据库;
    写入Mysql数据库-->返回响应信息
   
```
#### 2.用户登录接口
```mermaid
graph TD;
    Web层-->用户service层;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->密码加密;
    参数校验--失败-->返回响应信息;
    密码加密-->根据用户名查Mysql数据库;
    根据用户名查Mysql数据库--没查出结果-->返回响应信息;
    根据用户名查Mysql数据库--查出结果-->密码比对
    密码比对--失败-->返回响应信息
    密码比对--成功-->生成token
    生成token-->将用户登录次数记录在redis中并设置过期时间为一周string
    将用户登录次数记录在redis中并设置过期时间为一周string-->返回web
    返回web-->返回响应信息;
```
#### 3.用户信息接口
**服务链路**
```mermaid
graph TD;
    Web层-->用户service;
    用户service-->关注service;
```
**is_follow:
0-登录用户没有关注请求用户
1-登录用户关注请求用户
2-登录用户不知道是否关注了请求用户，需要后续到关注表中查**

**is_V:
true-是
false-不是**
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token;
    验证token--成功-->在redis中读取用户粉丝数用户关注数
    验证token--失败-->返回响应信息;
   在redis中读取用户粉丝数用户关注数-->判断请求用户是否为大V或活跃用户
   判断请求用户是否为大V或活跃用户--读不到 说明请求用户不是大V或活跃用户-->不是大V活跃-判断登录用户与请求用户是否为同一个人
   不是大V活跃-判断登录用户与请求用户是否为同一个人--是-->is_follow为true从用户表读取信息
   is_follow为true从用户表读取信息-->返回响应信息;
   不是大V活跃-判断登录用户与请求用户是否为同一个人--否-->读Mysql关注表和用户信息表
   读Mysql关注表和用户信息表-->返回响应信息;
   判断请求用户是否为大V或活跃用户--能读到说明请求用户是大V或活跃-->是大V活跃-判断登录用户与请求用户是否为同一个人
   是大V活跃-判断登录用户与请求用户是否为同一个人--是-->is_follow设为true从redis读取用户信息
  is_follow设为true从redis读取用户信息-->返回响应信息;
   是大V活跃-判断登录用户与请求用户是否为同一个人--否-->读大V活跃用户粉丝信息表和redis用户信息
   读大V活跃用户粉丝信息表和redis用户信息-->返回响应信息;
  
```
### 视频相关接口
#### 1.视频流接口
**服务链路**
```mermaid
graph TD;
    Web层-->视频service;
    视频service-->点赞service
    点赞service-->评论service
    评论service-->用户service
    用户service-->关注service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->进入视频service
    进入视频service-->读取redis最新发布的300个视频list
    读取redis最新发布的300个视频list--读不到-->去Mysql视频表读
    去Mysql视频表读-->点赞service
    点赞service-->redis中读视频对应的点赞数
    redis中读视频对应的点赞数-->Mysql点赞表中查看用户是否对视频点赞
    Mysql点赞表中查看用户是否对视频点赞-->评论service
    评论service-->redis中读视频对应的评论数
    redis中读视频对应的评论数-->走获取用户信息的那一套流程
    读取redis最新发布的300个视频list--读到了-->点赞service
    走获取用户信息的那一套流程-->返回响应信息
```
#### 2.发布列表 登录用户的视频发布列表
**服务链路**
```mermaid
graph TD;
    Web层-->视频service;
    视频service-->用户1service
    用户1service-->点赞service
    点赞service-->评论service
    评论service-->用户service
    用户service-->关注service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token--成功-->进入视频service
    进入视频service-->判断请求用户是否是大V
    判断请求用户是否是大V--是-->从redis中读取大V视频信息
    从redis中读取大V视频信息-->点赞service
    判断请求用户是否是大V--否-->用户service
    用户service-->redis中查看用户是否属于活跃用户
    redis中查看用户是否属于活跃用户--是-->从redis中读取活跃用户视频信息
    从redis中读取活跃用户视频信息-->点赞service
    redis中查看用户是否属于活跃用户--否-->从Mysql视频表中读取请求用户视频信息
    从Mysql视频表中读取请求用户视频信息-->点赞service
    点赞service-->redis中读视频对应的点赞数
    redis中读视频对应的点赞数-->Mysql点赞表中查看用户是否对视频点赞
    Mysql点赞表中查看用户是否对视频点赞-->评论service
    评论service-->redis中读视频对应的评论数
    redis中读视频对应的评论数-->走获取用户信息的那一套流程
    走获取用户信息的那一套流程-->返回响应信息
```
#### 3.视频投稿
**服务链路**
```mermaid
graph TD;
    Web层-->视频service;
    视频service-->关注service
    视频service-->用户service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token--成功-->进入视频service
    进入视频service-->生成视频ID
    生成视频ID-->根据视频信息生成视频url和封面url
    根据视频信息生成视频url和封面url-->进入关注service
    进入关注service-->读redis中用户的粉丝数判断是否为大V
    读redis中用户的粉丝数判断是否为大V--是-->返回视频service
    读redis中用户的粉丝数判断是否为大V--否-->进入用户service
    进入用户service-->判断投稿用户是否为活跃用户
    判断投稿用户是否为活跃用户-->返回视频service
    返回视频service--是大V-->将视频信息写入redis中的大V视频信息
    将视频信息写入redis中的大V视频信息-->将视频ID写入所有大V发布的视频ID列表
    将视频ID写入所有大V发布的视频ID列表-->将视频信息写入redis中的最新300个视频list
    返回视频service--不是大V是活跃用户-->将视频信息写入redis中的活跃者视频信息
    将视频信息写入redis中的活跃者视频信息-->将视频ID写入所有活跃用户发布的视频ID列表
    将视频ID写入所有活跃用户发布的视频ID列表-->将视频信息写入redis中的最新300个视频list
    返回视频service--不是大V不是活跃用户-->将视频信息写入redis中的最新300个视频list
    将视频信息写入redis中的最新300个视频list-->将视频信息写入Mysql视频表
    将视频信息写入Mysql视频表-->redis中用户视频数加一
    redis中用户视频数加一-->返回响应信息
```
### 点赞相关接口
#### 1.赞操作
**服务链路**
```mermaid
graph TD;
    Web层-->视频service;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断是点赞还是取消点赞
    判断是点赞还是取消点赞--点赞-->进入视频service点赞接口
    进入视频service点赞接口-->判断视频所属用户是否为大V或活跃用户
    判断视频所属用户是否为大V或活跃用户--是-->rocketMQ根视频信息存到redis只存userid
    判断视频所属用户是否为大V或活跃用户--否-->redis用户点赞视频数加一
    redis用户点赞视频数加一-->redis视频点赞数加一
    redis视频点赞数加一-->点赞关系写入Mysql点赞表
    点赞关系写入Mysql点赞表-->返回响应消息
    判断是点赞还是取消点赞--取消点赞-->进入视频service取消点赞接口
    进入视频service取消点赞接口-->判断视频所属用户是否为大V或活跃用户1
    判断视频所属用户是否为大V或活跃用户1--是-->rocketMQ视频信息存到redis只存userid1
    判断视频所属用户是否为大V或活跃用户1--否-->redis用户点赞视频数减一
    redis用户点赞视频数减一-->redis视频点赞数减一
    redis视频点赞数减一-->点赞关系从Mysql点赞表删除
    点赞关系从Mysql点赞表删除-->返回响应消息
```
#### 2.喜欢列表 登录用户的所有点赞视频
**服务链路**
```mermaid
graph TD;
    Web层-->视频service;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->进入关注service
   进入关注service-->判断请求用户是否是大V
    判断请求用户是否是大V--是-->从redis中读取大V喜欢视频信息
    从redis中读取大V喜欢视频信息-->点赞service
    判断请求用户是否是大V--否-->redis中查看用户是否属于活跃用户
    redis中查看用户是否属于活跃用户--是-->从redis中读取活跃用户喜欢视频信息
    从redis中读取活跃用户喜欢视频信息-->点赞service
    redis中查看用户是否属于活跃用户--否-->从Mysql视频表中读取请求用户喜欢视频信息
    从Mysql视频表中读取请求用户喜欢视频信息-->点赞service
    点赞service-->redis中读视频对应的点赞数
    redis中读视频对应的点赞数-->Mysql点赞表中查看用户是否对视频点赞
    Mysql点赞表中查看用户是否对视频点赞-->评论service
    评论service-->redis中读视频对应的评论数
    redis中读视频对应的评论数-->走获取用户信息的那一套流程
    走获取用户信息的那一套流程-->返回响应信息
```
### 评论相关接口
#### 1.评论操作
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断是添加评论还是删除评论
    判断是添加评论还是删除评论--添加评论-->进入评论service评论接口
    进入评论service评论接口-->判断视频所属用户是否为大V
    判断视频所属用户是否为大V--是-->评论信息存到redis
    评论信息存到redis-->redis视频评论数加一
    判断视频所属用户是否为大V--否-->redis视频评论数加一
    redis视频评论数加一-->评论写入Mysql评论表
    评论写入Mysql评论表-->返回响应消息
    判断是添加评论还是删除评论--删除评论-->进入评论service删除评论接口
    进入评论service删除评论接口-->判断视频所属用户是否为大V1
    判断视频所属用户是否为大V1--是-->删除redis中的评论信息
    删除redis中的评论信息-->redis视频评论数减一
    判断视频所属用户是否为大V1--否-->redis视频评论数减一
    redis视频评论数减一-->评论从Mysql评论表删除
    评论从Mysql评论表删除-->返回响应消息
```
#### 2.视频评论列表
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断视频所属用户是否是大V
    判断视频所属用户是否是大V--是-->读redis中大V视频评论
    读redis中大V视频评论-->获取用户信息那套
    判断视频所属用户是否是大V--否-->Mysql中读大V视频评论
    Mysql中读大V视频评论-->获取用户信息那套
    获取用户信息那套-->返回响应消息
```
### 关注相关接口
#### 1.关系操作
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断是关注还是取消关注
    判断是关注还是取消关注--关注-->判断主动关注的用户是否为大V或活跃用户
    判断主动关注的用户是否为大V或活跃用户--是-->从Mysql用户表读取被关注者的用户信息到大V活跃关注者信息redis
    从Mysql用户表读取被关注者的用户信息到大V活跃关注者信息redis-->进入关注service关注接口
    判断主动关注的用户是否为大V或活跃用户--否-->进入关注service关注接口
    进入关注service关注接口-->判断关注用户是否为大V或活跃用户
     判断关注用户是否为大V或活跃用户--是-->rocketMQredis大V活跃用户粉丝信息
    判断关注用户是否为大V或活跃用户--否-->redis用户关注数加一
    redis用户关注数加一-->redis用户粉丝数加一
    redis用户粉丝数加一-->关注关系写入Mysql关注表
    关注关系写入Mysql关注表-->返回响应消息
    判断是关注还是取消关注--取消关注-->判断取消关注的用户是否为大V或活跃用户
    判断取消关注的用户是否为大V或活跃用户--是-->删除大V活跃关注者信息redis对应列
    删除大V活跃关注者信息redis对应列-->进入关注service取消关注接口
    判断取消关注的用户是否为大V或活跃用户--否-->进入关注service取消关注接口
    进入关注service取消关注接口-->判断关注用户是否为大V活跃用户1
    判断关注用户是否为大V活跃用户1--是-->rocketMQ删redis大V活跃粉丝I信息
    判断关注用户是否为大V活跃用户1--否-->redis用户关注数减一
    redis用户关注数减一-->redis用户粉丝数减一
    redis用户粉丝数减一-->关注关系从Mysql关注表删除
    关注关系从Mysql关注表删除-->返回响应消息
```
#### 2.用户关注列表
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断用户是否为大V或活跃用户
    判断用户是否为大V或活跃用户--是-->从redis中读
    从redis中读-->获取用户其它信息那一套
    获取用户其它信息那一套-->返回响应信息
    判断用户是否为大V或活跃用户--否-->从mysql读
    从mysql读-->获取用户其它信息那一套
```
#### 3.用户粉丝列表
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断用户是否为大V或活跃用户
    判断用户是否为大V或活跃用户--是-->从redis中读
    从redis中读-->获取用户其它信息那一套
    获取用户其它信息那一套-->返回响应信息
    判断用户是否为大V或活跃用户--否-->从mysql读
    从mysql读-->获取用户其它信息那一套
```
#### 4.用户好友列表
**服务链路**
```mermaid
graph TD;
    Web层-->评论service;
    评论service-->视频service
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->判断用户是否为大V或活跃用户
    判断用户是否为大V或活跃用户--是-->从redis中读
    从redis中读-->获取用户其它信息那一套
    获取用户其它信息那一套-->返回响应信息
    判断用户是否为大V或活跃用户--否-->从mysql读
    从mysql读-->获取用户其它信息那一套
```
### 消息相关接口
#### 1.消息操作
**服务链路**
```mermaid
graph TD;
    Web层-->消息service;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->消息service
    消息service-->查看用户是否是消息活跃用户
    查看用户是否是消息活跃用户--是-->存入活跃用户消息redis
    存入活跃用户消息redis-->rocketMQ存入Mysql消息表
    rocketMQ存入Mysql消息表-->返回响应信息
    查看用户是否是消息活跃用户--否-->存入Mysql消息表
    存入Mysql消息表-->返回响应信息
```
#### 2.聊天记录
**服务链路**
```mermaid
graph TD;
    Web层-->消息service;
```
**流程**
```mermaid
graph TD;
    接收参数-->参数校验;
    参数校验--成功-->验证token
    验证token-->消息service
    消息service-->查看用户是否是消息活跃用户
    查看用户是否是消息活跃用户--是-->从redis中读消息
    从redis中读消息-->返回响应信息
    查看用户是否是消息活跃用户--否-->从Mysql中读消息
    从Mysql中读消息-->返回响应信息
```
## 代码示例，环境、中间件等配置
## 性能优化 不足之处
当活跃用户的粉丝数到20时，会自动升级为大V，这时候要删除在活跃用户中存的各种信息，若在判断用户是否为活跃用户、执行后续逻辑之前（如获取活跃用户信息）恰好删除，可能会读不到信息、报错