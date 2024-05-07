package aggtransport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	stdopentracing "github.com/opentracing/opentracing-go"

	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/godev/tolls/go-kit-example/aggsvc/aggendpoint"
	"github.com/godev/tolls/go-kit-example/aggsvc/aggservice"
)

func errorEncoder (ctx context.Context, err error, w http.ResponseWriter) {
	
}

func NewHTTPHandler(endpoints aggendpoint.Set, logger log.Logger) http.Handler {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	m := http.NewServeMux()
	m.Handle("/aggregate", httptransport.NewServer(
		endpoints.AggregateEndpoint,
		decodeHTTPAggRequest,
		encodeHTTPGenericResponse,
		options...
	))
	m.Handle("/invoice", httptransport.NewServer(
		endpoints.CalculateEndpoint,
		decodeHTTPCalculateRequest,
		encodeHTTPGenericResponse,
		options...
	))
	return m
}

func decodeHTTPAggRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.AggregateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func decodeHTTPCalculateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req aggendpoint.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	return req, err
}

func encodeHTTPGenericResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeHTTPGenericRequest(_ context.Context, r *http.Request, request interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func decodeHTTPAggResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.AggregateResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

func decodeHTTPConcatResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp aggendpoint.CalculateRequest
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}


func NewHTTPClient(instance string, otTracer stdopentracing.Tracer, logger log.Logger) (aggservice.Service, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))

	var options []httptransport.ClientOption
	var aggEndpoint endpoint.Endpoint
	{
		aggEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/aggregate"),
			encodeHTTPGenericRequest,
			decodeHTTPAggResponse,
			options...,
		).Endpoint()
		
		aggEndpoint = limiter(aggEndpoint)
		aggEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Aggregate",
			Timeout: 30 * time.Second,
		}))(aggEndpoint)
	}

	var calcEndpoint endpoint.Endpoint
	{
		calcEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/invoice"),
			encodeHTTPGenericRequest,
			decodeHTTPConcatResponse,
			options...
		).Endpoint()
		
		
		calcEndpoint = limiter(calcEndpoint)
		calcEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:    "Calculate",
			Timeout: 10 * time.Second,
		}))(calcEndpoint)
	}

	return aggendpoint.Set{
		AggregateEndpoint: aggEndpoint,
		CalculateEndpoint: calcEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}