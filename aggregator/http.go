package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/godev/tolls/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err error
}

type HTTPMetricHandler struct {
	errCounter prometheus.Counter
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func makeHTTPHandlerFunc(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Error()})
			}
		}
	}
}

// Error implement Error interface
func (e APIError) Error() string {
	return e.Err.Error()
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler{
	errCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "error_counter"),
		Name: "aggregator",
	})
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
		Name: "aggregator",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
		Name: "aggregator",
		Buckets: prometheus.LinearBuckets(0.1, 0.5, 5) ,
	})
	return &HTTPMetricHandler{
		errCounter: errCounter,
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func (w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(Start time.Time) {
			latency := time.Since(Start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err": err,
			}).Info()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errCounter.Inc()
			}
		}(time.Now())
		if err = next(w, r); err != nil {
			fmt.Println("--error---")
			return err
		}
		err = next(w, r)
		return err
	}
}

func handleGetInvoice(svc Aggregator)  HTTPFunc {
	return func (w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			return APIError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("missing obu id"),
			}
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("invalid obu id %s", values[0]),
			}
		}
		
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err: err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]any{"invoice": invoice})
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func (w http.ResponseWriter, r *http.Request) error {
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("invalid HTTP method %s", r.Method),
			}
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err: fmt.Errorf("Failed to decode response body: %s", err),
			}
		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err: err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}