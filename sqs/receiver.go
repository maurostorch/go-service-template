package sqs

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

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
	workersCtx, cancelWorkers := context.WithCancel(ctx)
	receiverCtx, cancelReceiver := context.WithCancel(ctx)
	incoming := make(chan *sqs.Message, workers)
	var wg sync.WaitGroup
	for i := 1; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Debug("Starting ", id)
			NewHandler(incoming, r.Client, r.QueueUrl, &handler).Start(workersCtx)
			log.Debug("Stopping ", id)
		}(i)
	}
	go func() {
		defer func() {
			cancelWorkers()
			wg.Done()
		}()
		wg.Add(1)
		r.receiveMessages(receiverCtx, incoming)
	}()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT)
	go func(signals chan os.Signal) {
		s := <-signals
		log.Info(s, " Signal received. Waiting receivers to finish...")
		cancelReceiver()
	}(signalCh)
	wg.Wait()
}

func (r *Receiver) receiveMessages(ctx context.Context, incoming chan *sqs.Message) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := r.Client.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:            r.QueueUrl.QueueUrl,
				MaxNumberOfMessages: aws.Int64(10),
				WaitTimeSeconds:     aws.Int64(10),
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
}
