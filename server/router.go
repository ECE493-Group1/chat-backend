package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"catchat.com/config"
	dto "catchat.com/dtos"
	"catchat.com/models/threads"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type HTTPServer struct {
	router        *gin.Engine
	threadManager *threads.ThreadManager
	config        *config.Config
}

func ginMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func NewHTTPServer(t *threads.ThreadManager, config *config.Config) *HTTPServer {
	server := &HTTPServer{threadManager: t, config: config}
	server.router = gin.New()
	server.router.Use(ginMiddleware(config.CORSOrigin))

	server.router.GET("/rooms", server.getRooms)
	server.router.GET("/subscribed", server.getSubscribedRooms)
	server.router.POST("/rooms", server.addRoom)
	server.router.GET("/room", server.getRoomInfo)
	server.router.POST("/room", server.updateRoomInfo)
	server.router.POST("/roomlist", server.getRoomsById)
	server.router.POST("/leave", server.leave)
	go server.keywordUpdateLoop()
	return server
}

func (s *HTTPServer) AddSocketRoutes(server *socketio.Server) {
	s.router.GET("/socket.io/*any", gin.WrapH(server))
	s.router.POST("/socket.io/*any", gin.WrapH(server))
}

func (s *HTTPServer) getRooms(g *gin.Context) {
	publicRooms := s.threadManager.GetAllRooms()
	roomListDTO := dto.ToRoomListDTO(publicRooms)
	g.JSON(200, roomListDTO)
	fmt.Println("Get all threads")
}

func (s *HTTPServer) getSubscribedRooms(g *gin.Context) {
	username := g.Query("username")
	if username == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter"})
		fmt.Printf("Could not find parameter")
		return
	}
	subbedRooms := s.threadManager.GetSubscribedRooms(username)
	subbedListDTO := dto.ToRoomListDTO(subbedRooms)

	g.JSON(200, subbedListDTO)
	fmt.Printf("Get %s subscriptions\n", username)
}

func (s *HTTPServer) addRoom(g *gin.Context) {
	var newRoomDTO dto.NewRoomDTO
	err := g.BindJSON(&newRoomDTO)
	if err != nil {
		log.Fatal(err)
	}
	newRoom := threads.NewThreadRoom(newRoomDTO.Title, newRoomDTO.Members, newRoomDTO.IsPrivate)
	s.threadManager.AddThread(newRoom)
	fmt.Println(newRoomDTO.Title)

	g.JSON(200, dto.RoomRequestDTO{
		Id: newRoom.Id,
	})
}

func (s *HTTPServer) getRoomInfo(g *gin.Context) {
	roomId := g.Query("id")
	if roomId == "" {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameter"})
		fmt.Printf("Could not find parameter")
		return
	}
	room := s.threadManager.GetRoomInfo(roomId)
	if room == nil {
		g.JSON(http.StatusBadRequest, gin.H{"error": "Id not Found"})
		fmt.Printf("Could not find room")
		return
	}
	fmt.Printf("Got room info for %s\n", roomId)
	g.JSON(200, *dto.ToRoomDTO(room))
}

// Adds users to room only. Does not remove
func (s *HTTPServer) updateRoomInfo(g *gin.Context) {
	var update dto.UpdateRoomDTO
	err := g.BindJSON(&update)
	if err != nil {
		log.Fatal(err)
	}
	s.threadManager.AddMembers(update.Id, update.Members)
	fmt.Printf("Added users to %s\n", update.Id)
	g.JSON(200, gin.H{})
}

func (s *HTTPServer) leave(g *gin.Context) {
	var leave dto.LeaveRoomDTO
	err := g.BindJSON(&leave)
	if err != nil {
		log.Fatal(err)
	}
	s.threadManager.RemoveMember(leave.Id, leave.Member)
	fmt.Printf("Removed user %s\n", leave.Member)
	g.JSON(200, gin.H{})
}

func (s *HTTPServer) getRoomsById(g *gin.Context) {
	var ids dto.RoomListRequestDTO
	err := g.BindJSON(&ids)
	if err != nil {
		log.Fatal(err)
	}
	rooms := s.threadManager.GetThreadsById(ids.Ids)
	roomListDTO := dto.ToRoomListDTO(rooms)
	fmt.Printf("Getting rooms by Id\n")
	g.JSON(200, roomListDTO)
}

// Sends Messages to the keyword process. This is run as
func (s *HTTPServer) keywordUpdateLoop() {
	t := s.threadManager
	for {
		item, err := t.MessageQueue.Get(1)
		if err != nil {
			fmt.Printf(("Error grabbing from keyword queue\n"))
			continue
		}
		message := item[0].(threads.Message)
		json_data, err := json.Marshal(dto.ToKeyWordUpdateDTO(message.Content, message.RoomId))
		if err != nil {
			fmt.Printf("Could not convert to JSON\n")
		}

		resp, err := http.Post(s.config.KeywordEndpoint, "application/json", bytes.NewReader(json_data))

		if err != nil || resp.StatusCode != 200 {
			fmt.Printf("Error sending message to keyword service\n")
		}
		fmt.Printf("Submitted message to KW service\n")

	}
}
