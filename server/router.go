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
	server.router.POST("/rooms", server.addRooms)
	server.router.GET("/room", server.getRoomInfo)

	return server
}

func (s *HTTPServer) AddSocketRoutes(server *socketio.Server) {
	s.router.GET("/socket.io/*any", gin.WrapH(server))
	s.router.POST("/socket.io/*any", gin.WrapH(server))
}

func (s *HTTPServer) getRooms(g *gin.Context) {
	chatRoomNames, chatRoomIds := s.threadManager.GetAllRooms()
	g.JSON(200, gin.H{
		"rooms": chatRoomNames,
		"ids":   chatRoomIds,
	})
	fmt.Println("Grabbed Rooms")
}

func (s *HTTPServer) addRooms(g *gin.Context) {
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
	fmt.Printf("HEre")
	g.JSON(200, *dto.ToRoomDTO(room))
}
