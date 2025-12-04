package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type LogEvent struct {
	Service   string    `json:"service"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Global variables for AWS Kinesis.
var (
	kinesisClient *kinesis.Client
	streamName    = "log_stream"
)

// initAWS sets up the AWS SDK and Kinesis client.
func initAWS() {
	// Load AWS configuration (region, credentials) from:
	// - aws configure
	// - environment variables
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS SDK config: %v", err)
	}

	kinesisClient = kinesis.NewFromConfig(cfg)
	log.Println("AWS Kinesis client initialised")
}

// sendToKinesis sends a single LogEvent to the Kinesis data stream.
func sendToKinesis(event LogEvent) error {
	// Convert the struct to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// PartitionKey decides which shard the record goes to.
	// Using the service name is simple and reasonable.
	_, err = kinesisClient.PutRecord(context.Background(), &kinesis.PutRecordInput{
		StreamName:   aws.String(streamName),
		Data:         data,
		PartitionKey: aws.String(event.Service),
	})

	return err
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var event LogEvent
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// If client doesn’t send a timestamp, set it to "now" in UTC
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Send to Kinesis
	if err := sendToKinesis(event); err != nil {
		log.Printf("failed to send to Kinesis: %v", err)
		http.Error(w, "failed to queue log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "queued"})
}

func main() {
	// Initialise AWS & Kinesis client
	initAWS()

	// Set up HTTP route
	http.HandleFunc("/logs", logHandler)

	log.Println("Starting log API on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
