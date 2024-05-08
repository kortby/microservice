package aggservice

import (
	"context"

	"github.com/godev/tolls/types"
)

type Middleware func(Service) Service

type LogginMiddleware struct {
	next Service
}

func NewLogginMiddleware() Middleware {
	return func(next Service) Service {
		return LogginMiddleware{
			next: next,
		}
	}
}

func (mw LogginMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw LogginMiddleware) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, obuID)
}

type InstrumentationMiddleware struct {
	next Service
}


func NewInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return InstrumentationMiddleware{
			next: next,
		}
	}
}

func (mw InstrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw InstrumentationMiddleware) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, obuID)
}