# 极简版抖音
## 一句话介绍
实现了基于微服务的极简版抖音的基础接口、互动接口、点赞接口，使用mysql存储数据，基于redis、rocketmq等组件对系统做了性能优化，增加了部分安全功能，并对系统做了功能测试和性能测试。
## 项目启动
### 启动前环境配置
#### 1.安装consul
本项目默认地址为"192.168.182.137:8500" 如需修改请修改配置文件
#### 2.安装nacos
本项目默认地址为"192.168.182.137:8848" 如需修改请修改配置文件
#### 3.安装rocketmq
本项目默认地址为"192.168.182.137:9876" 如需修改请修改配置文件
#### 4.安装dtm
本项目默认地址为"192.168.182.137:36790" 如需修改请修改配置文件
#### 5.安装ffmpeg相关环境
#### 6.配置阿里云oss
参见web层的配置文件
#### 7.配置mysql相关环境
#### 8.配置redis相关环境
### 默认启动端口
web 8080 service 9080-9085 具体参见配置文件
## 项目目录树形图
```
├─.idea
├─proto - service层protobuf文件
│  ├─comment - 评论service层protobuf文件
│  │  ├─api - 评论service层rpc接口
│  │  ├─request - 评论service层请求message
│  │  └─response - 评论service层响应message
│  ├─favorite - 点赞service层protobuf文件
│  │  ├─api - 点赞service层rpc接口
│  │  ├─request - 点赞service层请求message
│  │  └─response - 点赞service层响应message
│  ├─follow - 关系service层protobuf文件
│  │  ├─api - 关系service层rpc接口
│  │  ├─request - 关系service层请求message
│  │  └─response - 关系service层响应message
│  ├─message - 消息service层protobuf文件
│  │  ├─api - 消息service层rpc接口
│  │  ├─request - 消息service层请求message
│  │  └─response - 消息service层响应message
│  ├─user - 用户service层protobuf文件
│  │  ├─api - 用户service层rpc接口
│  │  ├─request - 用户service层请求message
│  │  └─response - 用户service层响应message
│  └─video - 视频service层protobuf文件
│      ├─api - 视频service层rpc接口
│      ├─request - 视频service层请求message
│      └─response - 视频service层响应message
├─service - 项目service层
│  ├─comment - 评论service层
│  │  ├─config - 配置文件
│  │  │  └─nacos - 配置中心nacos缓存文件
│  │  │      └─cache
│  │  │          └─config
│  │  ├─dao - 数据读写 mysql redis
│  │  │  ├─mysql
│  │  │  └─redis
│  │  ├─handler - handler层
│  │  ├─initialize - 项目初始化
│  │  │  ├─config - 配置初始化
│  │  │  ├─consul - 注册中心初始化
│  │  │  ├─grpc_client - rpc客户端初始化
│  │  │  ├─log - 日志初始化
│  │  │  ├─server - 开启grpc服务初始化
│  │  │  └─snowflake - 雪花算法初始化
│  │  ├─model - mysql表
│  │  ├─service - 业务逻辑层
│  │  └─util - 工具包
│  ├─favorite - 点赞service层
│  │  ├─config - 配置文件
│  │  │  └─nacos - 配置中心nacos缓存文件
│  │  │      └─cache
│  │  │          └─config
│  │  ├─dao - 数据读写 mysql redis
│  │  │  ├─mysql
│  │  │  └─redis
│  │  ├─handler - handler层
│  │  ├─initialize - 项目初始化
│  │  │  ├─config - 配置初始化
│  │  │  ├─consul - 注册中心初始化
│  │  │  ├─grpc_client - rpc客户端初始化
│  │  │  ├─log - 日志初始化
│  │  │  ├─rocketmq - 消息队列初始化
│  │  │  └─server - 开启grpc服务初始化
│  │  ├─middleware - 中间件层
│  │  │  └─rocketmq - 消息队列
│  │  ├─model - mysql表
│  │  └─service - 业务逻辑层
│  ├─follow - 关系service层
│  │  ├─config - 配置文件
│  │  │  └─nacos - 配置中心nacos缓存文件
│  │  │      └─cache
│  │  │          └─config
│  │  ├─dao - 数据读写 mysql redis
│  │  │  ├─mysql
│  │  │  └─redis
│  │  ├─handler - handler层
│  │  ├─initialize - 项目初始化
│  │  │  ├─config - 配置初始化
│  │  │  ├─consul - 注册中心初始化
│  │  │  ├─grpc_client - rpc客户端初始化
│  │  │  ├─log - 日志初始化
│  │  │  ├─rocketmq - 消息队列初始化
│  │  │  │  ├─costumer
│  │  │  │  └─producer
│  │  │  └─server - 开启grpc服务初始化
│  │  ├─middleware - 中间件层
│  │  │  └─rocketmq - 消息队列
│  │  ├─model - mysql表
│  │  └─service - 业务逻辑层
│  ├─message - 消息service层
│  │  ├─config - 配置文件
│  │  │  └─nacos - 配置中心nacos缓存文件
│  │  │      └─cache
│  │  │          └─config
│  │  ├─dao - 数据读写 mysql redis
│  │  │  ├─mysql
│  │  │  └─redis
│  │  ├─handler - handler层
│  │  ├─initialize - 项目初始化
│  │  │  ├─config - 配置初始化
│  │  │  ├─consul - 注册中心初始化
│  │  │  ├─job - 定时任务初始化
│  │  │  ├─log - 日志初始化
│  │  │  ├─rocketmq - 消息队列初始化
│  │  │  ├─server - 开启grpc服务初始化
│  │  │  └─snowflake - 雪花算法初始化
│  │  ├─middleware - 中间件层
│  │  │  ├─job - 定时任务
│  │  │  └─rocketmq - 消息队列
│  │  ├─model - mysql表
│  │  ├─service - 业务逻辑层
│  │  └─util - 工具包
│  ├─user - 用户service层
│  │  ├─config - 配置文件
│  │  │  └─nacos - 配置中心nacos缓存文件
│  │  │      └─cache
│  │  │          └─config
│  │  ├─dao - 数据读写 mysql redis
│  │  │  ├─mysql
│  │  │  └─redis
│  │  ├─handler - handler层
│  │  ├─initialize - 项目初始化
│  │  │  ├─config - 配置初始化
│  │  │  ├─consul - 注册中心初始化
│  │  │  ├─grpc_client - rpc客户端初始化
│  │  │  ├─job - 定时任务初始化
│  │  │  ├─log - 日志初始化
│  │  │  ├─server - 开启grpc服务初始化
│  │  │  └─snowflake - 雪花算法初始化
│  │  ├─middleware - 中间件层
│  │  │  └─job - 定时任务
│  │  ├─model - mysql表
│  │  ├─service - 业务逻辑层
│  │  └─util - 工具包
│  └─video - 视频service层
│      ├─config - 配置文件
│      │  └─nacos - 配置中心nacos缓存文件
│      │      └─cache
│      │          └─config
│      ├─dao - 数据读写 mysql redis
│      │  ├─mysql
│      │  └─redis
│      ├─handler - handler层
│      ├─initialize - 项目初始化
│      │  ├─config - 配置初始化
│      │  ├─consul - 注册中心初始化
│      │  ├─grpc_client - rpc客户端初始化
│      │  ├─log - 日志初始化
│      │  ├─server - 开启grpc服务初始化
│      │  └─snowflake - 雪花算法初始化
│      ├─model - mysql表
│      ├─service - 业务逻辑层
│      └─util - 工具包
├─test - 项目单元测试
└─web - 项目web层
    ├─config - 配置文件
    │  └─nacos - 配置中心nacos缓存文件
    │      └─cache
    │          └─config
    ├─gin_middleware - gin中间件
    ├─handler - handler层
    ├─initialize - 项目初始化
    │  ├─config - 配置初始化
    │  ├─consul - 注册中心初始化
    │  ├─grpc_client - rpc客户端初始化
    │  ├─log - 日志初始化
    │  ├─oss - 对象存储初始化
    │  └─rocketmq - 消息队列初始化
    ├─middleware - 中间件层
    ├─model - http请求响应结构体
    │  ├─request
    │  └─response
    ├─router - http路由
    ├─service - 业务逻辑层
    └─util - 工具包
```
