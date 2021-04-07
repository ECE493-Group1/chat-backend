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

type RoomListRequestDTO struct {
	Ids []string `json:"ids"`
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

type KeyWordUpdateDTO struct {
	Content string `json:"message"`
	RoomId  string `json:"chatroom_id"`
}

type MessageListDTO struct {
	Messages []MessageDTO `json:"messages"`
}

type RoomListItemDTO struct {
	RoomId          string   `json:"roomId"`
	Title           string   `json:"title"`
	LastMessageTime int64    `json:"lastMessageTime"`
	Members         []string `json:"members"`
	LastMessage     string   `json:"lastMessage"`
	IsPrivate       bool     `json:"isPrivate"`
}

type RoomListDTO struct {
	Rooms []RoomListItemDTO `json:"rooms"`
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
		RoomId:  new.RoomId,
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

func ToRoomListDTO(rooms []*threads.ThreadRoom) *RoomListDTO {
	roomList := &RoomListDTO{
		Rooms: make([]RoomListItemDTO, len(rooms)),
	}
	for i, room := range rooms {

		memberList := make([]string, 0)

		for member := range room.Members {
			memberList = append(memberList, member)
		}

		lastMessage := "New Thread"
		if len(room.Messages) > 0 {
			lastMessage = room.Messages[len(room.Messages)-1].Content
		}
		roomList.Rooms[i] = RoomListItemDTO{
			RoomId:          room.Id,
			Title:           room.Title,
			LastMessageTime: room.UpdateTime.Unix(),
			Members:         memberList,
			LastMessage:     lastMessage,
			IsPrivate:       room.IsPrivate,
		}
	}
	return roomList
}

func ToKeyWordUpdateDTO(content, roomId string) KeyWordUpdateDTO {
	return KeyWordUpdateDTO{
		Content: content,
		RoomId:  roomId,
	}
}
