package main

import (
	"time"

	"github.com/godev/tolls/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	errCounterAgg prometheus.Counter
	errCounterCal prometheus.Counter
	reqCounterAgg prometheus.Counter
	reqCounterCal prometheus.Counter
	reqLatencyAgg prometheus.Histogram
	reqLatencyCal prometheus.Histogram
	next	   Aggregator
}

type LogginMiddleware struct {
	next Aggregator
}

// Global metric definitions
var (
	errCounterAgg = promauto.NewCounter(prometheus.CounterOpts{
        Namespace: "aggregator_error_counter",
        Name: "aggregate",
    })
    errCounterCal = promauto.NewCounter(prometheus.CounterOpts{
        Namespace: "aggregator_error_counter",
        Name: "calculate",
    })
    reqCounterAgg = promauto.NewCounter(prometheus.CounterOpts{
        Namespace: "aggregator_request_counter",
        Name: "aggregate",
    })
    reqCounterCal = promauto.NewCounter(prometheus.CounterOpts{
        Namespace: "aggregator_request_counter",
        Name: "calculate",
    })
    reqLatencyAgg = promauto.NewHistogram(prometheus.HistogramOpts{
        Namespace: "aggregator_request_latency",
        Name: "aggregate",
        Buckets: prometheus.LinearBuckets(0.1, 0.5, 5),
    })
    reqLatencyCal = promauto.NewHistogram(prometheus.HistogramOpts{
        Namespace: "calculator_request_latency",
        Name: "calculate",
        Buckets: prometheus.LinearBuckets(0.1, 0.5, 5),
    })
)

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
    return &MetricsMiddleware{
        next: next,
        errCounterAgg: errCounterAgg,
        errCounterCal: errCounterCal,
		reqCounterAgg: reqCounterAgg,
        reqCounterCal: reqCounterCal,
        reqLatencyAgg: reqLatencyAgg,
        reqLatencyCal: reqLatencyCal,
    }
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}


func (m *MetricsMiddleware) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCal.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCal.Inc()
		if err != nil {
			m.errCounterCal.Inc()
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
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