package main

import (
	"time"

	"github.com/godev/tolls/types"
	"github.com/sirupsen/logrus"
)

type LogginMiddleware struct {
	next DataProducer
}

func NewLogginMiddleware(next DataProducer) *LogginMiddleware {
	return &LogginMiddleware{
		next: next,
	}
}

func (l *LogginMiddleware) ProduceData(data types.OBUData) error {
	defer func(start time.Time) {
		logrus.WithField("producing to kafka", logrus.Fields{
			"obuID": data.OBUID,
			"lat" : data.Lat,
			"lng" : data.Lng,
			"took": time.Since(start),
		}).Info("producing to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}