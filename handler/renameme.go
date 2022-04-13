package handler

import (
	"fmt"

	"github.com/maurostorch/go-service-template/model"
)

type Handler struct{}

func (h *Handler) Handle(msg *model.Message) {
	// Write your code here to process messages
	fmt.Println(">>>>>>>", msg)
}
