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

func (mw LogginMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw LogginMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
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

func (mw InstrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw InstrumentationMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}