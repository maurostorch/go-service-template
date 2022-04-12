package sqs

import (
	"os"

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

func (r *Receiver) Start() chan *sqs.Message {
	incoming := make(chan *sqs.Message, 10)
	go r.receiveMessages(incoming)
	return incoming
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
			os.Exit(1)
		}

		for _, message := range msg.Messages {
			incoming <- message
		}
	}
}
