package main

import (
	"time"

	"github.com/godev/tolls/types"
	"github.com/sirupsen/logrus"
)

type LogginMiddleware struct {
	next Aggregator
}

func NewLogginMiddleware(next Aggregator) *LogginMiddleware {
	return &LogginMiddleware{
		next: next,
	}
}

func (l *LogginMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithField("action", "producing to kafka").WithFields(logrus.Fields{
			"took": time.Since(start),
			"err" : err,
		}).Info("Aggregated distance")
	}(time.Now())

	return l.next.AggregateDistance(distance)
}

func (l *LogginMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		var (
			distance float64
			amount	 float64
		)
		if inv != nil {
			distance =  inv.TotalDistance
			amount = inv.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err" : err,
			"obuID": obuID,
			"amount": amount,
			"distance": distance,
		}).Info("Calculate invoice")
	}(time.Now())

	inv, err = l.next.CalculateInvoice(obuID)
	return
}