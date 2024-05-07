package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/godev/tolls/aggregator/client"
	"github.com/godev/tolls/types"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	flag.Parse()

	var (
		store = makeStore()
		svc = NewInvoiceAggregator(store)
		grpcListenAddr= os.Getenv("AGG_GRPC_ENDPOINT")
		httpListenAddr= os.Getenv("AGG_HTTP_ENDPOINT")
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogginMiddleware(svc)

	go makeGRPCTransport(grpcListenAddr, svc)
	time.Sleep(time.Second * 3)
	c, err := client.NewGRCPClient(grpcListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 344.34,
		Unix: time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
	makeHTTPTransport(httpListenAddr, svc)
}

func makeStore() Storer {
	storeType := os.Getenv("AGG_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store %s", storeType)
		return nil
	}
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRCP transport running on port ", listenAddr)
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// register to grpc package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	aggMatricsHandler := NewHTTPMetricHandler("aggregate")
	invMatricsHandler := NewHTTPMetricHandler("invoice")
	http.HandleFunc("/aggregate", aggMatricsHandler.instrument(handleAggregate(svc)))
	http.HandleFunc("/invoice", invMatricsHandler.instrument(handleGetInvoice(svc)))
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("HTTP transport running on port ", listenAddr)
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}



func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}