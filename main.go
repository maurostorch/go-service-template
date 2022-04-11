package main

import (
	"os"
	"sync"
	"time"

	"github.com/Kooltra/go-service-template/handler"
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.DebugLevel)
}

func main() {
	t1 := time.Now()

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	sqsClient := sqs.New(sess)

	queueName := os.Getenv("QUEUE_NAME")
	queueUrl, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Error("Error getting queue", queueName, " error", err)
	}

	incoming := make(chan *sqs.Message, 10)
	go receiveMessages(queueUrl.QueueUrl, sqsClient, incoming)

	var wg sync.WaitGroup

	for i := 1; i < 10; i++ {
		wg.Add(1)
		log.Info("Worker", i)
		go func() {
			defer wg.Done()
			handler.New(queueUrl, sqsClient, incoming).Start()
		}()
	}

	duration := time.Now().Sub(t1)
	log.Debug("Load in", duration)

	wg.Wait()
}

func receiveMessages(queueUrl *string, sqsClient *sqs.SQS, incoming chan<- *sqs.Message) {
	for {
		msg, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            queueUrl,
			MaxNumberOfMessages: aws.Int64(*aws.Int64(10)),
			WaitTimeSeconds:     aws.Int64(*aws.Int64(20)),
		})

		if err != nil {
			log.Error("Error receiveing message ", err)
			os.Exit(1)
		}

		for _, message := range msg.Messages {
			incoming <- message
		}
	}
}
