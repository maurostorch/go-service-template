package handler

import (
	"fmt"

	"github.com/Kooltra/go-service-template/sqs"
)

type Handler struct{}

func (h *Handler) Handle(msg *sqs.Message) {
	// Write your code here to process messages
	fmt.Println(">>>>>>>", msg)
}
