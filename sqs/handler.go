package sqs

import (
	"context"

	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Handler interface {
	Start()
}

type handler struct {
	ctx             context.Context
	inChannel       chan *sqs.Message
	handlerFunction func(msg *Message)
}

type Message struct {
	body string
}

func NewHandler(context context.Context, inChannel chan *sqs.Message, f func(msg *Message)) Handler {
	return &handler{
		ctx:             context,
		inChannel:       inChannel,
		handlerFunction: f,
	}
}

func (h *handler) Start() {
	for m := range h.inChannel {
		err := h.handleMessage(m)
		if err == nil {
			h.deleteMessage(m)
		} else {
			log.Error("Error processing message ", err)
		}
	}
}

func (h *handler) handleMessage(msg *sqs.Message) error {
	log.Debug("Received message: ", msg)
	h.handlerFunction(&Message{
		body: *msg.Body,
	})
	return nil
}

func (h *handler) deleteMessage(msg *sqs.Message) error {
	sqsContext := h.ctx.Value("sqs").(*SQSClient)
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
