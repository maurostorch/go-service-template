package handler

import (
	"fmt"

	"github.com/maurostorch/go-service-template/sqs"
)

type Handler struct{}

func (h *Handler) Handle(msg *sqs.Message) {
	// Write your code here to process messages
	fmt.Println(">>>>>>>", msg)
}
