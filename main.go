package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

func allowOrigin(r *http.Request) bool {

	return true

}
func main() {
	server, err := socketio.NewServer(&engineio.Options{

		Transports: []transport.Transport{

			&polling.Transport{

				Client: &http.Client{

					Timeout: time.Minute,
				},

				CheckOrigin: allowOrigin,
			},

			&websocket.Transport{

				CheckOrigin: allowOrigin,
			},
		},
	})

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

	http.Handle("/socket.io/", server)
	log.Println("Serving at localhost:8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
