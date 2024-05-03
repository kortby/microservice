package main

import (
	"time"

	"github.com/godev/tolls/types"
	"github.com/sirupsen/logrus"
)

type LogginMiddleware struct {
	next CalculatorServicer
}

func NewLogginMiddleware(next CalculatorServicer) *LogginMiddleware {
	return &LogginMiddleware{
		next: next,
	}
}

func (m *LogginMiddleware) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err": err,
			"distance": dist,
		}).Info("calculate distance")
	}(time.Now())
	dist, err = m.next.CalculateDistance(data)
	
	return
}