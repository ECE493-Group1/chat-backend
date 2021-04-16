package dto

import (
	"testing"

	"catchat.com/models/threads"
	"github.com/stretchr/testify/assert"
)

func TestToMessageListDTO(t *testing.T) {
	messages := []threads.Message{
		{
			Content: "Test",
			Sender:  "A",
			RoomId:  "1",
		},
	}
	listDTO := ToMessageListDTO(messages)

	assert.Equal(t, listDTO.Messages[0].Content, "Test")
	assert.Equal(t, listDTO.Messages[0].Username, "A")
}

func TestToMessage(t *testing.T) {
	messageDTO := NewMessageDTO{
		Content:  "Test",
		Username: "A",
		RoomId:   "1",
	}
	message := ToMessage(messageDTO)

	assert.Equal(t, message.Content, "Test")
	assert.Equal(t, message.Sender, "A")
	assert.Equal(t, message.RoomId, "1")
}

func TestToRoomDTO(t *testing.T) {
	threadRoom := threads.NewThreadRoom("Test", []string{"a", "b"}, false)
	roomDTO := ToRoomDTO(threadRoom)
	assert.Equal(t, len(roomDTO.Members), 2)
	assert.Equal(t, roomDTO.Title, "Test")
}

func TestToRoomListDTO(t *testing.T) {
	threadRoom1 := threads.NewThreadRoom("Test", []string{"a", "b"}, false)
	threadRoom2 := threads.NewThreadRoom("Test", []string{"a", "b"}, false)
	threadRoom3 := threads.NewThreadRoom("Test", []string{"a", "b"}, false)
	roomList := []*threads.ThreadRoom{threadRoom1, threadRoom2, threadRoom3}
	roomDTO := ToRoomListDTO(roomList)

	assert.Equal(t, len(roomDTO.Rooms), 3)
}

func TestKWUpdateDTO(t *testing.T) {
	kwDTO := ToKeyWordUpdateDTO("Test", "1")
	assert.Equal(t, kwDTO.Content, "Test")
	assert.Equal(t, kwDTO.RoomId, "1")
}
