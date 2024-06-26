# Toll Microservice

## Invoice Aggregation Service

## Overview
This microservice, **Invoice Aggregation Service**, is designed to process and aggregate distance data for invoicing purposes. It is built with robustness in mind, allowing for efficient handling and storage of distance metrics that are essential for generating accurate invoices.

## Features
- **Data Aggregation**: Aggregate distance data that can be used to calculate costs and generate detailed invoices.
- **REST API**: Expose a simple REST API for data submission and processing.
- **Fault Tolerance**: Implements basic error handling that ensures the service remains available and consistent even when downstream services fail.

## Technologies Used
- **Go (Golang)**: The primary programming language used.
- **Fiber**: A web framework for Go, used for setting up the HTTP server and routes.
- **Confluent Kafka**: Utilized for handling message queues and real-time data streaming.
- **MongoDB**: As a storage solution for the aggregated data.

## Getting Started

### Prerequisites
- Go 1.15+
- MongoDB running on localhost
- Kafka setup and accessible

### Running the Service
Clone the repository and navigate to the project directory:
```bash
git clone https://github.com/yourusername/invoice-aggregation-service.git
cd invoice-aggregation-service
```

# installing gRPC

```
brew install protobuf
```
```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

export PATH="$PATH:$(go env GOPATH)/bin"
```

```
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       types/ptypes.proto
```


## installing Prometheus
```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
```

```
github.com/prometheus/client_golang/prometheus
```

```
../prometheus/prometheus --config.file=.config/prometheus.yml
```