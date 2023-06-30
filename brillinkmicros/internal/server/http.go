package server

import (
	v1 "brillinkmicros/api/rc/v1"
	"brillinkmicros/internal/conf"
	"brillinkmicros/internal/data"
	"brillinkmicros/internal/midware"
	"brillinkmicros/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, data *data.Data, confData *conf.Data, rss *service.RcServiceService, logger log.Logger) *http.Server {
	openAPIhandler := openapiv2.NewHandler()
	var opts = []http.ServerOption{
		// 一个请求进入时的处理顺序为 Middleware 注册的顺序，而响应返回的处理顺序为注册顺序的倒序
		http.Middleware(
			recovery.Recovery(),
			midware.BlAuth(confData, data),
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
	srv.HandlePrefix("/q/", openAPIhandler)
	v1.RegisterRcServiceHTTPServer(srv, rss)
	return srv
}
