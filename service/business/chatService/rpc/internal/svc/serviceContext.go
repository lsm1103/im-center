package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"im-center/common/nodeHandler"
	"im-center/service/business/chatService/rpc/internal/config"
	localCache "im-center/service/model/cache"
	"im-center/service/model/database"
	"im-center/service/model/utils"
)

type ServiceContext struct {
	Config config.Config
	RedisClient *redis.Redis

	//model
	ModelHandle     utils.ModelHandle
	JoinTableQuery     database.JoinTableQuery
	FriendModel     database.FriendModel
	GroupModel     database.GroupModel
	UserGroupModel     database.UserGroupModel
	OfflineMsgModel     database.OfflineMsgModel
	SingleMsgModel     database.SingleMsgModel
	GroupMsgModel     database.GroupMsgModel

	//rpc服务
	RpcU *nodeHandler.RpcU
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			//r.Pass = c.Redis.Pass
		}),

		ModelHandle:     	utils.NewModelHandle(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		JoinTableQuery:     database.NewJoinTableQuery(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		FriendModel:     	database.NewFriendModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		GroupModel:         database.NewGroupModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserGroupModel:     database.NewUserGroupModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		OfflineMsgModel:    database.NewOfflineMsgModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		SingleMsgModel:     database.NewSingleMsgModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		GroupMsgModel:      database.NewGroupMsgModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),

		RpcU: nodeHandler.NewRpcU(
			localCache.NewRedisCache(redis.New(c.Redis.Host, func(r *redis.Redis) {
				r.Type = c.Redis.Type
				//r.Pass = c.Redis.Pass
			}) ) ),
	}
}
