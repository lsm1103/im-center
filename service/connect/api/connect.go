package api

import (
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"im-center/common/xerr"
	"im-center/service/connect/api/internal/handler"
	"im-center/service/connect/internal/config"
	"im-center/service/connect/internal/svc"
	"net/http"

	"github.com/zeromicro/go-zero/rest"
)


func Run(c config.Config, ctx *svc.ServiceContext) *rest.Server {
	server := rest.MustNewServer(ctx.Config.RestConf)
	handler.RegisterHandlers(server, ctx)

	// 自定义错误, 不能拦截auth错误
	httpx.SetErrorHandler(func(err error) (int, interface{}) {
		ctx.Cs.Errorf("SetErrorHandler %s", err.Error())
		switch e := err.(type) {
		case *xerr.CodeError:
			return http.StatusOK, e
		default:
			return http.StatusInternalServerError, nil
		}
	})

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	return server
}
