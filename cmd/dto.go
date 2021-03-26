package main

type NewRoomDTO struct {
	Name string `json:"roomname"`
}

type MessageDTO struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
}
