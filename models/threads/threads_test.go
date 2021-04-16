package threads

import (
	"testing"
)

func TestInit(t *testing.T) {
	newManager := NewThreadManager()
	if newManager == nil {
		t.Errorf("Could not initialize manager\n")
	}
}

func TestAddThread(t *testing.T) {
	title := "test"
	members := []string{"a", "b", "c"}
	newManager := NewThreadManager()
	newRoom := NewThreadRoom(title, members, false)

	id := newRoom.Id
	newManager.AddThread(newRoom)

	room := newManager.GetRoomInfo(id)

	for _, member := range members {
		if !room.Members[member] {
			t.Errorf("Missing Member\n")
		}
	}
	if id != room.Id {
		t.Error()
	}
}

func TestAddMessage(t *testing.T) {
	title := "test"
	members := []string{"a", "b", "c"}
	newManager := NewThreadManager()
	newRoom := NewThreadRoom(title, members, false)

	id := newRoom.Id
	newManager.AddThread(newRoom)

	message := &Message{
		RoomId:  id,
		Content: "TEST",
		Sender:  "a",
	}
	newManager.AddMessage(id, message)

	messages := newManager.GetThreadMessages(id)

	if messages[0].Content != "TEST" {
		t.Errorf("Invalid test message\n")
	}
}

func TestGetAllRooms(t *testing.T) {
	title := "test"
	members := []string{"a"}
	newManager := NewThreadManager()
	newRoom1 := NewThreadRoom(title, members, false)
	newRoom2 := NewThreadRoom(title, members, true)
	newRoom3 := NewThreadRoom(title, members, false)
	newRoom4 := NewThreadRoom(title, members, true)

	newManager.AddThread(newRoom1)
	newManager.AddThread(newRoom2)
	newManager.AddThread(newRoom3)
	newManager.AddThread(newRoom4)

	rooms := newManager.GetAllRooms()

	if len(rooms) != 2 {
		t.Error()
	}
}

func TestMembers(t *testing.T) {
	title := "test"
	members := []string{"a", "b"}
	newManager := NewThreadManager()
	newRoom := NewThreadRoom(title, []string{}, false)

	id := newRoom.Id
	newManager.AddThread(newRoom)

	newManager.AddMembers(id, members)
	updatedMembers := newManager.GetRoomInfo(id).Members
	if len(updatedMembers) != 2 {
		t.Errorf("Could not add members\n")
	}

	newManager.RemoveMember(id, "a")

	updatedMembers = newManager.GetRoomInfo(id).Members
	if !updatedMembers["b"] || updatedMembers["a"] {
		t.Errorf("Removal error\n")
	}
}

func TestSubscriptions(t *testing.T) {
	title := "test"
	members1 := []string{"a", "b"}
	members2 := []string{"b", "c"}
	newManager := NewThreadManager()
	newRoom1 := NewThreadRoom(title, members1, false)
	newRoom2 := NewThreadRoom(title, members2, false)

	newManager.AddThread(newRoom1)
	newManager.AddThread(newRoom2)

	bSubs := newManager.GetSubscribedRooms("b")
	aSubs := newManager.GetSubscribedRooms("c")
	if len(bSubs) != 2 || len(aSubs) != 1 {
		t.Errorf("Wrong number of subs\n")
	}
}

func TestGetThreadsById(t *testing.T) {
	title := "test"
	members := []string{"a", "b", "c"}
	newManager := NewThreadManager()
	newRoom1 := NewThreadRoom(title, members, false)
	newRoom2 := NewThreadRoom(title, members, false)

	id1 := newRoom1.Id
	newManager.AddThread(newRoom1)
	id2 := newRoom2.Id
	newManager.AddThread(newRoom2)

	ids := []string{id1, id2}

	rooms := newManager.GetThreadsById(ids)

	if len(rooms) != 2 {
		t.Error()
	}
}
