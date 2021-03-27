package dto

type NewRoomDTO struct {
	Title     string   `json:"title"`
	IsPrivate bool     `json:"isPrivate"`
	Members   []string `json:"members"`
}

type NewRoomResponseDTO struct {
	Id string `json:"id"`
}

type MessageDTO struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomId   string `json:"roomId"`
}
