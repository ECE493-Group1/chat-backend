package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"catchat.com/config"
	dto "catchat.com/dtos"
	"catchat.com/models/threads"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupThreadManager() *threads.ThreadManager {
	title := "test"
	members := []string{"a", "b", "c"}

	newRoom1 := threads.NewThreadRoom(title, members, false)
	newRoom2 := threads.NewThreadRoom(title, members, true)
	newRoom3 := threads.NewThreadRoom(title, members, false)
	newRoom4 := threads.NewThreadRoom(title, members, true)

	newManager := threads.NewThreadManager()

	newManager.AddThread(newRoom1)
	newManager.AddThread(newRoom2)
	newManager.AddThread(newRoom3)
	newManager.AddThread(newRoom4)
	return newManager
}

func TestServerInit(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false
	server := NewHTTPServer(setupThreadManager(), config)

	server.router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestGetRooms(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false
	server := NewHTTPServer(setupThreadManager(), config)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms", nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestSubscriptions(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	server := NewHTTPServer(setupThreadManager(), config)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/subscribed?username=a", nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAddRoom(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	server := NewHTTPServer(setupThreadManager(), config)

	newRoom := dto.NewRoomDTO{
		Title:     "Test",
		IsPrivate: false,
		Members:   []string{"e", "f"},
	}
	body, _ := json.Marshal(newRoom)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/rooms", bytes.NewReader(body))
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	rooms := server.threadManager.GetAllRooms()
	assert.Equal(t, 3, len(rooms))
}

func TestGetRoomInfo(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	manager := setupThreadManager()
	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	manager.AddThread(newRoom)

	server := NewHTTPServer(manager, config)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/room?id=%s", id), nil)
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	req, _ = http.NewRequest("GET", "/room", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	req, _ = http.NewRequest("GET", "/room?id=test", nil)
	w = httptest.NewRecorder()
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestUpdateRoom(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	manager := setupThreadManager()
	server := NewHTTPServer(manager, config)
	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	manager.AddThread(newRoom)

	roomDTO := dto.UpdateRoomDTO{
		Id:      id,
		Members: []string{"e", "f"},
	}
	body, _ := json.Marshal(roomDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/room", bytes.NewReader(body))
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestLeaveRoom(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	manager := setupThreadManager()
	server := NewHTTPServer(manager, config)
	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	manager.AddThread(newRoom)

	roomDTO := dto.LeaveRoomDTO{
		Id:     id,
		Member: "e",
	}
	body, _ := json.Marshal(roomDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/leave", bytes.NewReader(body))
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestRoomList(t *testing.T) {
	config := config.GetConfig()
	config.KeywordUpdates = false

	manager := setupThreadManager()
	server := NewHTTPServer(manager, config)
	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	manager.AddThread(newRoom)

	roomDTO := dto.RoomListRequestDTO{
		Ids: []string{id},
	}
	body, _ := json.Marshal(roomDTO)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/roomlist", bytes.NewReader(body))
	server.router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}
