package main

import (

	"github.com/godev/tolls/types"
	"github.com/sirupsen/logrus"
)

const basedPrice = 4.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuid":	distance.OBUID,
		"distance": distance.Value,
		"unix":		distance.Unix,
	}).Info("Aggregating distance")
    return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	dist, err:= i.store.Get(obuID)
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