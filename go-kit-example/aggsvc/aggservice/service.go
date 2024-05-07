package aggservice

import (
	"context"
	"fmt"

	"github.com/godev/tolls/types"
)

const basedPrice = 4.15

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

func NewBasicService(store Storer) Service {
	return &BasicService{
		store: store,
	}
}

func (svc *BasicService) Aggregate(_ context.Context, dist types.Distance) error {
	fmt.Println("Aggregate using go-kit")
    return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(_ context.Context, obuID int) (*types.Invoice, error) {
	dist, err:= svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	inv := &types.Invoice{
		OBUID: obuID,
		TotalDistance: dist,
		TotalAmount: basedPrice * dist,
	}
	return inv, nil
}

func New() Service {
	var svc Service
	svc = NewBasicService(NewMemoryStore())
	svc = NewLogginMiddleware()(svc)
	svc = NewInstrumentationMiddleware()(svc)
	return svc
}