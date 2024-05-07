package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/godev/tolls/go-kit-example/aggsvc/aggservice"
	"github.com/godev/tolls/types"
)

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID float64 `json:"obuID"`
	Unix float64 `json:"unix"`
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type Set struct {
	AggregateEndpoint    endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type AggregateResponse struct {
	Err				error		`json:"err"`
}

type CalculateResponse struct {
	OBUID 			int 		`json:"obuID"`
	TotalDistance	float64 	`json:"totalDistance"`
	TotalAmount		float64 	`json:"totalAmount"`
	Err				error		`json:"err"`
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: float64(dist.OBUID),
		Value: dist.Value,
		Unix: float64(dist.Unix),
	})
	return err
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{OBUID: obuID})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID: 			result.OBUID,
		TotalDistance: 	result.TotalDistance,
		TotalAmount: 	result.TotalAmount,
	}, nil
}

func MakeAggregateEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: int(req.OBUID),
			Value: req.Value,
			Unix: int64(req.Unix),
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeConcatEndpoint(s aggservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		inv, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			OBUID: req.OBUID,
			TotalDistance: inv.TotalDistance,
			TotalAmount: inv.TotalAmount,
			Err: err,
		}, nil
	}
}
