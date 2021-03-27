package main

import (
	"github.com/google/uuid"
)

type Message struct {
	sender  string
	content string
}

type ChatRoom struct {
	title     string
	messages  []Message
	id        string
	users     []string
	isPrivate bool
}

func NewChatRoom(title string, users []string, isPrivate bool) *ChatRoom {
	return &ChatRoom{
		id:        uuid.New().String(),
		title:     title,
		users:     users,
		isPrivate: isPrivate,
	}
}
