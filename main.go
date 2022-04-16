package main

import (
	"context"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maurostorch/go-service-template/handler"
	"github.com/maurostorch/go-service-template/sqs"
)

func init() {
	_, err := os.Create("/tmp/LIVE")
	if err != nil {
		os.Exit(1)
	}
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	ctx := context.Background()
	queueName := os.Getenv("QUEUE_NAME")
	sqs := sqs.NewReceiver(queueName)
	log.Info("Connected to queue ", queueName)
	sqs.Start(ctx, &handler.Handler{}, 10)
}
