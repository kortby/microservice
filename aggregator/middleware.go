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