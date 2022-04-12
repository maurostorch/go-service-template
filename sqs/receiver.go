package sqs

import (
	"context"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Key string

const (
	contextKey Key = "contextKey"
)

type Receiver interface {
	Start() (context.Context, chan *sqs.Message)
}

type receiver struct {
	client *sqs.SQS
	ctx    context.Context
}

type SQSClient struct {
	client   *sqs.SQS
	QueueUrl *sqs.GetQueueUrlOutput
}

func NewReceiver(queueName string, parentCtx *context.Context) Receiver {
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
	return &receiver{
		client: sqsClient,
		ctx: context.WithValue(*parentCtx, contextKey, &SQSClient{
			client:   sqsClient,
			QueueUrl: queueUrl,
		}),
	}
}

func (r *receiver) Start() (context.Context, chan *sqs.Message) {
	incoming := make(chan *sqs.Message, 10)
	go r.receiveMessages(incoming)
	return r.ctx, incoming
}

func (r *receiver) receiveMessages(incoming chan *sqs.Message) {
	for {
		msg, err := r.client.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            r.ctx.Value(contextKey).(*SQSClient).QueueUrl.QueueUrl,
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
