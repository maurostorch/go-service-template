package sqs

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageHandler interface {
	Handle(msg *Message)
}
type Handler interface {
	Start()
}

type handler struct {
	ctx       context.Context
	inChannel chan *sqs.Message
	handler   *MessageHandler
}

type Message struct {
	body string
}

func NewHandler(context context.Context, inChannel chan *sqs.Message, msgHandler *MessageHandler) Handler {
	return &handler{
		ctx:       context,
		inChannel: inChannel,
		handler:   msgHandler,
	}
}

func (h *handler) Start() {
	for {
		select {
		case m := <-h.inChannel:
			err := h.handleMessage(m)
			if err == nil {
				h.deleteMessage(m)
			} else {
				log.Error("Error processing message ", err)
			}
		case <-h.ctx.Done():
			return
		}
	}
}

func (h *handler) handleMessage(msg *sqs.Message) error {
	log.Debug("Received message: ", msg)
	(*h.handler).Handle(&Message{
		body: *msg.Body,
	})
	return nil
}

func (h *handler) deleteMessage(msg *sqs.Message) error {
	sqsContext := h.ctx.Value(contextKey).(*SQSClient)
	_, err := sqsContext.client.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      sqsContext.QueueUrl.QueueUrl,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Error("Unable to delete message ", msg)
	} else {
		log.Debug("Message delete ", msg)
	}
	return err
}
