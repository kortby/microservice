package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/godev/tolls/aggregator/client"
	"github.com/godev/tolls/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "listen address of the HTTP transport handler server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "listen address of the GRPC transport handler server")
	flag.Parse()

	var (
		store = NewMemoryStore()
		svc = NewInvoiceAggregator(store)
	)
	svc = NewLogginMiddleware(svc)

	go makeGRPCTransport(*grpcListenAddr, svc)
	time.Sleep(time.Second * 3)
	c, err := client.NewGRCPClient(*grpcListenAddr)
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
	makeHTTPTransport(*httpListenAddr, svc)
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
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetInvoice(svc Aggregator)  http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {

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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}