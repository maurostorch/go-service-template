package handler

import (
	log "github.com/Sirupsen/logrus"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Handler struct {
	client    *sqs.SQS
	queueUrl  *sqs.GetQueueUrlOutput
	inChannel chan *sqs.Message
}

func New(queueUrl *sqs.GetQueueUrlOutput, client *sqs.SQS, inChannel chan *sqs.Message) *Handler {
	return &Handler{
		client:    client,
		queueUrl:  queueUrl,
		inChannel: inChannel,
	}
}

func (h *Handler) Start() {
	for m := range h.inChannel {
		err := h.handleMessage(m)
		if err == nil {
			h.deleteMessage(m)
		} else {
			return
		}
	}
}

func (h *Handler) handleMessage(msg *sqs.Message) error {
	log.Info("Received message: ", msg)
	return nil
}

func (h *Handler) deleteMessage(msg *sqs.Message) error {
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
