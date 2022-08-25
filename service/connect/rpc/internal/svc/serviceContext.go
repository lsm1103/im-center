package svc

import (
	"im-center/service/connect/core"
	"im-center/service/connect/rpc/internal/config"
)

type ServiceContext struct {
	Config config.Config
	CS *core.ConnectServer
}

func NewServiceContext(c config.Config, cs *core.ConnectServer) *ServiceContext {
	return &ServiceContext{
		Config: c,
		CS: cs,
	}
}
