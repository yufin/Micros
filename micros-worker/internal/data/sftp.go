package data

import (
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"micros-worker/internal/conf"
	"sync"
)

type SftpClientPool struct {
	uri          string
	poolSize     int
	clientConfig *ssh.ClientConfig
	clients      []*sftp.Client
	mu           sync.Mutex
	cond         *sync.Cond
}

func (p *SftpClientPool) newClient() (*sftp.Client, error) {
	conn, err := ssh.Dial("tcp", p.uri, p.clientConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return sftpClient, nil
}

func (p *SftpClientPool) checkClient(client *sftp.Client) bool {
	_, err := client.Stat("/")
	return err == nil
}

func (p *SftpClientPool) GetClient() (*sftp.Client, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for len(p.clients) == 0 {
		p.cond.Wait()
	}

	client := p.clients[0]
	p.clients = p.clients[1:]

	if p.checkClient(client) {
		return client, nil
	}

	client.Close()
	return p.newClient()
}

func (p *SftpClientPool) ReturnClient(client *sftp.Client) {
	p.mu.Lock()
	p.clients = append(p.clients, client)
	p.mu.Unlock()
	p.cond.Signal()
}

func NewSftpClientPool(c *conf.Data) (*SftpClientPool, func(), error) {
	cliConf := &ssh.ClientConfig{
		User:            c.VzoomSftp.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(c.VzoomSftp.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10000000000,
	}

	pool := SftpClientPool{
		clientConfig: cliConf,
		uri:          c.VzoomSftp.Uri,
		clients:      make([]*sftp.Client, 0),
		poolSize:     int(c.VzoomSftp.PoolSize),
	}
	pool.cond = sync.NewCond(&pool.mu)

	cleanUp := func() {
		for _, cli := range pool.clients {
			if cli != nil {
				cli.Close()
			}
		}
	}

	for i := 0; i < pool.poolSize; i++ {
		client, err := pool.newClient()
		if err != nil {
			return nil, cleanUp, errors.WithStack(err)
		}
		pool.clients = append(pool.clients, client)
	}

	return &pool, cleanUp, nil
}
