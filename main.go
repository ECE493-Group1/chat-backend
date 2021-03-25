package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

var chatRoomMap = make(map[string]*ChatRoom)
var chatRoomNames = make([]string, 0)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
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

func main() {
	var port = ":8000"
	router := gin.New()
	server, err := socketio.NewServer(nil)

	if err != nil {
		log.Fatal(err)
	}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		s.Emit("connected", "")
		return nil
	})

	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		fmt.Println("User joined")
		s.Join(msg)
	})

	server.OnEvent("/", "outgoing", func(s socketio.Conn, msg string) {
		fmt.Println("Message Received:", msg)
		server.BroadcastToRoom("/", "messages", "incoming", msg)
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})
	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
	go server.Serve()
	defer server.Close()

	router.Use(GinMiddleware("http://localhost:8080"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))

	router.GET("/rooms", func(g *gin.Context) {
		g.JSON(200, gin.H{
			"rooms": chatRoomNames,
		})
		fmt.Println("Grabbed Rooms")
	})

	router.POST("/rooms", func(g *gin.Context) {
		var newRoomDTO NewRoomDTO
		err := g.BindJSON(&newRoomDTO)
		if err != nil {
			log.Fatal(err)
		}
		chatRoomMap[newRoomDTO.Name] = NewChatRoom(newRoomDTO.Name)
		chatRoomNames = append(chatRoomNames, newRoomDTO.Name)
		fmt.Println(newRoomDTO.Name)
	})
	router.Run(port)
}
