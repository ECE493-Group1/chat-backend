package main

import (
	"github.com/google/uuid"
)

type Message struct {
	sender  string
	content string
}
type ChatRoom struct {
	name     string
	messages []Message
	id       string
}

func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{id: uuid.New().String(), name: name}
}
