package data

import "micros-worker/internal/conf"
import "go.temporal.io/sdk/client"

type TemporalClient struct {
	Client      *client.Client
	InfraClient *client.Client
}

func NewTemporalClient(c *conf.Data) (*TemporalClient, func(), error) {
	cli, err := client.Dial(
		client.Options{
			HostPort: c.Temporal.ClientUri,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	infraCli, err := client.Dial(
		client.Options{
			HostPort:  c.Temporal.ClientUri,
			Namespace: "infra",
		},
	)
	if err != nil {
		return nil, nil, err
	}

	tc := TemporalClient{
		Client:      &cli,
		InfraClient: &infraCli,
	}

	return &tc, func() {
		cli.Close()
		infraCli.Close()
	}, nil
}
