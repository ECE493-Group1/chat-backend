package threads

import (
	"github.com/google/uuid"
)

type Message struct {
	Sender  string
	Content string
}

type ThreadRoom struct {
	Title     string
	Messages  []Message
	Members   []string
	Id        string
	IsPrivate bool
}

type ThreadManager struct {
	threadRooms map[string]*ThreadRoom
}

func NewThreadRoom(title string, users []string, isPrivate bool) *ThreadRoom {
	return &ThreadRoom{
		Id:        uuid.New().String(),
		Title:     title,
		Members:   users,
		IsPrivate: isPrivate,
	}
}

func NewThreadManager() *ThreadManager {
	return &ThreadManager{threadRooms: map[string]*ThreadRoom{}}
}

func (t *ThreadManager) AddThread(newThread *ThreadRoom) {
	t.threadRooms[newThread.Id] = newThread
}

func (t *ThreadManager) GetAllRooms() ([]string, []string) {
	titles := make([]string, len(t.threadRooms))
	ids := make([]string, len(t.threadRooms))

	for k, v := range t.threadRooms {
		titles[0] = v.Title
		ids[0] = k
	}
	return titles, ids
}

func (t *ThreadManager) AddMessage(id string, m *Message) {
	t.threadRooms[id].Messages = append(t.threadRooms[id].Messages, *m)
}

func (t *ThreadManager) GetThreadMessages(threadId string) []Message {
	return t.threadRooms[threadId].Messages
}

func (t *ThreadManager) GetRoomInfo(threadId string) *ThreadRoom {
	return t.threadRooms[threadId]
}
