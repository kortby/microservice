package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/godev/tolls/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type InvoiceHandler struct {
	client client.Client
}

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "the list address of the HTTP server")
	flag.Parse()
	var (

		client = client.NewHTTPClient("http://127.0.0.1:3000")
		invHandler = NewInvoiceHandler(client)
	)
	http.HandleFunc("/invoice", makeAPIFunc(invHandler.handleGetInvoice))
	logrus.Infof("gatway HTTP server running on port %s", *listenAddr)
	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	// access aggre client
	inv, err := h.client.GetInvoice(context.Background(), 3322)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}