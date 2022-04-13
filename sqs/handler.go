package sqs

import (
	"context"

	"github.com/maurostorch/go-service-template/model"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageHandler interface {
	Handle(msg *model.Message)
}
type Handler interface {
	Start(context.Context)
}

type handler struct {
	client    *sqs.SQS
	queueUrl  *sqs.GetQueueUrlOutput
	inChannel chan *sqs.Message
	handler   *MessageHandler
}

func NewHandler(inChannel chan *sqs.Message, client *sqs.SQS, queueUrl *sqs.GetQueueUrlOutput, msgHandler *MessageHandler) Handler {
	return &handler{
		client:    client,
		queueUrl:  queueUrl,
		inChannel: inChannel,
		handler:   msgHandler,
	}
}

func (h *handler) Start(ctx context.Context) {
	for {
		select {
		case m := <-h.inChannel:
			err := h.handleMessage(m)
			if err == nil {
				h.deleteMessage(m)
			} else {
				log.Error("Error processing message ", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (h *handler) handleMessage(msg *sqs.Message) error {
	log.Debug("Received message: ", msg)
	(*h.handler).Handle(&model.Message{
		Body: *msg.Body,
	})
	return nil
}

func (h *handler) deleteMessage(msg *sqs.Message) error {
	_, err := h.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      h.queueUrl.QueueUrl,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Error("Unable to delete message ", msg)
	} else {
		log.Debug("Message delete ", msg)
	}
	return err
}
