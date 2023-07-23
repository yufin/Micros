package server

import (
	gv1 "brillinkmicros/api/graph/v1"
	rcv1 "brillinkmicros/api/rc/v1"
	"brillinkmicros/internal/conf"
	"brillinkmicros/internal/data"
	"brillinkmicros/internal/midware"
	"brillinkmicros/internal/service"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server,
	data *data.Data,
	confData *conf.Data,
	rss *service.RcServiceServicer,
	tgs *service.TreeGraphServiceServicer,
	rrs *service.RcRdmServiceServicer,
	logger log.Logger) *http.Server {
	var opts = []http.ServerOption{
		// 一个请求进入时的处理顺序为 Middleware 注册的顺序，而响应返回的处理顺序为注册顺序的倒序
		http.Middleware(
			recovery.Recovery(),
			logging.Server(logger),
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
	openAPIHandler := openapiv2.NewHandler()
	srv.HandlePrefix("/q/", openAPIHandler)
	log.Infof("SwaggerUI DOC: %s/q/swagger-ui/", c.Http.Addr)
	rcv1.RegisterRcServiceHTTPServer(srv, rss)
	rcv1.RegisterRcRdmServiceHTTPServer(srv, rrs)
	gv1.RegisterTreeGraphServiceHTTPServer(srv, tgs)
	return srv
}
