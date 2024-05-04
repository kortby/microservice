package main

import (
	"context"

	"github.com/godev/tolls/types"
)

type GRPCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

// there is a difference between transport layer and business layer ( main type )
// transport layer JSON -> types.Distance
// 				   GRPC -> types.AggregateRequest
func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
    distance := types.Distance{
        OBUID: int(req.ObuID),
        Value: req.Value,
        Unix:  req.Unix,
    }
    err := s.svc.AggregateDistance(distance)
    if err != nil {
        return nil, err
    }
    return &types.None{}, nil
}
