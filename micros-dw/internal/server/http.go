package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"micros-dw/api/dw/v2"
	"micros-dw/internal/conf"
	"micros-dw/internal/service/dw/v2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	dws *v2.DwServiceServicer,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		// 一个请求进入时的处理顺序为 Middleware 注册的顺序，而响应返回的处理顺序为注册顺序的倒序
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	openAPIHandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIHandler)
	log.Infof("SwaggerUI DOC: %s/q/swagger-ui/", c.Http.Addr)
	dwV2.RegisterDwServiceHTTPServer(srv, dws)
	return srv
}
