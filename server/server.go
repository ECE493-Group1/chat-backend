package server

import (
	"catchat.com/config"
	"catchat.com/models/threads"
)

func Init() {
	config := config.GetConfig()
	threadManager := threads.NewThreadManager()
	socketServer := NewSocketServer(threadManager)
	go socketServer.server.Serve()
	defer socketServer.server.Close()

	httpServer := NewHTTPServer(threadManager, config.CORSOrigin)
	httpServer.AddSocketRoutes(socketServer.server)
	httpServer.router.Run(config.Port)
}
