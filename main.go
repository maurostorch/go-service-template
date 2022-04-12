package main

import (
	"context"
	"fmt"
	"os"
	"sync"

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
	receiver := sqs.NewReceiver(queueName, &ctx)
	sqsContext, channel := receiver.Start()
	log.Info("Connected to queue ", queueName)
	var wg sync.WaitGroup
	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			sqs.NewHandler(sqsContext, channel, handleMessage).Start()
		}()
	}
	wg.Wait()
}

func handleMessage(msg *sqs.Message) {
	fmt.Println(">>>>>>>", msg)
}
