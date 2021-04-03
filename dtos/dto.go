package dto

import "catchat.com/models/threads"

type NewRoomDTO struct {
	Title     string   `json:"title"`
	IsPrivate bool     `json:"isPrivate"`
	Members   []string `json:"members"`
}

type RoomDTO struct {
	Title     string   `json:"title"`
	IsPrivate bool     `json:"isPrivate"`
	Members   []string `json:"members"`
}

type UpdateRoomDTO struct {
	Id      string   `json:"id"`
	Members []string `json:"members"`
}

type LeaveRoomDTO struct {
	Id     string `json:"id"`
	Member string `json:"member"`
}

type RoomRequestDTO struct {
	Id string `json:"id"`
}

type NewMessageDTO struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
}

type MessageDTO struct {
	Content  string `json:"content"`
	Username string `json:"username"`
}

type MessageListDTO struct {
	Messages []MessageDTO `json:"messages"`
}

func ToMessageListDTO(m []threads.Message) *MessageListDTO {
	messages := make([]MessageDTO, len(m))
	for i := range m {
		messages[i].Content = m[i].Content
		messages[i].Username = m[i].Sender
	}
	return &MessageListDTO{Messages: messages}
}

func ToMessage(new NewMessageDTO) *threads.Message {
	return &threads.Message{
		Content: new.Content,
		Sender:  new.Username,
	}
}

func ToRoomDTO(t *threads.ThreadRoom) *RoomDTO {
	members := make([]string, len(t.Members))
	i := 0
	for member := range t.Members {
		members[i] = member
		i++
	}
	return &RoomDTO{
		Title:     t.Title,
		Members:   members,
		IsPrivate: t.IsPrivate,
	}
}
