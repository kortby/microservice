package main

import (
	"log"

	"github.com/godev/tolls/aggregator/client"
)

const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogginMiddleware(svc)

	// httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGRCPClient(aggregatorEndpoint)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}