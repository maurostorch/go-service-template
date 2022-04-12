package main

import (
	"context"
	"os"
	"sync"

	"github.com/Kooltra/go-service-template/handler"
	"github.com/Kooltra/go-service-template/sqs"
	log "github.com/Sirupsen/logrus"
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
