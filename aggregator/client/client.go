package client

import (
	"context"

	"github.com/godev/tolls/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
	GetInvoice(context.Context, int) (*types.Invoice, error)
}