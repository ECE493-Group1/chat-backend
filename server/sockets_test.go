package server

import (
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"testing"

	dto "catchat.com/dtos"
	threads "catchat.com/models/threads"
	"github.com/stretchr/testify/mock"
)

// Conn is an autogenerated mock type for the Conn type
type Conn struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Conn) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Context provides a mock function with given fields:
func (_m *Conn) Context() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// Emit provides a mock function with given fields: msg, v
func (_m *Conn) Emit(msg string, v ...interface{}) {
	var _ca []interface{}
	_ca = append(_ca, msg)
	_ca = append(_ca, v...)
	_m.Called(_ca...)
}

// ID provides a mock function with given fields:
func (_m *Conn) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Join provides a mock function with given fields: room
func (_m *Conn) Join(room string) {
	_m.Called(room)
}

// Leave provides a mock function with given fields: room
func (_m *Conn) Leave(room string) {
	_m.Called(room)
}

// LeaveAll provides a mock function with given fields:
func (_m *Conn) LeaveAll() {
	_m.Called()
}

// LocalAddr provides a mock function with given fields:
func (_m *Conn) LocalAddr() net.Addr {
	ret := _m.Called()

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// Namespace provides a mock function with given fields:
func (_m *Conn) Namespace() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RemoteAddr provides a mock function with given fields:
func (_m *Conn) RemoteAddr() net.Addr {
	ret := _m.Called()

	var r0 net.Addr
	if rf, ok := ret.Get(0).(func() net.Addr); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(net.Addr)
		}
	}

	return r0
}

// RemoteHeader provides a mock function with given fields:
func (_m *Conn) RemoteHeader() http.Header {
	ret := _m.Called()

	var r0 http.Header
	if rf, ok := ret.Get(0).(func() http.Header); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(http.Header)
		}
	}

	return r0
}

// Rooms provides a mock function with given fields:
func (_m *Conn) Rooms() []string {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// SetContext provides a mock function with given fields: v
func (_m *Conn) SetContext(v interface{}) {
	_m.Called(v)
}

// URL provides a mock function with given fields:
func (_m *Conn) URL() url.URL {
	ret := _m.Called()

	var r0 url.URL
	if rf, ok := ret.Get(0).(func() url.URL); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(url.URL)
	}

	return r0
}

func TestSocketInit(t *testing.T) {
	threads := setupThreadManager()
	NewSocketServer(threads)
}

func TestOnConnect(t *testing.T) {
	threads := setupThreadManager()
	socketServer := NewSocketServer(threads)
	testConnection := new(Conn)
	testConnection.On("SetContext", "").Return(nil)
	testConnection.On("ID").Return("")
	testConnection.On("Emit", "connected", "").Return(nil)
	socketServer.onConnect(testConnection)
}

func TestOnEnter(t *testing.T) {
	threadManager := setupThreadManager()
	socketServer := NewSocketServer(threadManager)
	testConnection := new(Conn)

	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	threadManager.AddThread(newRoom)

	testConnection.On("Join", id).Return(nil)
	testConnection.On("Emit", "RECEIVE_PREV", dto.ToMessageListDTO([]threads.Message{})).Return(nil)
	socketServer.onEnter(testConnection, id)
}

func TestOnSend(t *testing.T) {
	threadManager := setupThreadManager()
	socketServer := NewSocketServer(threadManager)
	testConnection := new(Conn)

	newRoom := threads.NewThreadRoom("test", []string{"e"}, false)
	id := newRoom.Id
	threadManager.AddThread(newRoom)

	newMessage := dto.NewMessageDTO{
		Content:  "test",
		Username: "e",
		RoomId:   id,
	}

	body, _ := json.Marshal(newMessage)

	testConnection.On("BroadcastToRoom", "/", id, "RECEIVE", "test").Return(nil)
	socketServer.onSend(testConnection, string(body))
}

func TestOnLeave(t *testing.T) {
	threads := setupThreadManager()
	socketServer := NewSocketServer(threads)
	testConnection := new(Conn)
	testConnection.On("Context").Return("")
	testConnection.On("Emit", "bye", "").Return(nil)
	testConnection.On("Close").Return(nil)
	socketServer.onLeave(testConnection)
}