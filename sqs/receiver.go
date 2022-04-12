package sqs

import (
	"context"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Receiver struct {
	Client   *sqs.SQS
	QueueUrl *sqs.GetQueueUrlOutput
}

func NewReceiver(queueName string) Receiver {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})

	sqsClient := sqs.New(sess)

	queueUrl, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: &queueName,
	})
	if err != nil {
		log.Error("Error getting queue", queueName, " error", err)
		os.Exit(1)
	}
	return Receiver{
		Client:   sqsClient,
		QueueUrl: queueUrl,
	}
}

func (r *Receiver) Start(ctx context.Context, handler MessageHandler, workers int) {
	localCtx, cancel := context.WithCancel(ctx)
	incoming := make(chan *sqs.Message, workers)
	var wg sync.WaitGroup
	for i := 1; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Debug("Starting ", id)
			NewHandler(incoming, r.Client, r.QueueUrl, &handler).Start(localCtx)
			log.Debug("Stopping ", id)
		}(i)
	}
	r.receiveMessages(incoming)
	cancel()
	wg.Wait()
}

func (r *Receiver) receiveMessages(incoming chan *sqs.Message) {
	for {
		msg, err := r.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            r.QueueUrl.QueueUrl,
			MaxNumberOfMessages: aws.Int64(*aws.Int64(10)),
			WaitTimeSeconds:     aws.Int64(*aws.Int64(20)),
		})

		if err != nil {
			log.Error("Error receiveing message ", err)
			return
		}

		for _, message := range msg.Messages {
			incoming <- message
		}
	}
}
