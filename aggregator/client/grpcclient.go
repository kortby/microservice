package client

import (
	"github.com/godev/tolls/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)


type GRCPClient struct {
	Endpoint string
	types.AggregatorClient
}

func NewGRCPClient(endpoint string) (*GRCPClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	c := types.NewAggregatorClient(conn)
	return &GRCPClient{
		Endpoint: endpoint,
		AggregatorClient: c,
	}, nil
}



