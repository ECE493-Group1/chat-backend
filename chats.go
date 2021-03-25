package main

import "time"

type Message struct {
	sender  string
	content string
	time    time.Time
}
type ChatRoom struct {
	name     string
	messages []Message
	id       string
}
