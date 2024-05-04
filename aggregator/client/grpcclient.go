package client

import (
	"context"

	"github.com/godev/tolls/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type GRCPClient struct {
	Endpoint string
	client 	 types.AggregatorClient
}

func NewGRCPClient(endpoint string) (*GRCPClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRCPClient{
		Endpoint: endpoint,
		client: c,
	}, nil
}

func (c *GRCPClient) Aggregate(ctx context.Context, req *types.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}

