package domain

import (
	"github.com/google/uuid"
)

type Message struct {
	FromEmail string
	ToEmail   string
	Subject   string
	Content   string
}

type MessageID string

func NewMessageID() MessageID {
	return MessageID(uuid.NewString())
}

func (m MessageID) String() string {
	return string(m)
}
