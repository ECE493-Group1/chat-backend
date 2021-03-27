package server

import (
	"encoding/json"
	"fmt"
	"log"

	dto "catchat.com/dtos"
	threads "catchat.com/models/threads"
	socketio "github.com/googollee/go-socket.io"
)

type SocketServer struct {
	server        *socketio.Server
	threadManager *threads.ThreadManager
}

func NewSocketServer(t *threads.ThreadManager) *SocketServer {
	server := &SocketServer{threadManager: t}
	server.init()
	return server
}

func (s *SocketServer) init() {
	var err error
	s.server, err = socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	s.server.OnConnect("/", s.onConnect)

	s.server.OnEvent("/", "JOIN", s.onJoin)
	s.server.OnEvent("/", "SEND", s.onSend)

	s.server.OnEvent("/", "LEAVE", s.onLeave)
	s.server.OnError("/", s.onError)
	s.server.OnDisconnect("/", s.onDisconnect)
}

func (server *SocketServer) onConnect(s socketio.Conn) error {
	s.SetContext("")
	fmt.Println("connected:", s.ID())
	s.Emit("connected", "")
	return nil
}

func (server *SocketServer) onJoin(s socketio.Conn, msg string) {
	fmt.Println("User joined: ", msg)
	s.Join(msg)
}

func (server *SocketServer) onSend(s socketio.Conn, data string) {
	var message dto.MessageDTO
	err := json.Unmarshal([]byte(data), &message)
	if err != nil {
		fmt.Println("Error unmarshaling ", data)
	}
	fmt.Println("Message Received:", message.Content)
	server.server.BroadcastToRoom("/", message.RoomId, "RECEIVE", message.Content)
}

func (server *SocketServer) onLeave(s socketio.Conn) string {
	last := s.Context().(string)
	s.Emit("bye", last)
	s.Close()
	return last
}

func (server *SocketServer) onError(s socketio.Conn, e error) {
	fmt.Println("Socket error:", e)
}

func (server *SocketServer) onDisconnect(s socketio.Conn, reason string) {
	fmt.Println("closed", reason)
}
