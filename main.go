package main

import (
	"context"
	"os"
	"sync"

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
	channel := receiver.Start()
	log.Info("Connected to queue ", queueName)
	var handler sqs.MessageHandler = &handler.Handler{}
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Debug("Starting ", id)
			sqs.NewHandler(channel, receiver.Client, receiver.QueueUrl, &handler).Start(ctx)
			log.Debug("Stopping ", id)
		}(i)
	}
	wg.Wait()
}
