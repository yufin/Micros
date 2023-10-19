package data

import (
	"github.com/pkg/errors"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"micros-graph/internal/conf"
)

type NebulaDb struct {
	Pool     *nebula.ConnectionPool
	addrs    []nebula.HostAddress
	pwd      string
	usr      string
	useHttp2 bool
}

func (n *NebulaDb) getSession() (*nebula.Session, error) {
	session, err := n.Pool.GetSession(n.usr, n.pwd)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (n *NebulaDb) NewSessionPool(usr string, pwd string, space string) (*nebula.SessionPool, error) {
	config, err := nebula.NewSessionPoolConf(
		n.usr,
		n.pwd,
		n.addrs,
		space,
		nebula.WithHTTP2(n.useHttp2),
	)
	if err != nil {
		return nil, err
	}
	return nebula.NewSessionPool(*config, nebula.DefaultLogger{})
}

func (*NebulaDb) checkResultSet(res *nebula.ResultSet) error {
	if !res.IsSucceed() {
		return errors.Errorf("ErrorCode: %v, ErrorMsg: %s", res.GetErrorCode(), res.GetErrorMsg())
	}
	return nil
}

func NewNebulaDb(c *conf.Data) (*NebulaDb, func(), error) {
	hostAddr := nebula.HostAddress{
		Host: c.NebulaDb.Addr,
		Port: int(c.NebulaDb.Port),
	}
	hostList := []nebula.HostAddress{hostAddr}
	poolConfig := nebula.GetDefaultConf()
	poolConfig.UseHTTP2 = c.NebulaDb.UseHttp2

	log := nebula.DefaultLogger{}
	pool, err := nebula.NewConnectionPool(hostList, poolConfig, log)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		pool.Close()
	}
	return &NebulaDb{
		Pool:     pool,
		addrs:    hostList,
		usr:      c.NebulaDb.User,
		pwd:      c.NebulaDb.Password,
		useHttp2: c.NebulaDb.UseHttp2,
	}, cleanup, nil
}
