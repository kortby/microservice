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

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler{
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
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (h *HTTPMetricHandler) instrument(next http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		defer func(Start time.Time) {
			latency := time.Since(Start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
			}).Info()
			h.reqLatency.Observe(latency)
		}(time.Now())
		h.reqCounter.Inc()
		next(w, r)
	}
}

func handleGetInvoice(svc Aggregator)  http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Method not supported"})
			return
		}
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu id"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu ID"})
			return
		}
		
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"invoice": invoice})
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Method Post not supported"})
			return
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}