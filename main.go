package main

import (
	"context"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/maurostorch/go-service-template/handler"
	"github.com/maurostorch/go-service-template/sqs"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	ctx := context.Background()
	queueName := os.Getenv("QUEUE_NAME")
	receiver := sqs.NewReceiver(queueName)
	log.Info("Connected to queue ", queueName)
	receiver.Start(ctx, &handler.Handler{}, 10)
}
