package main

import (
	"fmt"

	"github.com/godev/tolls/types"
)

type Aggregator interface {
	AggregateDistances(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistances(distance types.Distance) error {
    fmt.Println("processing ..", distance)
    return i.store.Insert(distance)
}