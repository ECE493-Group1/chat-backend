package threads

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Sender  string
	Content string
}

type ThreadRoom struct {
	Title      string
	Messages   []Message
	Members    map[string]bool
	Id         string
	IsPrivate  bool
	UpdateTime time.Time
}

type ThreadManager struct {
	threadRooms map[string]*ThreadRoom
}

func NewThreadRoom(title string, members []string, isPrivate bool) *ThreadRoom {
	memberMap := make(map[string]bool)
	for _, member := range members {
		memberMap[member] = true
	}
	return &ThreadRoom{
		Id:         uuid.New().String(),
		Title:      title,
		Members:    memberMap,
		IsPrivate:  isPrivate,
		UpdateTime: time.Now(),
	}
}

func NewThreadManager() *ThreadManager {
	return &ThreadManager{threadRooms: map[string]*ThreadRoom{}}
}

func (t *ThreadManager) AddThread(newThread *ThreadRoom) {
	t.threadRooms[newThread.Id] = newThread
}

func (t *ThreadManager) GetAllRooms() []*ThreadRoom {
	publicRooms := make([]*ThreadRoom, 0)
	for _, room := range t.threadRooms {
		if !room.IsPrivate {
			publicRooms = append(publicRooms, room)
		}
	}
	return publicRooms
}

func (t *ThreadManager) GetSubscribedRooms(username string) []*ThreadRoom {
	subbedRooms := make([]*ThreadRoom, 0)
	for _, room := range t.threadRooms {
		if room.Members[username] {
			subbedRooms = append(subbedRooms, room)
		}
	}
	return subbedRooms
}

func (t *ThreadManager) AddMessage(id string, m *Message) {
	t.threadRooms[id].Messages = append(t.threadRooms[id].Messages, *m)
	t.threadRooms[id].UpdateTime = time.Now()
}

func (t *ThreadManager) GetThreadMessages(threadId string) []Message {
	return t.threadRooms[threadId].Messages
}

func (t *ThreadManager) GetRoomInfo(threadId string) *ThreadRoom {
	return t.threadRooms[threadId]
}

func (t *ThreadManager) AddMembers(roomId string, members []string) {
	if t.threadRooms[roomId] == nil {
		return
	}
	for _, member := range members {
		// check if member is in room
		t.threadRooms[roomId].Members[member] = true
	}
}

func (t *ThreadManager) RemoveMember(roomId, member string) {
	if t.threadRooms[roomId] == nil {
		return
	}
	delete(t.threadRooms[roomId].Members, member)
}
