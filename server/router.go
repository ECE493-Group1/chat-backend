package server

import (
	"fmt"
	"log"
	"net/http"

	dto "catchat.com/dtos"
	"catchat.com/models/threads"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type HTTPServer struct {
	router        *gin.Engine
	threadManager *threads.ThreadManager
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

func NewHTTPServer(t *threads.ThreadManager, CORSOrigin string) *HTTPServer {
	server := &HTTPServer{threadManager: t}
	server.router = gin.New()
	server.router.Use(ginMiddleware(CORSOrigin))

	server.router.GET("/rooms", server.getRooms)
	server.router.GET("/subscribed", server.getSubscribedRooms)
	server.router.POST("/rooms", server.addRoom)
	server.router.GET("/room", server.getRoomInfo)
	server.router.POST("/room", server.updateRoomInfo)
	server.router.POST("/leave", server.leave)
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
